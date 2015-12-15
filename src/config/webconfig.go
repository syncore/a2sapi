package config

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultAllowDirectUserQueries = false
	defaultMaxHostsPerAPIQuery    = 12
	defaultAPIWebTimeout          = 7
	defaultAPIWebPort             = 40080
)

// CfgWeb represents web-related API configuration options.
type CfgWeb struct {
	AllowDirectUserQueries  bool `json:"allowDirectUserQueries"`
	APIWebPort              int  `json:"apiWebPort"`
	APIWebTimeout           int  `json:"apiWebTimeout"`
	MaximumHostsPerAPIQuery int  `json:"maxHostsPerAPIQuery"`
}

func configureDirectQueries(reader *bufio.Reader) bool {
	valid := false
	var val bool
	var err error
	prompt := fmt.Sprintf(
		"\nAllow users to directly query *any* IP address, not just those in the serverID database?\nThis is mainly for testing and has some issues depending on the game.\nIt also may have security implications so enable with caution.\nNote: if you have timed master queries disabled & your server ID database is empty, then without this option there will be no way for users to make queries.\n>> 'yes' or 'no' [default: %s]: ",
		getBoolString(defaultAllowDirectUserQueries))

	input := func(r *bufio.Reader) (bool, error) {
		enable, err := r.ReadString('\n')
		if err != nil {
			return defaultAllowDirectUserQueries,
				fmt.Errorf("Unable to read respone: %s", err)
		}
		if enable == newline {
			return defaultAllowDirectUserQueries, nil
		}
		response := strings.Trim(enable, newline)
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			return true, nil
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response,
			"no") {
			return false, nil
		} else {
			return defaultAllowDirectUserQueries,
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

func configureMaxHostsPerAPIQuery(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf(
		"\nEnter the maximum number of servers that users may query at a time via the API.\n>> [default: %d]: ", defaultMaxHostsPerAPIQuery)

	input := func(r *bufio.Reader) (int, error) {
		hostsval, err := r.ReadString('\n')
		if err != nil {
			return defaultMaxHostsPerAPIQuery, fmt.Errorf("Unable to read response: %s", err)
		}
		if hostsval == newline {
			return defaultMaxHostsPerAPIQuery, nil
		}
		response, err := strconv.Atoi(strings.Trim(hostsval, newline))
		if err != nil {
			return defaultMaxHostsPerAPIQuery,
				fmt.Errorf("[ERROR] Maximum hosts to allow per API query must be a positive number")
		}
		if response <= 0 {
			return defaultMaxHostsPerAPIQuery,
				fmt.Errorf("[ERROR] Maximum hosts to allow per API query must be a positive number")
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

func configureWebServerPort(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf(
		"\nEnter the port number on which the API web server will listen.\n>> [default: %d]: ",
		defaultAPIWebPort)

	input := func(r *bufio.Reader) (int, error) {
		portval, err := r.ReadString('\n')
		if err != nil {
			return defaultAPIWebPort, fmt.Errorf("Unable to read response: %s", err)
		}
		if portval == newline {
			return defaultAPIWebPort, nil
		}
		response, err := strconv.Atoi(strings.Trim(portval, newline))
		if err != nil {
			return defaultAPIWebPort,
				fmt.Errorf("[ERROR] API webserver port must be between 1 and 65535")
		}
		if response < 1 || response > 65535 {
			return defaultAPIWebPort,
				fmt.Errorf("[ERROR] API webserver port must be between 1 and 65535")
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

func configureWebTimeout(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf("\nEnter the time in seconds before an HTTP request times out.\nThis must be at least 5 seconds; don't set this too low or the response will not be returned to the user.\n>> [default: %d]: ",
		defaultAPIWebTimeout)

	input := func(r *bufio.Reader) (int, error) {
		timeoutval, err := r.ReadString('\n')
		if err != nil {
			return defaultAPIWebTimeout, fmt.Errorf("Unable to read response: %s", err)
		}
		if timeoutval == newline {
			return defaultAPIWebTimeout, nil
		}
		response, err := strconv.Atoi(strings.Trim(timeoutval, newline))
		if err != nil {
			return defaultAPIWebTimeout,
				fmt.Errorf("[ERROR] API timeout cannot be less than 5 seconds")
		}
		if response < 5 {
			return defaultAPIWebTimeout,
				fmt.Errorf("[ERROR] API timeout cannot be less than 5 seconds")
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
