// steamtest.go - testing full list and individual retrieval of steam server data
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"steamtest/db"
	"steamtest/steam"
	"steamtest/steam/filters"
	"strconv"
	"sync"
)

type requestType int

type serverList struct {
	ServerCount   int       `json:"serverCount"`
	Servers       []*server `json:"servers"`
	FailedCount   int       `json:"failedCount"`
	FailedServers []string  `json:"failedServers"`
}
type server struct {
	Host        string      `json:"address"`
	IP          string      `json:"ip"`
	Port        int         `json:"port"`
	CountryInfo *db.Country `json:"location"`
	// 'Info' by default was *steam.ServerInfo, but nil pointers are encoded as
	// 'null' in JSON instead of an empty object, so use interface and handle appropriately
	Info    interface{}         `json:"info"`
	Players []*steam.PlayerInfo `json:"players"`
	Rules   map[string]string   `json:"rules"`
}

const (
	ruleRequest requestType = iota
	playerRequest
	infoRequest
)

func main() {
	//singleServerTest("85.229.197.211:25797", steam.QueryTimeout)

	errch := make(chan error, 1)
	go retrieve(errch, filters.NewFilter(filters.SrAll,
		[]filters.SrvFilter{filters.GameQuakeLive},
		[]filters.IgnoredRequest{}))

	err := <-errch
	if err == nil {
		fmt.Print("All good! Got no errors from retrieve goroutine!\n")
	} else {
		fmt.Printf("Retrieve channel error: %s\n", err)
	}
}

func singleServerTest(host string, timeout int) {
	ii, err := steam.GetInfoForServer(host, timeout)
	if err != nil && err != steam.NoInfoError {
		fmt.Printf("%s", err)
	} else if err == steam.NoInfoError {
		fmt.Print("no info found error detected\n")
	} else {
		fmt.Printf(`protocol:%d, name:%s, mapname:%s, folder:%s, game:%s, id:%d,
	 players:%d, maxplayers:%d, bots:%d, servertype:%s, environment:%s,
	 visibility:%d, vac:%d, version:%s, port:%d, steamid:%d, sourcetvport:%d,
	 sourcetvname:%s, keywords:%s, gameid:%d
	 `,
			ii.Protocol, ii.Name, ii.Map, ii.Folder,
			ii.Game, ii.ID, ii.Players, ii.MaxPlayers,
			ii.Bots, ii.ServerType, ii.Environment,
			ii.Visibility, ii.VAC, ii.Version,
			ii.ExtraData.Port, ii.ExtraData.SteamID,
			ii.ExtraData.SourceTVPort, ii.ExtraData.SourceTVName,
			ii.ExtraData.Keywords, ii.ExtraData.GameID)
	}

	pi, err := steam.GetPlayersForServer(host, timeout)
	if err != nil && err != steam.NoPlayersError {
		fmt.Printf("%s", err)
	} else if err == steam.NoPlayersError {
		fmt.Print("no players found error detected\n")
	} else {
		for _, p := range pi {
			fmt.Printf("Name: %s, Score: %d, Connected for: %s\n",
				p.Name, p.Score, p.TimeConnectedTot)
		}
	}

	ri, err := steam.GetRulesForServer(host, timeout)
	if err != nil && err != steam.NoRulesError {
		fmt.Printf("%s", err)
	} else if err == steam.NoRulesError {
		fmt.Print("no rules found error detected\n")
	} else {
		for _, r := range ri {
			fmt.Printf("%s\n", r)
		}
	}
}

func getInfoForServers(serverlist []string) map[string]*steam.ServerInfo {
	m := make(map[string]*steam.ServerInfo)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string

	for _, h := range serverlist {
		wg.Add(1)
		go func(host string) {
			serverinfo, err := steam.GetInfoForServer(host, steam.QueryTimeout)
			if err != nil {
				mut.Lock()
				failed = append(failed, host)
				mut.Unlock()
				wg.Done()
				return
			}
			mut.Lock()
			m[host] = serverinfo
			mut.Unlock()
			wg.Done()
		}(h)
	}
	wg.Wait()
	retried := steam.RetryFailedInfoReq(failed, 3)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

func getPlayersForServers(serverlist []string) map[string][]*steam.PlayerInfo {
	m := make(map[string][]*steam.PlayerInfo)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string

	for _, h := range serverlist {
		wg.Add(1)
		go func(host string) {
			players, err := steam.GetPlayersForServer(host, steam.QueryTimeout)
			if err != nil {
				// server could just be empty
				if err != steam.NoPlayersError {
					mut.Lock()
					failed = append(failed, host)
					mut.Unlock()
					wg.Done()
					return
				}
			}
			mut.Lock()
			m[host] = players
			mut.Unlock()
			wg.Done()
		}(h)
	}
	wg.Wait()
	retried := steam.RetryFailedPlayersReq(failed, steam.QueryRetryCount)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

func getRulesForServers(serverlist []string) map[string]map[string]string {
	m := make(map[string]map[string]string)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string
	for _, h := range serverlist {
		wg.Add(1)
		go func(host string) {
			rules, err := steam.GetRulesForServer(host, steam.QueryTimeout)
			if err != nil {
				// server might have no rules
				if err != steam.NoRulesError {
					mut.Lock()
					failed = append(failed, host)
					mut.Unlock()
					wg.Done()
					return
				}
			}
			mut.Lock()
			m[host] = rules
			mut.Unlock()
			wg.Done()
		}(h)
	}
	wg.Wait()
	retried := steam.RetryFailedRulesReq(failed, steam.QueryRetryCount)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

func buildServerList(filter *filters.Filter, servers []string,
	infomap map[string]*steam.ServerInfo, rulemap map[string]map[string]string,
	playermap map[string][]*steam.PlayerInfo) (*serverList, error) {

	// No point in ignoring all three requests
	if filter.HasIgnoreInfo && filter.HasIgnorePlayers && filter.HasIgnoreRules {
		return nil, fmt.Errorf("Cannot ignore all three A2S_ requests!")
	}

	sl := &serverList{
		Servers: make([]*server, 0),
	}

	var success bool
	var useEmptyInfo bool
	successcount := 0
	cdb, err := db.OpenCountryDB()
	if err != nil {
		return nil, err
	}
	defer cdb.Close()

	for _, host := range servers {
		var i interface{}
		info, iok := infomap[host]
		players, pok := playermap[host]
		if players == nil {
			// return empty array instead of nil pointers (null) in json
			players = make([]*steam.PlayerInfo, 0)
		}
		rules, rok := rulemap[host]

		// default, unless we should skip
		success = iok && rok && pok

		if filter.HasIgnoreInfo {
			success = pok && rok
			useEmptyInfo = true
			i = make(map[string]int, 0)
		}
		if filter.HasIgnorePlayers {
			success = iok && rok
		}
		if filter.HasIgnoreRules {
			rules = make(map[string]string, 0)
			success = iok && pok
		}
		if filter.HasIgnoreInfo && filter.HasIgnorePlayers {
			success = rok
		}
		if filter.HasIgnoreInfo && filter.HasIgnoreRules {
			success = pok
		}
		if filter.HasIgnorePlayers && filter.HasIgnoreRules {
			success = iok
		}

		if success {
			srv := &server{
				Players: players,
				Rules:   rules,
			}
			// this is needed to return the omitted info as an empty object in JSON
			if useEmptyInfo {
				srv.Info = i
			} else {
				srv.Info = info
			}

			ip, port, err := net.SplitHostPort(host)
			if err == nil {
				srv.IP = ip
				if info.ExtraData.Port != 0 {
					srv.Host = fmt.Sprintf("%s:%d", ip, info.ExtraData.Port)
				} else {
					srv.Host = host
				}
				p, err := strconv.Atoi(port)
				if err == nil {
					srv.Port = p
				}
				loc := make(chan *db.Country, 1)
				go db.GetCountryInfo(loc, cdb, ip)
				srv.CountryInfo = <-loc
			}
			sl.Servers = append(sl.Servers, srv)
			successcount++
		} else {
			sl.FailedServers = append(sl.FailedServers, host)
		}
	}
	sl.ServerCount = len(sl.Servers)
	sl.FailedCount = len(sl.FailedServers)

	fmt.Printf("%d servers were successfully queried!\n", successcount)
	return sl, nil
}

func retrieve(errors chan<- error, filter *filters.Filter) {
	defer close(errors)
	servers, err := steam.GetServersFromMaster(filter)
	if err != nil {
		errors <- fmt.Errorf("Master server error: %s\n", err)
		return
	}

	if filter.HasIgnoreInfo && filter.HasIgnorePlayers && filter.HasIgnoreRules {
		errors <- fmt.Errorf("Cannot ignore all three AS2 requests!")
		return
	}

	var players map[string][]*steam.PlayerInfo
	var rules map[string]map[string]string
	var info map[string]*steam.ServerInfo
	// Order of retrieval is by amount of work that must be done (1 = 2, 3)
	// 1. players (request chal #, recv chal #, req players, recv players)
	// 2. rules (request chal #, recv chal #, req rules, recv rules)
	// 3. rules: just request rules & receive rules

	// Some servers (i.e. new beta games) don't have all 3 of AS2_RULES/PLAYER/INFO
	if !filter.HasIgnorePlayers {
		players = getPlayersForServers(servers)
	}
	if !filter.HasIgnoreRules {
		rules = getRulesForServers(servers)
	}
	if !filter.HasIgnoreInfo {
		info = getInfoForServers(servers)
	}

	serverlist, err := buildServerList(filter, servers, info, rules, players)
	if err != nil {
		errors <- err
		return
	}

	j, err := json.Marshal(serverlist)
	if err != nil {
		errors <- fmt.Errorf("Error marshaling json: %s\n", err)
		return
	}
	file, err := os.Create("servers.json")
	if err != nil {
		errors <- fmt.Errorf("Error creating json file: %s\n", err)
		return
	}
	defer file.Close()
	file.Sync()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(j)
	if err != nil {
		errors <- fmt.Errorf("Error writing json file: %s\n", err)
		return
	}
	writer.Flush()
}
