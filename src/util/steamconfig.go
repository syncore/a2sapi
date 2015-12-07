package util

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultBuggedPlayerTime         = 7
	defaultMaxHostsToReceive        = 3800
	defaultTimeBetweenMasterQueries = 60
)

// CfgSteam represents Steam-related configuration options.
type CfgSteam struct {
	SteamBugPlayerTime       int `json:"steamBugPlayerHours"`
	MaximumHostsToReceive    int `json:"maxHostsToReceive"`
	TimeBetweenMasterQueries int `json:"timeBetweenMasterQueries"`
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

func configureTimeBetweenQueries(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf("\nEnter the time, in seconds, between Master Server queries.\nThis must be greater than 20 & should not be too low as receiving servers can take a while.\n>> [default: %d]: ",
		defaultTimeBetweenMasterQueries)

	input := func(r *bufio.Reader) (int, error) {
		timeval, err := r.ReadString('\n')
		if err != nil {
			return defaultTimeBetweenMasterQueries,
				fmt.Errorf("Unable to read response: %s", err)
		}
		if timeval == newline {
			return defaultTimeBetweenMasterQueries, nil
		}
		response, err := strconv.Atoi(strings.Trim(timeval, newline))
		if err != nil {
			return defaultTimeBetweenMasterQueries,
				fmt.Errorf("[ERROR] Time between Steam master server queries must be a number greater than 20")
		}
		if response < 20 {
			return defaultTimeBetweenMasterQueries,
				fmt.Errorf("[ERROR] Time between Steam master server queries must be a number greater than 20")

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
