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

func configureDirectQueries(reader *bufio.Reader, timedEnabled bool) bool {
	valid, val := false, false
	note := ""
	if !timedEnabled {
		note = `
Note: you disabled timed master queries in the previous option, so
if you do not enable this option then there will be no way for users to make
queries (if your server ID database is empty.)`
	}
	prompt := fmt.Sprintf(`
Allow users to directly query *any* IP address, not just those in the serverID
database? This may have some issues for some games and it also may have abuse
implications since your server could query unknown, user-specified IP addresses.%s
%s`, note, promptColor("> 'yes' or 'no' [default: %s]: ",
		getBoolString(defaultAllowDirectUserQueries)))

	input := func(r *bufio.Reader) (bool, error) {
		enable, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultAllowDirectUserQueries,
				fmt.Errorf("Unable to read respone: %s", rserr)
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

func configureMaxHostsPerAPIQuery(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the maximum number of servers that users may query at a time via the API.
%s`, promptColor("> [default: %d]: ", defaultMaxHostsPerAPIQuery))

	input := func(r *bufio.Reader) (int, error) {
		hostsval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultMaxHostsPerAPIQuery, fmt.Errorf("Unable to read response: %s",
				rserr)
		}
		if hostsval == newline {
			return defaultMaxHostsPerAPIQuery, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(hostsval, newline))
		if rserr != nil {
			return defaultMaxHostsPerAPIQuery,
				fmt.Errorf("[ERROR] Maximum hosts to allow per API query must be a positive number")
		}
		if response <= 0 {
			return defaultMaxHostsPerAPIQuery,
				fmt.Errorf("[ERROR] Maximum hosts to allow per API query must be a positive number")
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

func configureWebServerPort(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the port number on which the API web server will listen.
%s`, promptColor("> [default: %d]: ", defaultAPIWebPort))

	input := func(r *bufio.Reader) (int, error) {
		portval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultAPIWebPort, fmt.Errorf("Unable to read response: %s", rserr)
		}
		if portval == newline {
			return defaultAPIWebPort, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(portval, newline))
		if rserr != nil {
			return defaultAPIWebPort,
				fmt.Errorf("[ERROR] API webserver port must be between 1 and 65535")
		}
		if response < 1 || response > 65535 {
			return defaultAPIWebPort,
				fmt.Errorf("[ERROR] API webserver port must be between 1 and 65535")
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

func configureWebTimeout(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the time in seconds before an HTTP request times out. This must be at
least 5 seconds; don't set this too low or the response will not be returned
to the user.
%s`, promptColor("> [default: %d]: ", defaultAPIWebTimeout))

	input := func(r *bufio.Reader) (int, error) {
		timeoutval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultAPIWebTimeout, fmt.Errorf("Unable to read response: %s",
				rserr)
		}
		if timeoutval == newline {
			return defaultAPIWebTimeout, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(timeoutval, newline))
		if rserr != nil {
			return defaultAPIWebTimeout,
				fmt.Errorf("[ERROR] API timeout cannot be less than 5 seconds")
		}
		if response < 5 {
			return defaultAPIWebTimeout,
				fmt.Errorf("[ERROR] API timeout cannot be less than 5 seconds")
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
