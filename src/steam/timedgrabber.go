package steam

// timedgrabber.go - Timed retrieval of servers from the Steam Master server.

import (
	"a2sapi/src/config"
	"a2sapi/src/constants"
	"a2sapi/src/logger"
	"a2sapi/src/models"
	"a2sapi/src/steam/filters"
	"a2sapi/src/util"
	"encoding/json"
	"fmt"
	"time"
)

func retrieve(filter filters.Filter) (*models.APIServerList, error) {
	mq, err := NewMasterQuery(filter)
	if err != nil {
		return nil, logger.LogSteamErrorf("Master server error: %s", err)
	}

	if filter.Game.IgnoreInfo && filter.Game.IgnorePlayers && filter.Game.IgnoreRules {
		return nil, logger.LogAppErrorf("Cannot ignore all three AS2 requests!")
	}

	data := a2sData{}
	hg := make(map[string]filters.Game, len(mq.Servers))
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
		return nil, logger.LogAppError(err)
	}

	if config.Config.DebugConfig.EnableServerDump {
		if err := dumpServersToDisk(filter.Game.Name, serverlist); err != nil {
			logger.LogAppError(err)
		}
	}

	return serverlist, nil
}

func dumpServersToDisk(gamename string, sl *models.APIServerList) error {
	j, err := json.Marshal(sl)
	if err != nil {
		return logger.LogAppErrorf("Error marshaling json: %s", err)
	}
	t := time.Now()
	if err := util.CreateDirectory(constants.DumpDirectory); err != nil {
		return logger.LogAppErrorf("Couldn't create '%s' dir: %s\n",
			constants.DumpDirectory, err)
	}
	// Windows doesn't allow ":" in filename so use '-' separators for time
	err = util.CreateByteFile(j, constants.DumpFileFullPath(
		fmt.Sprintf("%s-servers-%d-%02d-%02d.%02d-%02d-%02d.json",
			gamename, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())),
		true)
	if err != nil {
		return logger.LogAppErrorf("Error creating server dump file: %s", err)
	}
	return nil
}

// StartMasterRetrieval starts a timed retrieval of servers specified by a given
// filter from the Steam Master server after an initial delay of initialDelay
// seconds. It retrieves the list every timeBetweenQueries seconds thereafter.
// A bool can be sent to the stop channel to cancel all timed retrievals.
func StartMasterRetrieval(stop chan bool, filter filters.Filter,
	initialDelay int, timeBetweenQueries int) {
	retrticker := time.NewTicker(time.Duration(timeBetweenQueries) * time.Second)

	logger.WriteDebug(
		"Waiting %d seconds before grabbing %s servers. Will retrieve servers every %d secs afterwards.", initialDelay, filter.Game.Name, timeBetweenQueries)

	logger.LogAppInfo(
		"Waiting %d seconds before grabbing %s servers from master. Will retrieve every %d secs afterwards.", initialDelay, filter.Game.Name, timeBetweenQueries)

	firstretrieval := time.NewTimer(time.Duration(initialDelay) * time.Second)
	<-firstretrieval.C
	logger.WriteDebug("Starting first retrieval of %s servers from master.",
		filter.Game.Name)
	sl, err := retrieve(filter)
	if err != nil {
		logger.LogAppErrorf("Error when performing timed master retrieval: %s", err)
	}
	models.MasterList = sl

	for {
		select {
		case <-retrticker.C:
			go func(filters.Filter) {
				logger.WriteDebug("%s: Starting %s master server query", time.Now().Format(
					"Mon Jan 2 15:04:05 2006 EST"), filter.Game.Name)
				logger.LogAppInfo("%s: Starting %s master server query", time.Now().Format(
					"Mon Jan 2 15:04:05 2006 EST"), filter.Game.Name)
				sl, err := retrieve(filter)
				if err != nil {
					logger.LogAppErrorf("Error when performing timed master retrieval: %s",
						err)
				}
				models.MasterList = sl
			}(filter)
		case <-stop:
			retrticker.Stop()
			return
		}
	}
}
