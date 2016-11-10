package steam

// steammasterweb.go - valve web "master" server list
// Note: this is now intended to be the method of retrieving all of the servers because
// Valve apparently shut down their master servers on November 8, 2016. Previously, a call
// was made via steammaster.go to the master server (i.e. hl2master.steampowered.com) - now we
// use this undocumented web endpoint to retrieve the servers, since the master server method of
//  retrieval is no longer available, and it is uncertain whether it will ever be available again
// in the future.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/syncore/a2sapi/src/config"
	"github.com/syncore/a2sapi/src/logger"
	"github.com/syncore/a2sapi/src/steam/filters"
)

var steamWebAPIURL = func(webAPIKey, filter string, limit int) string {
	return fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetServerList/v1/?key=%s&format=json&filter=%s&limit=%d",
		webAPIKey, filter, limit)
}

// webGameServerList repersents the response returned from the Steam Web API that includes the
// server addresses (and some extra information that we are not interested in)
type webGameServerList struct {
	Response struct {
		Servers []struct {
			Addr       string `json:"addr"`
			Gameport   int    `json:"gameport"`
			Steamid    string `json:"steamid"`
			Name       string `json:"name"`
			Appid      int    `json:"appid"`
			Gamedir    string `json:"gamedir"`
			Version    string `json:"version"`
			Product    string `json:"product"`
			Region     int    `json:"region"`
			Players    int    `json:"players"`
			MaxPlayers int    `json:"maxPlayers"`
			Bots       int    `json:"bots"`
			Map        string `json:"map"`
			Secure     bool   `json:"secure"`
			Dedicated  bool   `json:"dedicated"`
			Os         string `json:"os"`
			Gametype   string `json:"gametype"`
		} `json:"servers"`
	} `json:"response"`
}

func getServersWeb(filter filters.Filter) ([]string, error) {
	var fsl []string
	for _, f := range filter.Filters {
		fsl = append(fsl, string(f))
	}
	filterStr := strings.Join(fsl, "")
	response, err := http.Get(steamWebAPIURL(config.Config.SteamConfig.SteamWebAPIKey, filterStr,
		config.Config.SteamConfig.MaximumHostsToReceive))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var webAPIResponseModel webGameServerList
	var servers []string
	apiResult := json.NewDecoder(response.Body)
	if err := apiResult.Decode(&webAPIResponseModel); err != nil {
		logger.WriteDebug("Error decoding Steam Web API response: %s", err)
		return nil, err
	}
	for _, server := range webAPIResponseModel.Response.Servers {
		servers = append(servers, server.Addr)
	}
	return servers, nil
}

// NewMasterWebQuery initiates a new Steam "Master" server query using the Steam Web API for a
// given filter, returning a MasterQuery struct containing the hosts retrieved in the event of
// success or an empty struct and an error in the event of failure.
func NewMasterWebQuery(filter filters.Filter) (MasterQuery, error) {
	sl, err := getServersWeb(filter)
	if err != nil {
		return MasterQuery{}, err
	}
	logger.LogSteamInfo("*** Retrieved %d %s servers.", len(sl), filter.Game.Name)

	return MasterQuery{Servers: sl}, nil
}
