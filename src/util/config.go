package util

// config.go - configuration operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var (
	newline        = getNewLineForOS()
	ConfigFullPath = path.Join(ConfigDirectory, ConfigFileName)
)

const (
	defaultEnableDebugMessages      = false
	defaultEnableAppLogging         = true
	defaultEnableSteamLogging       = false
	defaultEnableWebLogging         = true
	defaultMaxLogSize               = 5120
	defaultMaxLogCount              = 6
	defaultMaxHostsToReceive        = 3800
	defaultMaxHostsPerAPIQuery      = 12
	defaultTimeBetweenMasterQueries = 60
	defaultAPIWebPort               = 40080
	ConfigDirectory                 = "conf"
	ConfigFileName                  = "config.conf"
	Version                         = "0.1"
)

type Config struct {
	EnableDebugMessages      bool  `json:"debugMessages"`
	EnableAppLogging         bool  `json:"enableAppLogging"`
	EnableSteamLogging       bool  `json:"enableSteamLogging"`
	EnableWebLogging         bool  `json:"enableWebLogging"`
	MaximumLogSize           int64 `json:"maxLogFilesize"`
	MaximumLogCount          int   `json:"maxLogCount"`
	MaximumHostsPerAPIQuery  int   `json:"maxHostsPerAPIQuery"`
	MaximumHostsToReceive    int   `json:"maxHostsToReceive"`
	TimeBetweenMasterQueries int   `json:"timeBetweenMasterQueries"`
	APIWebPort               int   `json:"apiWebPort"`
}

func getNewLineForOS() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func getLoggingValue(r *bufio.Reader, lt logType) (bool, error) {
	enable, err := r.ReadString('\n')
	var defaultVal bool
	if lt == App {
		defaultVal = defaultEnableAppLogging
	} else if lt == Steam {
		defaultVal = defaultEnableSteamLogging
	} else if lt == Web {
		defaultVal = defaultEnableWebLogging
	}

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

func getMaxLogSizeValue(r *bufio.Reader) (int64, error) {
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

func getMaxLogCountValue(r *bufio.Reader) (int, error) {
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

func getMaxHostsPerAPIQueryValue(r *bufio.Reader) (int, error) {
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
			fmt.Errorf("[ERROR] Maximum hosts to allow per API query must be a positive number.")
	}
	if response <= 0 {
		return defaultMaxHostsPerAPIQuery,
			fmt.Errorf("[ERROR] Maximum hosts to allow per API query must be a positive number.")
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

func getAPIWebPortValue(r *bufio.Reader) (int, error) {
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

func createConfigDir() error {
	if DirExists(ConfigDirectory) {
		return nil
	}
	if err := os.Mkdir(ConfigDirectory, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func writeConfig(cfg *Config) error {
	if err := createConfigDir(); err != nil {
		fmt.Printf("Couldn't create '%s' dir: %s\n", ConfigDirectory, err)
		return err
	}
	f, err := os.Create(ConfigFullPath)
	if err != nil {
		fmt.Printf("Couldn't create '%s' file: %s\n", ConfigFullPath, err)
		return err
	}
	defer f.Close()
	c, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling JSON for '%s' file: %s\n", ConfigFullPath, err)
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.Write(c)
	if err != nil {
		fmt.Printf("Error writing contents of config file to disk: %s\n", err)
		return err
	}
	f.Sync()
	w.Flush()
	fmt.Printf("Successfuly wrote config file to %s\n", ConfigFullPath)
	return nil
}

func ReadConfig() (*Config, error) {
	f, err := os.Open(ConfigFullPath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	d := json.NewDecoder(r)
	cfg := &Config{}
	if err := d.Decode(cfg); err != nil {
		fmt.Printf("Error decoding config file: %s\n", err)
		return nil, err
	}
	return cfg, nil
}

func getBoolLogString(lt logType) string {
	var val string
	switch lt {
	case App:
		if defaultEnableAppLogging {
			val = "yes"
		}
		val = "no"
	case Steam:
		if defaultEnableSteamLogging {
			val = "yes"
		}
		val = "no"
	case Web:
		if defaultEnableWebLogging {
			val = "yes"
		}
		val = "no"
	}
	return val
}

func CreateConfig() error {
	reader := bufio.NewReader(os.Stdin)
	cfg := &Config{}
	fmt.Printf("steamtest v%s - configuration file creation\n", Version)
	fmt.Print(
		"Type a value and press 'ENTER'. Leave a value empty and press 'ENTER' to use the default value.\n\n")

	// Configure app-related logging
	validLogAppEnableVal := false
	var logAppEnableVal bool
	for !validLogAppEnableVal {
		var err error
		fmt.Printf(
			"\nLog application-related info and error messages to disk?\n>> 'yes' or 'no' [default: %s]: ",
			getBoolLogString(App))
		logAppEnableVal, err = getLoggingValue(reader, App)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.EnableAppLogging = logAppEnableVal
			validLogAppEnableVal = true
		}
	}
	// Configure Steam-related logging
	validLogSteamEnableVal := false
	var logSteamEnableVal bool
	for !validLogSteamEnableVal {
		var err error
		fmt.Printf(
			"\nLog Steam connection info and error messages to disk?\nNOTE: this can dramatically increase resource usage and should only be used for debugging.\n>> 'yes' or 'no' [default: %s]: ",
			getBoolLogString(Steam))
		logSteamEnableVal, err = getLoggingValue(reader, Steam)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.EnableWebLogging = logSteamEnableVal
			validLogSteamEnableVal = true
		}
	}
	// Configure Web-related logging
	validLogWebEnableVal := false
	var logWebEnableVal bool
	for !validLogWebEnableVal {
		var err error
		fmt.Printf(
			"\nShould API web-related info and error messages should be logged to disk?\n>> 'yes' or 'no' [default: %s]: ",
			getBoolLogString(Web))
		logWebEnableVal, err = getLoggingValue(reader, Web)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.EnableWebLogging = logWebEnableVal
			validLogWebEnableVal = true
		}
	}
	// Configure max log size and max log count
	if logAppEnableVal || logSteamEnableVal || logWebEnableVal {
		validLogMaxSizeVal := false
		for !validLogMaxSizeVal {
			fmt.Printf(
				"\nEnter the maximum file size for log files in Kilobytes.\n>> By default this is %d, or %d megabyte(s). [default: %d]: ",
				defaultMaxLogSize, defaultMaxLogSize/1024, defaultMaxLogSize)
			logMaxSizeVal, err := getMaxLogSizeValue(reader)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumLogSize = logMaxSizeVal
				validLogMaxSizeVal = true
			}
		}
		validLogMaxCountVal := false
		for !validLogMaxCountVal {
			fmt.Printf(
				"\nEnter the maximum number of log files to keep.\n>> [default: %d]: ",
				defaultMaxLogCount)
			logMaxCountVal, err := getMaxLogCountValue(reader)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumLogCount = logMaxCountVal
				validLogMaxCountVal = true
			}
		}
	}
	// Configure maximum # of servers to retrieve
	validMaxMasterHostsVal := false
	for !validMaxMasterHostsVal {
		fmt.Printf(
			"\nEnter the maximum number of servers to retrieve from the Steam Master Server at a time.\nThis can be no more than 6930.\n>> [default: %d]: ", defaultMaxHostsToReceive)
		maxHostVal, err := getMaxMasterServerHostsValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.MaximumHostsToReceive = maxHostVal
			validMaxMasterHostsVal = true
		}
	}
	// Configure maximum number of servers to allow users to query via API
	validMaxAPIHostsVal := false
	for !validMaxAPIHostsVal {
		fmt.Printf(
			"\nEnter the maximum number of servers that users may query at a time via the API.\n>> [default: %d]: ", defaultMaxHostsPerAPIQuery)
		maxHostVal, err := getMaxHostsPerAPIQueryValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.MaximumHostsPerAPIQuery = maxHostVal
			validMaxAPIHostsVal = true
		}
	}
	// Configure time between master server queries
	validTimeBetweenVal := false
	for !validTimeBetweenVal {
		fmt.Printf("\nEnter the time, in seconds, between Master Server queries.\nThis must be greater than 20 & should not be too low as receiving servers can take a while.\n>> [default: %d]: ",
			defaultTimeBetweenMasterQueries)
		timeBetweenVal, err := getQueryTimeValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.TimeBetweenMasterQueries = timeBetweenVal
			validTimeBetweenVal = true
		}
	}
	// Configure webserver port
	validAPIPortVal := false
	for !validAPIPortVal {
		fmt.Printf("\nEnter the port number on which the API web server will listen.\n>> [default: %d]: ",
			defaultAPIWebPort)
		apiPortVal, err := getAPIWebPortValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.APIWebPort = apiPortVal
			validAPIPortVal = true
		}
	}
	// Debug mode for testing (no user option to enable)
	cfg.EnableDebugMessages = defaultEnableDebugMessages
	if err := writeConfig(cfg); err != nil {
		return err
	}
	return nil
}
