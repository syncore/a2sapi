package config

import (
	"a2sapi/src/constants"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultEnableAppLogging   = false
	defaultEnableSteamLogging = false
	defaultEnableWebLogging   = false
	defaultMaxLogSize         = 5120
	defaultMaxLogCount        = 5
)

// CfgLog represents logging-related configuration options.
type CfgLog struct {
	EnableAppLogging   bool  `json:"enableAppLogging"`
	EnableSteamLogging bool  `json:"enableSteamLogging"`
	EnableWebLogging   bool  `json:"enableWebLogging"`
	MaximumLogSize     int64 `json:"maxLogFilesize"`
	MaximumLogCount    int   `json:"maxLogCount"`
}

func configureLoggingEnable(reader *bufio.Reader, logt constants.LogType) bool {
	valid, val, defaultVal := false, false, false
	var prompt string

	switch logt {
	case constants.LTypeApp:
		defaultVal = defaultEnableAppLogging
		prompt = fmt.Sprintf(`
Log application-related info and error messages to disk?
%s`, promptColor("> 'yes' or 'no' [default: %s]: ",
			getBoolString(defaultEnableAppLogging)))

	case constants.LTypeSteam:
		defaultVal = defaultEnableSteamLogging
		prompt = fmt.Sprintf(`
Log Steam connection info and error messages to disk?
NOTE: this can dramatically increase resource usage and should only be
used for debugging.
%s`, promptColor("> 'yes' or 'no' [default: %s]: ",
			getBoolString(defaultEnableSteamLogging)))

	case constants.LTypeWeb:
		defaultVal = defaultEnableWebLogging
		prompt = fmt.Sprintf(`
Should API web-related info and error messages be
logged to disk?
%s`, promptColor("> 'yes' or 'no' [default: %s]: ",
			getBoolString(defaultEnableWebLogging)))
	}
	input := func(r *bufio.Reader, lt constants.LogType) (bool, error) {
		enable, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultVal, fmt.Errorf("Unable to read respone: %s", rserr)
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
				fmt.Errorf("[ERROR] Invalid response. Valid responses: y, yes, n, no")
		}
	}
	var err error
	for !valid {
		fmt.Print(prompt)
		val, err = input(reader, logt)
		if err != nil {
			errorColor(err)
		} else {
			valid = true
		}
	}
	return val
}

func configureMaxLogSize(reader *bufio.Reader) int64 {
	valid := false
	var val int64
	prompt := fmt.Sprintf(`
Enter the maximum file size for log files in Kilobytes.
By default this is %d, or %d megabyte(s).
%s`, defaultMaxLogSize, defaultMaxLogSize/1024,
		promptColor("> [default: %d]: ", defaultMaxLogSize))

	input := func(r *bufio.Reader) (int64, error) {
		sizeval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultMaxLogSize, fmt.Errorf("Unable to read response: %s", rserr)
		}
		if sizeval == newline {
			return defaultMaxLogSize, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(sizeval, newline))
		if rserr != nil {
			return defaultMaxLogSize,
				fmt.Errorf("[ERROR] Maximum log file size must be a positive number")
		}
		if response <= 0 {
			return defaultMaxLogSize,
				fmt.Errorf("[ERROR] Maximum log file size must be a positive number")
		}
		return int64(response), nil
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

func configureMaxLogCount(reader *bufio.Reader) int {
	valid := false
	var val int
	prompt := fmt.Sprintf(`
Enter the maximum number of log files to keep.
%s`, promptColor("> [default: %d]: ", defaultMaxLogCount))

	input := func(r *bufio.Reader) (int, error) {
		sizeval, rserr := r.ReadString('\n')
		if rserr != nil {
			return defaultMaxLogCount, fmt.Errorf("Unable to read response: %s", rserr)
		}
		if sizeval == newline {
			return defaultMaxLogCount, nil
		}
		response, rserr := strconv.Atoi(strings.Trim(sizeval, newline))
		if rserr != nil {
			return defaultMaxLogCount,
				fmt.Errorf("[ERROR] Maximum log count must be a positive number")
		}
		if response <= 0 {
			return defaultMaxLogCount,
				fmt.Errorf("[ERROR] Maximum log count must be a positive number")
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
