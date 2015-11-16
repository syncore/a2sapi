// config.go - configuration options and operations
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"steamtest/util"
	"strconv"
	"strings"
)

var (
	newline = getNewLineForOS()
)

const (
	defaultEnableLogging            = true
	defaultMaxLogSize               = 1024
	defaultMaxHostsToReceive        = 3800
	defaultTimeBetweenMasterQueries = 60
	defaultAPIWebPort               = 40080
	ConfigDirectory                 = "conf"
	ConfigFileName                  = "config.conf"
	Version                         = "0.1"
)

type Config struct {
	EnableLogging            bool  `json:"enableLogging"`
	MaximumLogSize           int64 `json:"maxLogFilesize"`
	MaximumHostsToReceive    int   `json:"maxHostsToReceive"`
	TimeBetweenMasterQueries int   `json:"timeBetweenMasterQueries"`
	APIWebPort               int   `json:"apiWebPort"`
}

func getNewLineForOS() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

func getLoggingValue(r *bufio.Reader) (bool, error) {
	enable, err := r.ReadString('\n')
	if err != nil {
		return defaultEnableLogging, fmt.Errorf("Unable to read respone: %s\n", err)
	}
	if enable == newline {
		return defaultEnableLogging, nil
	}
	response := strings.Trim(enable, newline)
	if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
		return true, nil
	} else if strings.EqualFold(response, "n") || strings.EqualFold(response,
		"no") {
		return false, nil
	} else {
		return defaultEnableLogging,
			fmt.Errorf("Invalid response. Valid responses: y, yes, n, no\n")
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
			fmt.Errorf("[ERROR] Maximum log file size must be a positive number.")
	}
	if response <= 0 {
		return defaultMaxLogSize,
			fmt.Errorf("[ERROR] Maximum log file size must be a positive number.")
	}
	return int64(response), nil
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
			fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930.")
	}
	if response < 500 || response > 6930 {
		return defaultMaxHostsToReceive,
			fmt.Errorf("[ERROR] Maximum hosts to receive from master server must be between 500 and 6930.")
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
			fmt.Errorf("[ERROR] Time between Steam master server queries must be a number greater than 20.")
	}
	if response < 20 {
		return defaultTimeBetweenMasterQueries,
			fmt.Errorf("[ERROR] Time between Steam master server queries must be a number greater than 20.")

	}
	return response, nil
}

func getApiWebPortValue(r *bufio.Reader) (int, error) {
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
			fmt.Errorf("[ERROR] API webserver port must be between 1 and 65535.")
	}
	if response < 1 || response > 65535 {
		return defaultAPIWebPort,
			fmt.Errorf("[ERROR] API webserver port must be between 1 and 65535.")
	}
	return response, nil
}

func createConfigDir() error {
	if util.DirExists(ConfigDirectory) {
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
	cfgpath := path.Join(ConfigDirectory, ConfigFileName)
	f, err := os.Create(cfgpath)
	if err != nil {
		fmt.Printf("Couldn't create '%s' file: %s\n", cfgpath, err)
		return err
	}
	defer f.Close()
	c, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling JSON for '%s' file: %s\n", cfgpath, err)
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
	fmt.Printf("Successfuly wrote config file to %s\n", cfgpath)
	return nil
}

func ReadConfig() (*Config, error) {
	cfgpath := path.Join(ConfigDirectory, ConfigFileName)
	f, err := os.Open(cfgpath)
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

func CreateConfig() error {
	reader := bufio.NewReader(os.Stdin)
	cfg := &Config{}
	fmt.Printf("steamtest v%s - configuration file creation\n", Version)
	fmt.Print(
		"Type a value and press 'ENTER'. Leave a value empty and press 'ENTER' to use the default value.\n\n")

	var defLogstr string
	if defaultEnableLogging {
		defLogstr = "yes"
	} else {
		defLogstr = "no"
	}
	validLogEnableVal := false
	var logEnableVal bool
	for !validLogEnableVal {
		var err error
		fmt.Printf(
			"Enter whether info and error messages should be logged to disk; 'yes' or 'no' [default: %s]: ",
			defLogstr)
		logEnableVal, err = getLoggingValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.EnableLogging = logEnableVal
			validLogEnableVal = true
		}
	}
	if logEnableVal {
		validLogMaxSizeVal := false
		for !validLogMaxSizeVal {
			fmt.Printf(
				"Enter the maximum file size for log files in Kilobytes. By default this is 1024, which is %d megabytes. [default: %d]: ",
				defaultMaxLogSize/1024, defaultMaxLogSize)
			logMaxSizeVal, err := getMaxLogSizeValue(reader)
			if err != nil {
				fmt.Println(err)
			} else {
				cfg.MaximumLogSize = logMaxSizeVal
				validLogMaxSizeVal = true
			}
		}
	}
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
	validApiPortVal := false
	for !validApiPortVal {
		fmt.Printf("Enter the port number on which the API web server will listen. [default: %d]: ",
			defaultAPIWebPort)
		apiPortVal, err := getApiWebPortValue(reader)
		if err != nil {
			fmt.Println(err)
		} else {
			cfg.APIWebPort = apiPortVal
			validApiPortVal = true
		}
	}
	if err := writeConfig(cfg); err != nil {
		return err
	}
	return nil
}
