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
	defaultBuggedPlayerTime         = 7
	defaultMaxHostsToReceive        = 4000
	defaultAutoQueryMaster          = false
	defaultTimeBetweenMasterQueries = 90
	// this value is not used in JSON, only in the config dialog
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
	var err error
	prompt := fmt.Sprintf(
		"\nPerform a timed automatic retrieval of game servers from the Steam master server?\nThis is necessary if you want the API to maintain a filterable list of game servers and allow users to query a\nserver by ID, however info can still be queried by address even without this.\nPlease note the relability of retrieving all servers generally decreases as the total number of servers increases.\nAlso note: Valve will throttle your requests if more than 6930 servers are returned per minute.\n>> 'yes' or 'no' [default: %s]: ",
		getBoolString(defaultAutoQueryMaster))

	input := func(r *bufio.Reader) (bool, error) {
		enable, err := r.ReadString('\n')
		if err != nil {
			return defaultAutoQueryMaster,
				fmt.Errorf("Unable to read respone: %s", err)
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
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimedQueryGame(reader *bufio.Reader) string {
	valid := false
	var val string
	var err error
	games := strings.Join(filters.GetGameNames(), ", ")
	prompt := fmt.Sprintf(
		"\nChoose the game you would like to automatically retrieve servers for at timed intervals.\nPossible choices are: %s\nMore games can be added via the %s file.\n>> [default: %s]: ",
		games, constants.GameFileFullPath, "NONE")

	input := func(r *bufio.Reader) (string, error) {
		gameval, err := r.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Unable to read respone: %s", err)
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

	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureMaxServersToRetrieve(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf(
		"\nEnter the maximum number of servers to retrieve from the Steam Master Server at a time.\nThis can be no more than 6930.\n>> [default: %d]: ", defaultMaxHostsToReceive)

	input := func(r *bufio.Reader) (int, error) {
		hostsval, err := r.ReadString('\n')
		if err != nil {
			return defaultMaxHostsToReceive, fmt.Errorf("Unable to read response: %s", err)
		}
		if hostsval == newline {
			return defaultMaxHostsToReceive, nil
		}
		response, err := strconv.Atoi(strings.Trim(hostsval, newline))
		if err != nil {
			return defaultMaxHostsToReceive,
				fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930")
		}
		if response < 500 || response > 6930 {
			return defaultMaxHostsToReceive,
				fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930")
		}
		return response, nil
	}
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimeBetweenQueries(reader *bufio.Reader, game string) int {
	valid := false
	var val int
	var err error
	defaultVal := defaultTimeBetweenMasterQueries
	highServerCountGame := filters.HasHighServerCount(game)
	if highServerCountGame {
		defaultVal = defaultTimeForHighServerCount
	}
	prompt := fmt.Sprintf("\nEnter the time, in seconds, between requests to grab all servers from the master server.\nFor many games this needs to be at least 60. For some games this will need to be even higher.\nNote: if the game returns more than 6930 servers/min, Valve will throttle future requests for 1 min.\n>> [default: %d]: ", defaultVal)

	input := func(r *bufio.Reader) (int, error) {
		timeval, err := r.ReadString('\n')
		if err != nil {
			return defaultVal,
				fmt.Errorf("Unable to read response: %s", err)
		}
		if timeval == newline {
			return defaultVal, nil
		}
		response, err := strconv.Atoi(strings.Trim(timeval, newline))
		if err != nil {
			return defaultVal,
				fmt.Errorf("[ERROR] Time between Steam aster server queries must be at least 60")
		}
		if response < 60 {
			return defaultVal,
				fmt.Errorf("[ERROR] Time between Steam master server queries must be at least 60")
		}
		if highServerCountGame && response < defaultTimeForHighServerCount {
			return defaultVal,
				fmt.Errorf("[ERROR] Game %s typically returns more than 6930 servers so the time between Steam master server queries will need to be at least %d",
					game, defaultTimeForHighServerCount)
		}
		return response, nil
	}
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureSteamBugPlayerTime(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf("\nEnter the time, in hours, before a player is considered \"bugged\" or stuck on a server. This can filter out bots and \"bugged\" non-real players. This is to address the well-known Steam issue where players with very high playing time get \"stuck\" in the player list. This value must be at least 3 hours.\n>> [default: %d]: ", defaultBuggedPlayerTime)

	input := func(r *bufio.Reader) (int, error) {
		timeval, err := r.ReadString('\n')
		if err != nil {
			return defaultBuggedPlayerTime, fmt.Errorf("Unable to read response: %s", err)
		}
		if timeval == newline {
			return defaultBuggedPlayerTime, nil
		}
		response, err := strconv.Atoi(strings.Trim(timeval, newline))
		if err != nil {
			return defaultBuggedPlayerTime, fmt.Errorf("[ERROR] Bugged player time must be at least 3 hours")
		}
		if response < 3 {
			return defaultBuggedPlayerTime, fmt.Errorf("[ERROR] Bugged player time must be at least 3 hours")
		}
		return response, nil
	}
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}
