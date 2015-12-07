package util

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultEnableDebugMessages = false
	defaultEnableAppLogging    = true
	defaultEnableSteamLogging  = false
	defaultEnableWebLogging    = true
	defaultMaxLogSize          = 5120
	defaultMaxLogCount         = 6
)

// CfgLog represents logging-related configuration options.
type CfgLog struct {
	EnableDebugMessages bool  `json:"debugMessages"`
	EnableAppLogging    bool  `json:"enableAppLogging"`
	EnableSteamLogging  bool  `json:"enableSteamLogging"`
	EnableWebLogging    bool  `json:"enableWebLogging"`
	MaximumLogSize      int64 `json:"maxLogFilesize"`
	MaximumLogCount     int   `json:"maxLogCount"`
}

func configureLoggingEnable(reader *bufio.Reader, logt logType) bool {
	valid := false
	var val bool
	var err error
	var prompt string
	var defaultVal bool
	switch logt {
	case App:
		defaultVal = defaultEnableAppLogging
		prompt = fmt.Sprintf(
			"\nLog application-related info and error messages to disk?\n>> 'yes' or 'no' [default: %s]: ",
			getBoolString(defaultEnableAppLogging))
	case Steam:
		defaultVal = defaultEnableSteamLogging
		prompt = fmt.Sprintf(
			"\nLog Steam connection info and error messages to disk?\nNOTE: this can dramatically increase resource usage and should only be used for debugging.\n>> 'yes' or 'no' [default: %s]: ",
			getBoolString(defaultEnableSteamLogging))
	case Web:
		defaultVal = defaultEnableWebLogging
		prompt = fmt.Sprintf(
			"\nShould API web-related info and error messages be logged to disk?\n>> 'yes' or 'no' [default: %s]: ",
			getBoolString(defaultEnableWebLogging))
	}
	input := func(r *bufio.Reader, lt logType) (bool, error) {
		enable, err := r.ReadString('\n')
		if err != nil {
			return defaultVal, fmt.Errorf("Unable to read respone: %s", err)
		}
		if enable == newline {
			return defaultVal, nil
		}
		response := strings.Trim(enable, newline)
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			return true, nil
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response,
			"no") {
			return false, nil
		} else {
			return defaultVal,
				fmt.Errorf("Invalid response. Valid responses: y, yes, n, no")
		}
	}
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader, logt)
		if err != nil {
			fmt.Println(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureMaxLogSize(reader *bufio.Reader) int64 {
	valid := false
	var val int64
	var err error
	prompt := fmt.Sprintf(
		"\nEnter the maximum file size for log files in Kilobytes.\n>> By default this is %d, or %d megabyte(s). [default: %d]: ",
		defaultMaxLogSize, defaultMaxLogSize/1024, defaultMaxLogSize)

	input := func(r *bufio.Reader) (int64, error) {
		sizeval, err := r.ReadString('\n')
		if err != nil {
			return defaultMaxLogSize, fmt.Errorf("Unable to read response: %s", err)
		}
		if sizeval == newline {
			return defaultMaxLogSize, nil
		}
		response, err := strconv.Atoi(strings.Trim(sizeval, newline))
		if err != nil {
			return defaultMaxLogSize,
				fmt.Errorf("[ERROR] Maximum log file size must be a positive number")
		}
		if response <= 0 {
			return defaultMaxLogSize,
				fmt.Errorf("[ERROR] Maximum log file size must be a positive number")
		}
		return int64(response), nil
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

func configureMaxLogCount(reader *bufio.Reader) int {
	valid := false
	var val int
	var err error
	prompt := fmt.Sprintf(
		"\nEnter the maximum number of log files to keep.\n>> [default: %d]: ",
		defaultMaxLogCount)

	input := func(r *bufio.Reader) (int, error) {
		sizeval, err := r.ReadString('\n')
		if err != nil {
			return defaultMaxLogCount, fmt.Errorf("Unable to read response: %s", err)
		}
		if sizeval == newline {
			return defaultMaxLogCount, nil
		}
		response, err := strconv.Atoi(strings.Trim(sizeval, newline))
		if err != nil {
			return defaultMaxLogCount,
				fmt.Errorf("[ERROR] Maximum log count must be a positive number")
		}
		if response <= 0 {
			return defaultMaxLogCount,
				fmt.Errorf("[ERROR] Maximum log count must be a positive number")
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
