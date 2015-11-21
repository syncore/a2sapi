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
	defaultEnableAppLogging         = true
	defaultEnableWebLogging         = true
	defaultMaxAppLogSize            = 1024
	defaultMaxAppLogCount           = 5
	defaultMaxWebLogSize            = 1024
	defaultMaxWebLogCount           = 5
	defaultMaxHostsToReceive        = 3800
	defaultTimeBetweenMasterQueries = 60
	defaultAPIWebPort               = 40080
	ConfigDirectory                 = "conf"
	ConfigFileName                  = "config.conf"
	Version                         = "0.1"
)

type Config struct {
	EnableAppLogging         bool  `json:"enableAppLogging"`
	EnableWebLogging         bool  `json:"enableWebLogging"`
	MaximumAppLogSize        int64 `json:"maxAppLogFilesize"`
	MaximumWebLogSize        int64 `json:"maxWebLogFilesize"`
	MaximumAppLogCount       int   `json:"maxAppLogCount"`
	MaximumWebLogCount       int   `json:"maxWebLogCount"`
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
	} else {
		defaultVal = defaultEnableWebLogging
	}
	if err != nil {
		return defaultVal, fmt.Errorf("Unable to read respone: %s\n", err)
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
			fmt.Errorf("Invalid response. Valid responses: y, yes, n, no\n")
	}
}

func getMaxLogSizeValue(r *bufio.Reader, lt logType) (int64, error) {
	sizeval, err := r.ReadString('\n')
	var defaultVal int64
	if lt == App {
		defaultVal = defaultMaxAppLogSize
	} else {
		defaultVal = defaultMaxWebLogSize
	}
	if err != nil {
		return defaultVal, fmt.Errorf("Unable to read response: %s", err)
	}
	if sizeval == newline {
		return defaultVal, nil
	}
	response, err := strconv.Atoi(strings.Trim(sizeval, newline))
	if err != nil {
		return defaultVal,
			fmt.Errorf("[ERROR] Maximum log file size must be a positive number")
	}
	if response <= 0 {
		return defaultVal,
			fmt.Errorf("[ERROR] Maximum log file size must be a positive number")
	}
	return int64(response), nil
}

func getMaxLogCountValue(r *bufio.Reader, lt logType) (int, error) {
	sizeval, err := r.ReadString('\n')
	var defaultVal int
	if lt == App {
		defaultVal = defaultMaxAppLogCount
	} else {
		defaultVal = defaultMaxWebLogCount
	}
	if err != nil {
		return defaultVal, fmt.Errorf("Unable to read response: %s", err)
	}
	if sizeval == newline {
		return defaultVal, nil
	}
	response, err := strconv.Atoi(strings.Trim(sizeval, newline))
	if err != nil {
		return defaultVal,
			fmt.Errorf("[ERROR] Maximum log count must be a positive number")
	}
	if response <= 0 {
		return defaultVal,
			fmt.Errorf("[ERROR] Maximum log count must be a positive number")
	}
	return response, nil
}

func getMaxHostsValue(r *bufio.Reader) (int, error) {
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
	if err := os.Mkdir(ConfigDirectory, os.ModeDir); err != nil {
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
	if lt == App {
		if defaultEnableAppLogging {
			return "yes"
		}
		return "no"
	}
	if defaultEnableWebLogging {
		return "yes"
	}
	return "no"
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
			"Enter whether application-related info and error messages should be logged to disk; 'yes' or 'no' [default: %s]: ",
			getBoolLogString(App))
		logAppEnableVal, err = getLoggingValue(reader, App)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.EnableAppLogging = logAppEnableVal
			validLogAppEnableVal = true
		}
	}
	if logAppEnableVal {
		validLogAppMaxSizeVal := false
		for !validLogAppMaxSizeVal {
			fmt.Printf(
				"Enter the maximum file size for application log files in Kilobytes. By default this is 1024 or %d megabyte(s). [default: %d]: ",
				defaultMaxAppLogSize/1024, defaultMaxAppLogSize)
			logMaxSizeVal, err := getMaxLogSizeValue(reader, App)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumAppLogSize = logMaxSizeVal
				validLogAppMaxSizeVal = true
			}
		}
		validLogAppMaxCountVal := false
		for !validLogAppMaxCountVal {
			fmt.Printf(
				"Enter the maximum number of application log files to keep. [default: %d]: ",
				defaultMaxAppLogCount)
			logMaxCountVal, err := getMaxLogCountValue(reader, App)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumAppLogCount = logMaxCountVal
				validLogAppMaxCountVal = true
			}
		}
	}
	// Configure Web-related logging
	validLogWebEnableVal := false
	var logWebEnableVal bool
	for !validLogWebEnableVal {
		var err error
		fmt.Printf(
			"Enter whether API web-related info and error messages should be logged to disk; 'yes' or 'no' [default: %s]: ",
			getBoolLogString(Web))
		logWebEnableVal, err = getLoggingValue(reader, Web)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.EnableWebLogging = logWebEnableVal
			validLogWebEnableVal = true
		}
	}
	if logWebEnableVal {
		validLogWebMaxSizeVal := false
		for !validLogWebMaxSizeVal {
			fmt.Printf(
				"Enter the maximum file size for API web log files in Kilobytes. By default this is 1024, or %d megabyte(s). [default: %d]: ",
				defaultMaxWebLogSize/1024, defaultMaxWebLogSize)
			logMaxSizeVal, err := getMaxLogSizeValue(reader, Web)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumWebLogSize = logMaxSizeVal
				validLogWebMaxSizeVal = true
			}
		}
		validLogWebMaxCountVal := false
		for !validLogWebMaxCountVal {
			fmt.Printf(
				"Enter the maximum number of API web log files to keep. [default: %d]: ",
				defaultMaxWebLogCount)
			logMaxCountVal, err := getMaxLogCountValue(reader, Web)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumWebLogCount = logMaxCountVal
				validLogWebMaxCountVal = true
			}
		}
	}
	// Configure maximum # of servers to retrieve
	validMaxHostsVal := false
	for !validMaxHostsVal {
		fmt.Println(
			"Enter the maximum number of servers to retrieve from the Steam Master Server at a time.")
		fmt.Printf("This can be no more than 6930. [default: %d]: ", defaultMaxHostsToReceive)
		maxHostVal, err := getMaxHostsValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.MaximumHostsToReceive = maxHostVal
			validMaxHostsVal = true
		}
	}
	// Configure time between master server queries
	validTimeBetweenVal := false
	for !validTimeBetweenVal {
		fmt.Println("Enter the time, in seconds, between Master Server queries.")
		fmt.Printf(
			"This must be greater than 20 & should not be too low as receiving servers can take a while. [default: %d]: ",
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
		fmt.Printf("Enter the port number on which the API web server will listen. [default: %d]: ",
			defaultAPIWebPort)
		apiPortVal, err := getAPIWebPortValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.APIWebPort = apiPortVal
			validAPIPortVal = true
		}
	}
	if err := writeConfig(cfg); err != nil {
		return err
	}
	return nil
}
