package config

import (
	"bufio"
	"fmt"
	"steamtest/src/constants"
	"steamtest/src/steam/filters"
	"strconv"
	"strings"
)

const (
	defaultBuggedPlayerTime         = 6
	defaultMaxHostsToReceive        = 4000
	defaultAutoQueryMaster          = false
	defaultTimeBetweenMasterQueries = 90
	// defaultTimeForHighServerCount: not used in JSON, only in the config dialog
	defaultTimeForHighServerCount = 120
)

// CfgSteam represents Steam-related configuration options.
type CfgSteam struct {
	SteamBugPlayerTime       int    `json:"steamBugPlayerHours"`
	AutoQueryMaster          bool   `json:"timedMasterServerQuery"`
	AutoQueryGame            string `json:"gameForTimedMasterQuery"`
	TimeBetweenMasterQueries int    `json:"timeBetweenMasterQueries"`
	MaximumHostsToReceive    int    `json:"maxHostsToReceive"`
}

func configureTimedMasterQuery(reader *bufio.Reader) bool {
	valid := false
	var val bool
	prompt := fmt.Sprintf(`
Perform a timed automatic retrieval of game servers from the Steam
master server? This is necessary if you want the API to maintain a
filterable list of game servers and allow users to query a server
by ID, however info can still be queried by address even without this.
Please note the reliability of retrieving all servers generally decreases
as the total number of servers increases. Also note: Valve will throttle
your requests if more than 6930 servers are returned per minute.
%s`, promptColor("> 'yes' or 'no' [default: %s]: ",
		getBoolString(defaultAutoQueryMaster)))

	input := func(r *bufio.Reader) (bool, error) {
		enable, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultAutoQueryMaster,
				fmt.Errorf("Unable to read respone: %s", rserr)
		}
		if enable == newline {
			return defaultAutoQueryMaster, nil
		}
		response := strings.Trim(enable, newline)
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			return true, nil
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response,
			"no") {
			return false, nil
		} else {
			return defaultAutoQueryMaster,
				fmt.Errorf("[ERROR] Invalid response. Valid responses: y, yes, n, no")
		}
	}
	var err error
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimedQueryGame(reader *bufio.Reader) string {
	valid := false
	var val string
	games := strings.Join(filters.GetGameNames(), ", ")
	prompt := fmt.Sprintf(`
Choose the game you would like to automatically retrieve servers for at timed
intervals. Possible choices are: %s
More games can be added via the %s file.
%s`, games, constants.GameFileFullPath, promptColor("> [default: NONE]: "))

	input := func(r *bufio.Reader) (string, error) {
		gameval, rserr := r.ReadString('\n')
		if rserr != nil {
			return "", fmt.Errorf("Unable to read respone: %s", rserr)
		}
		if gameval == newline {
			return "", fmt.Errorf("[ERROR] Invalid response. Valid responses: %s", games)
		}
		response := strings.Trim(gameval, newline)
		if filters.IsValidGame(response) {
			// format the capitalization
			return filters.GetGameByName(response).Name, nil
		}
		return "", fmt.Errorf("[ERROR] Invalid response. Valid responses: %s", games)
	}
	var err error
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureMaxServersToRetrieve(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the maximum number of servers to retrieve from the Steam Master Server at
a time. This can be no more than 6930.
%s`, promptColor("> [default: %d]: ", defaultMaxHostsToReceive))

	input := func(r *bufio.Reader) (int, error) {
		hostsval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultMaxHostsToReceive, fmt.Errorf("Unable to read response: %s", rserr)
		}
		if hostsval == newline {
			return defaultMaxHostsToReceive, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(hostsval, newline))
		if rserr != nil {
			return defaultMaxHostsToReceive,
				fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930")
		}
		if response < 500 || response > 6930 {
			return defaultMaxHostsToReceive,
				fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930")
		}
		return response, nil
	}
	var err error
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimeBetweenQueries(reader *bufio.Reader, game string) int {
	valid := false
	var val int
	defaultVal := defaultTimeBetweenMasterQueries
	highServerCountGame := filters.HasHighServerCount(game)
	if highServerCountGame {
		defaultVal = defaultTimeForHighServerCount
	}
	prompt := fmt.Sprintf(`
Enter the time, in seconds, between requests to grab all servers from the master
server. For many games this needs to be at least 60. For some games this will
need to be even higher. Note: if the game returns more than 6930 servers/min,
Valve will throttle future requests for 1 min.
%s `, promptColor("> [default: %d]: ", defaultVal))

	input := func(r *bufio.Reader) (int, error) {
		timeval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultVal,
				fmt.Errorf("Unable to read response: %s", rserr)
		}
		if timeval == newline {
			return defaultVal, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(timeval, newline))
		if rserr != nil {
			return defaultVal,
				fmt.Errorf("[ERROR] Time between Steam aster server queries must be at least 60")
		}
		if response < 60 {
			return defaultVal,
				fmt.Errorf("[ERROR] Time between Steam master server queries must be at least 60")
		}
		if highServerCountGame && response < defaultTimeForHighServerCount {
			return defaultVal, fmt.Errorf(`
[ERROR] Game %s typically returns more than 6930 servers so the time between
Steam master server queries will need to be at least %d`, game,
				defaultTimeForHighServerCount)
		}
		return response, nil
	}
	var err error
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureSteamBugPlayerTime(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the time, in hours, before a player is considered "bugged" or stuck on a
server. This can filter out bots and "bugged" non-real players. This is to
address the well-known issue in certain games (i.e. Quake Live) where game servers
do not receive the Steam de-auth message which causes players to get "stuck" in
the player list, long after they've disconnected. This value must be at least 3 hours.
%s`, promptColor("> [default: %d]: ", defaultBuggedPlayerTime))

	input := func(r *bufio.Reader) (int, error) {
		timeval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultBuggedPlayerTime, fmt.Errorf("Unable to read response: %s",
				rserr)
		}
		if timeval == newline {
			return defaultBuggedPlayerTime, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(timeval, newline))
		if rserr != nil {
			return defaultBuggedPlayerTime,
				fmt.Errorf("[ERROR] Bugged player time must be at least 3 hours")
		}
		if response < 3 {
			return defaultBuggedPlayerTime,
				fmt.Errorf("[ERROR] Bugged player time must be at least 3 hours")
		}
		return response, nil
	}
	var err error
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}
