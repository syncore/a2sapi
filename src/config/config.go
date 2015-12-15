package config

// config.go - configuration operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"steamtest/src/constants"
	"steamtest/src/steam/filters"
	"steamtest/src/util"
)

var newline = getNewLineForOS()

// Config represents logging, steam-related, and API-related options.
type Config struct {
	LogConfig   CfgLog   `json:"logConfig"`
	SteamConfig CfgSteam `json:"steamConfig"`
	WebConfig   CfgWeb   `json:"webConfig"`
}

func getNewLineForOS() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// ReadConfig reads the configuration file from disk and returns a pointer to
// a struct that contains the various configuration values if successful, otherwise
// panics.
func ReadConfig() *Config {
	f, err := os.Open(constants.ConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf(
			"Error reading config file.\nYou might need to recreate with --config switch.\nError: %s",
			err))
	}
	defer f.Close()
	r := bufio.NewReader(f)
	d := json.NewDecoder(r)
	cfg := &Config{}
	if err := d.Decode(cfg); err != nil {
		panic(fmt.Sprintf(
			"Error decoding config file.\nYou might need to recreate with --config switch.\nError: %s",
			err))
	}
	return cfg
}

func getBoolString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// CreateConfig initiates the configuration creation process by collecting user
// input for various configuration values and then writes the configuration file
// to disk if successful, otherwise panics.
func CreateConfig() {
	reader := bufio.NewReader(os.Stdin)
	cfg := &Config{
		LogConfig:   CfgLog{},
		SteamConfig: CfgSteam{},
		WebConfig:   CfgWeb{},
	}
	fmt.Printf("%s - configuration file creation\n", constants.AppInfo)
	fmt.Print(
		"Type a value and press 'ENTER'. Leave a value empty and press 'ENTER' to use the default value.\n\n")

	// Logging configuration
	// Determine if application, Steam, and/or web API logging should be enabled
	cfg.LogConfig.EnableAppLogging = configureLoggingEnable(reader, constants.LTypeApp)
	cfg.LogConfig.EnableSteamLogging = configureLoggingEnable(reader, constants.LTypeSteam)
	cfg.LogConfig.EnableWebLogging = configureLoggingEnable(reader, constants.LTypeWeb)
	// Debug mode for testing (no user option to enable)
	cfg.LogConfig.EnableDebugMessages = defaultEnableDebugMessages
	// Configure max log size and max log count if logging is enabled
	if cfg.LogConfig.EnableAppLogging || cfg.LogConfig.EnableSteamLogging ||
		cfg.LogConfig.EnableWebLogging {
		cfg.LogConfig.MaximumLogSize = configureMaxLogSize(reader)
		cfg.LogConfig.MaximumLogCount = configureMaxLogCount(reader)
	} else {
		cfg.LogConfig.MaximumLogSize = defaultMaxLogSize
		cfg.LogConfig.MaximumLogCount = defaultMaxLogCount
	}

	// Steam configuration
	// Query the master server automatically at timed intervals
	cfg.SteamConfig.AutoQueryMaster = configureTimedMasterQuery(reader)
	if cfg.SteamConfig.AutoQueryMaster {
		// The game to automatically query the master server for at timed intervals
		cfg.SteamConfig.AutoQueryGame = configureTimedQueryGame(reader)
		// Time between Steam Master server queries
		cfg.SteamConfig.TimeBetweenMasterQueries = configureTimeBetweenQueries(reader,
			cfg.SteamConfig.AutoQueryGame)
		// Maximum # of servers to retrieve from Steam Master server
		cfg.SteamConfig.MaximumHostsToReceive = configureMaxServersToRetrieve(reader)
	} else {
		cfg.SteamConfig.AutoQueryGame = filters.GameQuakeLive.Name
		cfg.SteamConfig.TimeBetweenMasterQueries = defaultTimeBetweenMasterQueries
		cfg.SteamConfig.MaximumHostsToReceive = defaultMaxHostsToReceive
	}
	// # hours before bugged "stuck" players are filtered out from the results
	cfg.SteamConfig.SteamBugPlayerTime = configureSteamBugPlayerTime(reader)

	// Web API configuration
	// Direct queries: whether users can query any host (not just those with IDs)
	cfg.WebConfig.AllowDirectUserQueries = configureDirectQueries(reader)
	// Maximum number of servers to allow users to query via API
	cfg.WebConfig.MaximumHostsPerAPIQuery = configureMaxHostsPerAPIQuery(reader)
	// Time in seconds before HTTP requests time out
	cfg.WebConfig.APIWebTimeout = configureWebTimeout(reader)
	// Port that API's web server will listen on
	cfg.WebConfig.APIWebPort = configureWebServerPort(reader)

	if err := util.WriteJSONConfig(cfg, constants.ConfigDirectory,
		constants.ConfigFilePath); err != nil {
		panic(err)
	}
}
