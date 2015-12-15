package steam

// timedgrabber.go - Timed retrieval of servers from the Steam Master server.

import (
	"bufio"
	"encoding/json"
	"os"
	"steamtest/src/logger"
	"steamtest/src/steam/filters"
	"time"
)

func retrieve(filter *filters.Filter) error {
	mq, err := NewMasterQuery(filter)
	if err != nil {
		return logger.LogSteamErrorf("Master server error: %s", err)
	}

	if filter.Game.IgnoreInfo && filter.Game.IgnorePlayers && filter.Game.IgnoreRules {
		return logger.LogAppErrorf("Cannot ignore all three AS2 requests!")
	}

	data := &a2sData{}
	hg := make(map[string]*filters.Game, len(mq.Servers))
	for _, h := range mq.Servers {
		hg[h] = filter.Game
	}
	data.HostsGames = hg

	// Order of retrieval is by amount of work that must be done (generally 1, 2, 3)
	// 1. rules (request chal #, recv chal #, req rules, recv rules)
	// games with multi-packet A2S_RULES replies do the most work; otherwise 1 = 2, 3
	// 2. players (request chal #, recv chal #, req players, recv players)
	// 3. info: just request info & receive info
	// Note: some servers (i.e. new beta games) don't have all 3 of AS2_RULES/PLAYER/INFO
	if !filter.Game.IgnoreRules {
		data.Rules = batchRuleQuery(mq.Servers)
	}
	if !filter.Game.IgnorePlayers {
		data.Players = batchPlayerQuery(mq.Servers)
	}
	if !filter.Game.IgnoreInfo {
		data.Info = batchInfoQuery(mq.Servers)
	}

	serverlist, err := buildServerList(data, true)
	if err != nil {
		return logger.LogAppError(err)
	}

	// TODO: a debugMode in the configuration which if enabled will dump servers.json
	// if not, then it won't (for when master server list is stored in memory)
	j, err := json.Marshal(serverlist)
	if err != nil {
		return logger.LogAppErrorf("Error marshaling json: %s", err)
	}
	file, err := os.Create("servers.json")
	if err != nil {
		return logger.LogAppErrorf("Error creating json file: %s", err)
	}
	defer file.Close()
	file.Sync()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(j)
	if err != nil {
		return logger.LogAppErrorf("Error writing json file: %s", err)
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

	logger.WriteDebug(
		"Waiting %d seconds before grabbing %s servers. Will retrieve servers every %d secs afterwards.",
		initialDelay, filter.Game.Name, timeBetweenQueries)

	logger.LogAppInfo(
		"Waiting %d seconds before grabbing %s servers from master. Will retrieve every %d secs afterwards.",
		initialDelay, filter.Game.Name, timeBetweenQueries)

	firstretrieval := time.NewTimer(time.Duration(initialDelay) * time.Second)
	<-firstretrieval.C
	_ = retrieve(filter)

	for {
		select {
		case <-retrticker.C:
			go func(*filters.Filter) {
				logger.WriteDebug("%s: Starting %s master server query", time.Now().Format(
					"Mon Jan 2 15:04:05 2006 EST"), filter.Game.Name)
				logger.LogAppInfo("%s: Starting %s master server query", time.Now().Format(
					"Mon Jan 2 15:04:05 2006 EST"), filter.Game.Name)
				_ = retrieve(filter)
			}(filter)
		case <-stop:
			retrticker.Stop()
			return
		}
	}
}
