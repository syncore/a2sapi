package util

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultMaxHostsToReceive        = 3800
	defaultTimeBetweenMasterQueries = 60
)

// CfgSteam represents Steam-related configuration options.
type CfgSteam struct {
	MaximumHostsToReceive    int `json:"maxHostsToReceive"`
	TimeBetweenMasterQueries int `json:"timeBetweenMasterQueries"`
}

func getMaxMasterServerHostsValue(r *bufio.Reader) (int, error) {
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

func getQueryTimeValue(r *bufio.Reader) (int, error) {
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

func configureMaxServersToRetrieve(r *bufio.Reader) int {
	valid := false
	var val int
	var err error
	for !valid {
		fmt.Printf(
			"\nEnter the maximum number of servers to retrieve from the Steam Master Server at a time.\nThis can be no more than 6930.\n>> [default: %d]: ", defaultMaxHostsToReceive)
		val, err = getMaxMasterServerHostsValue(r)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureTimeBetweenQueries(r *bufio.Reader) int {
	valid := false
	var val int
	var err error
	for !valid {
		fmt.Printf("\nEnter the time, in seconds, between Master Server queries.\nThis must be greater than 20 & should not be too low as receiving servers can take a while.\n>> [default: %d]: ",
			defaultTimeBetweenMasterQueries)
		val, err = getQueryTimeValue(r)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}
