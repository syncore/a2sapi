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

	err := make(chan error, 1)
	go retrieve(err, filters.SrAll, filters.GameQuakeLive)
	if <-err == nil {
		fmt.Print("Got no errors from retrieve goroutine!")
	} else {
		fmt.Printf("Retrieve channel error: %s", <-err)
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
	Host    string              `json:"address"`
	IP      string              `json:"ip"`
	Port    int                 `json:"port"`
	Info    *steam.ServerInfo   `json:"info"`
	Players []*steam.PlayerInfo `json:"players"`
	Rules   map[string]string   `json:"rules"`
}

func buildServerList(servers []string, infomap map[string]*steam.ServerInfo,
	rulemap map[string]map[string]string,
	playermap map[string][]*steam.PlayerInfo) *serverList {

	sl := &serverList{
		Servers: make([]*server, 0),
	}
	successcount := 0
	for _, host := range servers {
		info, iok := infomap[host]
		players, pok := playermap[host]
		if players == nil {
			// return empty arrays instead of nil pointers in json
			players = make([]*steam.PlayerInfo, 0)
		}
		rules, rok := rulemap[host]

		if iok && rok && pok {
			//fmt.Printf("Success: all data exists for host: %s\n", host)
			srv := &server{
				Players: players,
				Rules:   rules,
				Info:    info,
			}
			srv.Rules = rules
			srv.Info = info
			srv.Host = host
			ip, port, err := net.SplitHostPort(host)
			if err == nil {
				srv.IP = ip
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

func retrieve(errors chan<- error, region filters.ServerRegion,
	filters ...filters.ServerFilter) {
	defer close(errors)
	servers, err := steam.GetServerListWithLiveData(region, filters...)
	if err != nil {
		errors <- fmt.Errorf("Master server error: %s\n", err)
		return
	}

	// Retrieved by amount of work that must be done (1 = 2, 3)
	// 1. players (request chal #, recv chal #, req players, recv players)
	// 2. rules (request chal #, recv chal #, req rules, recv rules)
	// 3. rules: just request rules & receive rules
	players := getPlayersForServers(servers)
	rules := getRulesForServers(servers)
	info := getInfoForServers(servers)
	serverlist := buildServerList(servers, info, rules, players)

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
