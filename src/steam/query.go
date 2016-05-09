package steam

// query.go - Used for querying individual game servers to retrieve their info
// for building a list to return to the API

import (
	"sync"

	"github.com/syncore/a2sapi/src/logger"
	"github.com/syncore/a2sapi/src/models"
	"github.com/syncore/a2sapi/src/steam/filters"
)

type a2sData struct {
	HostsGames map[string]filters.Game
	Info       map[string]models.SteamServerInfo
	Rules      map[string]map[string]string
	Players    map[string][]models.SteamPlayerInfo
}

func batchInfoQuery(servers []string) map[string]models.SteamServerInfo {
	m := make(map[string]models.SteamServerInfo)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string

	for _, h := range servers {
		wg.Add(1)
		go func(host string) {
			serverinfo, err := GetInfoForServer(host, QueryTimeout)
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
	retried := RetryFailedInfoReq(failed, 3)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

func batchPlayerQuery(servers []string) map[string][]models.SteamPlayerInfo {
	m := make(map[string][]models.SteamPlayerInfo)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string

	for _, h := range servers {
		wg.Add(1)
		go func(host string) {
			players, err := GetPlayersForServer(host, QueryTimeout)
			if err != nil {
				// server could just be empty
				if err != ErrNoPlayers {
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
	retried := RetryFailedPlayersReq(failed, QueryRetryCount)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

func batchRuleQuery(servers []string) map[string]map[string]string {
	m := make(map[string]map[string]string)
	var wg sync.WaitGroup
	var mut sync.Mutex
	var failed []string
	for _, h := range servers {
		wg.Add(1)
		go func(host string) {
			rules, err := GetRulesForServer(host, QueryTimeout)
			if err != nil {
				// server might have no rules
				if err != ErrNoRules {
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
	retried := RetryFailedRulesReq(failed, QueryRetryCount)
	for k, v := range retried {
		m[k] = v
	}
	return m
}

// DirectQuery allows a user to query any host even if it is not in the internal
// server ID database. It is primarily intended for testing as it has two main
// issues: 1) obvious security implications, 2) determining which game a user-
// supplied host represents rests on potentially unreliable assumptions, which if
// not true would cause games with incomplete support for all three A2S queries
// (e.g. Reflex) to always fail. A production environment should use Query() instead.
func DirectQuery(hosts []string) (*models.APIServerList, error) {
	hg := make(map[string]filters.Game, len(hosts))

	// Try to account for the fact that we can't determine the game ahead of time
	// for user-specified direct host queries -- a number of assumptions:
	// (1) A2S_INFO for game/host, (2) extra data A2S_INFO flag & field w/ appid,
	//(3) game has been defined in game.go with the correct AppID and A2S ignore flags
	info := batchInfoQuery(hosts)
	needsRules := make([]string, len(hosts))
	needsPlayers := make([]string, len(hosts))

	for _, h := range hosts {
		logger.WriteDebug("direct query for %s. will try to figure out needed queries", h)
		if (info[h] != models.SteamServerInfo{}) {
			logger.WriteDebug("A2S_INFO not empty. got gameid: %d", info[h].ExtraData.GameID)
			fg := filters.GetGameByAppID(info[h].ExtraData.GameID)
			hg[h] = fg
			if !fg.IgnoreRules {
				logger.WriteDebug("based on game %s for %s, will need to get A2S_RULES",
					fg.Name, h)
				needsRules = append(needsRules, h)
			}
			if !fg.IgnorePlayers {
				logger.WriteDebug("based on game %s for %s, will need to get A2S_PLAYERS",
					fg.Name, h)
				needsPlayers = append(needsPlayers, h)
			}
		} else {
			logger.WriteDebug("A2S_INFO is nil. game will be unspecified; results may vary")
			hg[h] = filters.GameUnspecified
		}
	}
	data := a2sData{
		HostsGames: hg,
		Info:       info,
		Rules:      batchRuleQuery(needsRules),
		Players:    batchPlayerQuery(needsPlayers),
	}
	sl, err := buildServerList(data, true)
	if err != nil {
		return models.GetDefaultServerList(), logger.LogAppError(err)
	}
	return sl, nil
}

// Query retrieves the server information for a given set of host to game pairs
// and returns it in a format that is presented to the API. It takes a map consisting
// of host(s) and their corresponding game names (i.e: k:127.0.0.1:27960, v:"QuakeLive")
func Query(hostsgames map[string]string) (*models.APIServerList, error) {
	hg := make(map[string]filters.Game, len(hostsgames))
	needsPlayers := make([]string, len(hostsgames))
	needsRules := make([]string, len(hostsgames))
	needsInfo := make([]string, len(hostsgames))

	for host, game := range hostsgames {
		fg := filters.GetGameByName(game)
		hg[host] = fg
		if !fg.IgnoreRules {
			needsRules = append(needsRules, host)
		}
		if !fg.IgnorePlayers {
			needsPlayers = append(needsPlayers, host)
		}
		if !fg.IgnoreInfo {
			needsInfo = append(needsInfo, host)
		}
	}
	data := a2sData{
		HostsGames: hg,
		Info:       batchInfoQuery(needsInfo),
		Rules:      batchRuleQuery(needsRules),
		Players:    batchPlayerQuery(needsPlayers),
	}

	sl, err := buildServerList(data, true)
	if err != nil {
		return models.GetDefaultServerList(), logger.LogAppError(err)
	}
	return sl, nil
}
