package steam

// listbuilder.go - Functions for building the list of servers & their details
// in resposne to a retrieval of all servers from the Steam Master server
// or in response to a user's specific query from the API.

import (
	"database/sql"
	"net"
	"steamtest/src/db"
	"steamtest/src/models"
	"steamtest/src/util"
	"strconv"
	"time"
)

// TODO: large -- refactor
func buildServerList(data *a2sData, addtoServerDB bool) (*models.APIServerList,
	error) {
	// Cannot ignore all three requests
	for _, g := range data.HostsGames {
		if g.IgnoreInfo && g.IgnorePlayers && g.IgnoreRules {
			return nil, util.LogAppErrorf("Cannot ignore all three A2S_ requests!")
		}
	}
	successcount := 0
	var success bool
	var useEmptyInfo bool
	var sdb *sql.DB
	var sdbhosts map[string]string
	sl := &models.APIServerList{
		Servers:       make([]*models.APIServer, 0),
		FailedServers: make([]string, 0),
	}

	if addtoServerDB {
		sdbhosts = make(map[string]string, len(data.HostsGames))
	}

	cdb, err := db.OpenCountryDB()
	if err != nil {
		return nil, util.LogAppError(err)
	}
	defer cdb.Close()

	for host, game := range data.HostsGames {
		var i interface{}
		info, iok := data.Info[host]
		players, pok := data.Players[host]
		if players == nil {
			// return empty array instead of nil pointers (null) in json
			players = make([]*models.SteamPlayerInfo, 0)
		}
		rules, rok := data.Rules[host]
		success = iok && rok && pok

		if game.IgnoreInfo {
			success = pok && rok
			useEmptyInfo = true
			i = make(map[string]int, 0)
		}
		if game.IgnorePlayers {
			success = iok && rok
		}
		if game.IgnoreRules {
			rules = make(map[string]string, 0)
			success = iok && pok
		}
		if game.IgnoreInfo && game.IgnorePlayers {
			success = rok
		}
		if game.IgnoreInfo && game.IgnoreRules {
			success = pok
		}
		if game.IgnorePlayers && game.IgnoreRules {
			success = iok
		}

		if success {
			srv := &models.APIServer{
				Game:        game.Name,
				Players:     players,
				RealPlayers: removeBuggedPlayers(players),
				Rules:       rules,
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
				srv.Host = host
				p, err := strconv.Atoi(port)
				if err == nil {
					srv.Port = p
				}
				if addtoServerDB {
					sdbhosts[host] = game.Name
				}
				loc := make(chan *models.DbCountry, 1)
				go db.GetCountryInfo(loc, cdb, ip)
				srv.CountryInfo = <-loc
			}
			sl.Servers = append(sl.Servers, srv)
			successcount++
		} else {
			util.WriteDebug("failed consists of: %s", host)
			sl.FailedServers = append(sl.FailedServers, host)
		}
	}

	sl.RetrievedAt = time.Now().Format("Mon Jan _2 15:04:05 2006 EST")
	sl.RetrievedTimeStamp = time.Now().Unix()
	sl.ServerCount = len(sl.Servers)
	sl.FailedCount = len(sl.FailedServers)

	if addtoServerDB {
		sdb, err = db.OpenServerDB()
		if err != nil {
			return nil, util.LogAppError(err)
		}
		defer sdb.Close()
		go db.AddServersToDB(sdb, sdbhosts)
		sl = setServerIDForList(sdb, sl)
	}

	util.LogAppInfo("Specific query: %d servers were successfully queried!",
		successcount)
	return sl, nil
}

func removeBuggedPlayers(players []*models.SteamPlayerInfo) *models.RealPlayerInfo {
	rpi := &models.RealPlayerInfo{
		RealPlayerCount: len(players),
		Players:         players,
	}
	cfg, err := util.ReadConfig()
	if err != nil {
		util.LogAppError(err)
		return rpi
	}

	var filtered []*models.SteamPlayerInfo
	for _, p := range players {
		if int(p.TimeConnectedSecs) < (3600 * cfg.SteamConfig.SteamBugPlayerTime) {
			filtered = append(filtered, p)
		}
	}
	// Empty players (nil) displayed as empty array in JSON, not null
	if len(filtered) == 0 {
		rpi.RealPlayerCount = 0
		rpi.Players = make([]*models.SteamPlayerInfo, 0)
	} else {
		rpi.RealPlayerCount = len(filtered)
		rpi.Players = filtered
	}
	return rpi
}

func setServerIDForList(sdb *sql.DB,
	sl *models.APIServerList) *models.APIServerList {
	toSet := make(map[string]string, len(sl.Servers))
	for _, s := range sl.Servers {
		toSet[s.Host] = s.Game
	}
	result := make(chan map[string]int64, 1)
	go db.GetIDsForServerList(result, sdb, toSet)
	m := <-result

	for _, s := range sl.Servers {
		if m[s.Host] != 0 {
			s.ID = m[s.Host]
		}
	}
	return sl
}
