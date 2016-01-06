package steam

// listbuilder.go - Functions for building the list of servers & their details
// in resposne to a retrieval of all servers from the Steam Master server
// or in response to a user's specific query from the API.

import (
	"a2sapi/src/db"
	"a2sapi/src/logger"
	"a2sapi/src/models"
	"database/sql"
	"net"
	"strconv"
	"time"
)

func buildServerList(data *a2sData, addtoServerDB bool) (*models.APIServerList,
	error) {
	// Cannot ignore all three requests
	for _, g := range data.HostsGames {
		if g.IgnoreInfo && g.IgnorePlayers && g.IgnoreRules {
			return nil, logger.LogAppErrorf("Cannot ignore all three A2S_ requests!")
		}
	}
	successcount := 0
	var success bool
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
		return nil, logger.LogAppError(err)
	}
	defer cdb.Close()

	for host, game := range data.HostsGames {
		info, iok := data.Info[host]
		if info == nil {
			info = &models.SteamServerInfo{}
		}
		players, pok := data.Players[host]
		if players == nil {
			// return empty array instead of nil pointers (null) in json
			players = make([]*models.SteamPlayerInfo, 0)
		}
		rules, rok := data.Rules[host]
		success = iok && rok && pok

		if game.IgnoreInfo {
			success = pok && rok
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
				Game:            game.Name,
				Players:         players,
				FilteredPlayers: removeBuggedPlayers(players),
				Rules:           rules,
				Info:            info,
			}

			ip, port, serr := net.SplitHostPort(host)
			if serr == nil {
				srv.IP = ip
				srv.Host = host
				p, perr := strconv.Atoi(port)
				if perr == nil {
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
			sl.FailedServers = append(sl.FailedServers, host)
		}
	}

	sl.RetrievedAt = time.Now().Format("Mon Jan 2 15:04:05 2006 EST")
	sl.RetrievedTimeStamp = time.Now().Unix()
	sl.ServerCount = len(sl.Servers)
	sl.FailedCount = len(sl.FailedServers)

	if addtoServerDB {
		sdb, err = db.OpenServerDB()
		if err != nil {
			return nil, logger.LogAppError(err)
		}
		defer sdb.Close()
		go db.AddServersToDB(sdb, sdbhosts)
		sl = setServerIDForList(sdb, sl)
	}

	logger.LogAppInfo(
		"Successfully queried (%d/%d) servers. %d timed out or otherwise failed.",
		successcount, len(data.HostsGames), sl.FailedCount)
	logger.WriteDebug("Server Queries: Successful: (%d/%d) servers\tFailed: %d servers",
		successcount, len(data.HostsGames), sl.FailedCount)
	return sl, nil
}

// removeBuggedPlayers filters the players to removes "bugged" or stuck players
// from the player list in games like Quake Live where certain servers do not
// correctly send the Steam de-auth message, causing "ghost" or phantom players
// to exist on servers.
func removeBuggedPlayers(players []*models.SteamPlayerInfo) *models.FilteredPlayerInfo {
	rpi := &models.FilteredPlayerInfo{
		FilteredPlayerCount: len(players),
		FilteredPlayers:     players,
	}
	var filtered []*models.SteamPlayerInfo
	// 4 hour threshold with no score; also can leave bots intact in list
	for _, p := range players {
		if p.Score == 0 && int(p.TimeConnectedSecs) > (3600*4) {
			continue
		}
		filtered = append(filtered, p)
	}
	// Empty players (nil) displayed as empty array in JSON, not null
	if len(filtered) == 0 {
		rpi.FilteredPlayerCount = 0
		rpi.FilteredPlayers = make([]*models.SteamPlayerInfo, 0)
	} else {
		rpi.FilteredPlayerCount = len(filtered)
		rpi.FilteredPlayers = filtered
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
