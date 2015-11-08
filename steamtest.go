// steamtest.go - testing full list and individual retrieval of steam server data
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"steamtest/steam"
	"steamtest/steam/filters"
	"strconv"
	"sync"
)

type retryType int

func main() {
	//singleTest("85.229.197.211:25797", steam.QueryTimeout)

	errch := make(chan error, 1)
	go retrieve(errch, filters.NewFilter(filters.SrAll,
		[]filters.SrvFilter{filters.GameQuakeLive},
		[]filters.IgnoredRequest{}))

	err := <-errch
	if err == nil {
		fmt.Print("Got no errors from retrieve goroutine!")
	} else {
		fmt.Printf("Retrieve channel error: %s", err)
	}
}

func singleTest(host string, timeout int) {
	ii, err := steam.GetServerInfoWithLiveData(host, timeout)
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

	pi, err := steam.GetPlayersWithLiveData(host, timeout)
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

	ri, err := steam.GetRulesWithLiveData(host, timeout)
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
			serverinfo, err := steam.GetServerInfoWithLiveData(host, steam.QueryTimeout)
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
	retried := retryFailedInfoReq(failed, 3)
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
			players, err := steam.GetPlayersWithLiveData(host, steam.QueryTimeout)
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
	retried := retryFailedPlayersReq(failed, steam.QueryRetryCount)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

func retryFailedRulesReq(failed []string, rtcount int) map[string]map[string]string {
	m := make(map[string]map[string]string)
	var f []string
	var wg sync.WaitGroup
	var mut sync.Mutex
	for i := 0; i < rtcount; i++ {
		if i == 0 {
			f = failed
		}
		wg.Add(len(f))
		for _, host := range f {
			go func(h string) {
				defer wg.Done()
				r, err := steam.GetRulesWithLiveData(h, steam.QueryTimeout)
				if err != nil {
					if err != steam.NoRulesError {
						fmt.Printf("Host: %s failed on retry-rules request.\n", h)
						return
					}
				}
				mut.Lock()
				m[h] = r
				f = removeFailedHost(f, h)
				mut.Unlock()
			}(host)
		}
		wg.Wait()
	}
	return m
}

func retryFailedPlayersReq(failed []string, rtcount int) map[string][]*steam.PlayerInfo {
	m := make(map[string][]*steam.PlayerInfo)
	var f []string
	var wg sync.WaitGroup
	var mut sync.Mutex
	for i := 0; i < rtcount; i++ {
		if i == 0 {
			f = failed
		}
		wg.Add(len(f))
		for _, host := range f {
			go func(h string) {
				defer wg.Done()
				r, err := steam.GetPlayersWithLiveData(h, steam.QueryTimeout)
				if err != nil {
					if err != steam.NoPlayersError {
						fmt.Printf("Host: %s failed on players-retry request.\n", h)
						return
					}
				}
				mut.Lock()
				m[h] = r
				f = removeFailedHost(f, h)
				mut.Unlock()
			}(host)
		}
		wg.Wait()
	}
	return m
}

func retryFailedInfoReq(failed []string, rtcount int) map[string]*steam.ServerInfo {
	m := make(map[string]*steam.ServerInfo)
	var f []string
	var wg sync.WaitGroup
	var mut sync.Mutex
	for i := 0; i < rtcount; i++ {
		if i == 0 {
			f = failed
		}
		wg.Add(len(f))
		for _, host := range f {
			go func(h string) {
				defer wg.Done()
				r, err := steam.GetServerInfoWithLiveData(h, steam.QueryTimeout)
				if err != nil {
					if err != steam.NoPlayersError {
						fmt.Printf("Host: %s failed on info-retry request.\n", h)
						return
					}
				}
				mut.Lock()
				m[h] = r
				f = removeFailedHost(f, h)
				mut.Unlock()
			}(host)
		}
		wg.Wait()
	}
	return m
}

func removeFailedHost(failed []string, host string) []string {
	for i, v := range failed {
		if v == host {
			failed = append(failed[:i], failed[i+1:]...)
			fmt.Printf("removeFailedHost: removed: %s\n", host)
			fmt.Printf("removeFailedHost: new failed length: %d\n", len(failed))
			break
		}
	}
	return failed
}

func getRulesForServers(serverlist []string) map[string]map[string]string {
	m := make(map[string]map[string]string)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string
	for _, h := range serverlist {
		wg.Add(1)
		go func(host string) {
			rules, err := steam.GetRulesWithLiveData(host, steam.QueryTimeout)
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
	retried := retryFailedRulesReq(failed, steam.QueryRetryCount)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

type serverList struct {
	ServerCount   int       `json:"serverCount"`
	Servers       []*server `json:"servers"`
	FailedCount   int       `json:"failedCount"`
	FailedServers []string  `json:"failedServers"`
}

type server struct {
	Host string `json:"address"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
	// 'Info' by default was *steam.ServerInfo, but nil pointers are encoded as
	// 'null' in JSON instead of an empty object, so use interface and handle appropriately
	Info    interface{}         `json:"info"`
	Players []*steam.PlayerInfo `json:"players"`
	Rules   map[string]string   `json:"rules"`
}

func buildServerList(filter *filters.Filter, servers []string,
	infomap map[string]*steam.ServerInfo, rulemap map[string]map[string]string,
	playermap map[string][]*steam.PlayerInfo) *serverList {

	// No point in ignoring all three requests
	if filter.HasIgnoreInfo && filter.HasIgnorePlayers && filter.HasIgnoreRules {
		return nil
	}

	sl := &serverList{
		Servers: make([]*server, 0),
	}

	var success bool
	var useEmptyInfo bool
	successcount := 0
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
	return sl
}

func retrieve(errors chan<- error, filter *filters.Filter) {
	defer close(errors)
	servers, err := steam.GetServerListWithLiveData(filter)
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

	serverlist := buildServerList(filter, servers, info, rules, players)
	j, err := json.Marshal(serverlist)
	if err != nil {
		errors <- fmt.Errorf("Error marshaling json: %s", err)
		return
	}
	file, err := os.Create("servers.json")
	if err != nil {
		errors <- fmt.Errorf("Error creating json file: %s", err)
		return
	}
	defer file.Close()
	file.Sync()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(j)
	if err != nil {
		errors <- fmt.Errorf("Error writing json file: %s", err)
		return
	}
	writer.Flush()
}
