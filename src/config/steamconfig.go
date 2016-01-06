package config

import (
	"a2sapi/src/constants"
	"a2sapi/src/steam/filters"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultMaxHostsToReceive        = 4000
	defaultAutoQueryMaster          = true
	defaultTimeBetweenMasterQueries = 90
	// defaultTimeForHighServerCount: not used in JSON, only in the config dialog
	defaultTimeForHighServerCount = 120
)

// CfgSteam represents Steam-related configuration options.
type CfgSteam struct {
	AutoQueryMaster          bool   `json:"timedMasterServerQuery"`
	AutoQueryGame            string `json:"gameForTimedMasterQuery"`
	TimeBetweenMasterQueries int    `json:"timeBetweenMasterQueries"`
	MaximumHostsToReceive    int    `json:"maxHostsToReceive"`
}

func configureTimedMasterQuery(reader *bufio.Reader) bool {
	valid, val := false, false
	prompt := fmt.Sprintf(`
Perform a timed automatic retrieval of game servers from the Steam
master server? This is necessary if you want the API to maintain a
filterable list of game servers and allow users to query a server
by ID, however info can still be queried by address even without this, if
you enable it in the next option. Please note the reliability of retrieving
all servers generally decreases as the total number of servers increases.
Also note: Valve will throttle your requests if more than 6930 servers are
returned per minute.
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
