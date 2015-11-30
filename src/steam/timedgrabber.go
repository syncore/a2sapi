package steam

// timedgrabber.go - Timed retrieval of servers from the Steam Master server.

import (
	"bufio"
	"encoding/json"
	"os"
	"steamtest/src/models"
	"steamtest/src/steam/filters"
	"steamtest/src/util"
	"time"
)

func retrieve(filter *filters.Filter) error {
	mq, err := NewMasterQuery(filter)
	if err != nil {
		return util.LogSteamErrorf("Master server error: %s", err)
	}

	if filter.Game.IgnoreInfo && filter.Game.IgnorePlayers && filter.Game.IgnoreRules {
		return util.LogAppErrorf("Cannot ignore all three AS2 requests!")
	}

	var players map[string][]*models.SteamPlayerInfo
	var rules map[string]map[string]string
	var info map[string]*models.SteamServerInfo
	// Order of retrieval is by amount of work that must be done (1 = 2, 3)
	// 1. players (request chal #, recv chal #, req players, recv players)
	// 2. rules (request chal #, recv chal #, req rules, recv rules)
	// 3. info: just request info & receive info
	// Note: some servers (i.e. new beta games) don't have all 3 of AS2_RULES/PLAYER/INFO
	if !filter.Game.IgnorePlayers {
		players = batchPlayerQuery(mq.Servers)
	}
	if !filter.Game.IgnoreRules {
		rules = batchRuleQuery(mq.Servers)
	}
	if !filter.Game.IgnoreInfo {
		info = batchInfoQuery(mq.Servers)
	}

	serverlist, err := buildMasterServerList(filter.Game, mq.Servers, info, rules,
		players)
	if err != nil {
		return util.LogAppError(err)
	}

	j, err := json.Marshal(serverlist)
	if err != nil {
		return util.LogAppErrorf("Error marshaling json: %s", err)
	}
	file, err := os.Create("servers.json")
	if err != nil {
		return util.LogAppErrorf("Error creating json file: %s", err)
	}
	defer file.Close()
	file.Sync()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(j)
	if err != nil {
		return util.LogAppErrorf("Error writing json file: %s", err)
	}
	writer.Flush()
	return nil
}

// StartMasterRetrieval starts a timed retrieval of servers specified by a given
// filter from the Steam Master server after an initial delay of initialDelay
// seconds. It retrieves the list every timeBetweenQueries seconds thereafter.
// A bool can be sent to the stop channel to cancel all timed retrievals.
func StartMasterRetrieval(stop chan bool, filter *filters.Filter,
	initialDelay int, timeBetweenQueries int) {
	retrticker := time.NewTicker(time.Duration(timeBetweenQueries) * time.Second)
	util.LogAppInfo("Waiting %d seconds before attempting first retrieval...",
		initialDelay)

	firstretrieval := time.NewTimer(time.Duration(initialDelay) * time.Second)
	<-firstretrieval.C
	_ = retrieve(filter)

	for {
		select {
		case <-retrticker.C:
			go func(*filters.Filter) {
				util.LogAppInfo("%s: Starting %s master server query", time.Now().Format(
					"Mon Jan _2 15:04:05 2006 EST"), filter.Game.Name)
				_ = retrieve(filter)
			}(filter)
		case <-stop:
			retrticker.Stop()
			return
		}
	}
}
