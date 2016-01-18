package config

// config.go - configuration operations

import (
	"a2sapi/src/constants"
	"a2sapi/src/steam/filters"
	"a2sapi/src/util"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
)

var promptColor = color.New(color.FgHiGreen).SprintfFunc()
var errorColor = color.New(color.FgHiRed).PrintlnFunc()
var newline = getNewLineForOS()

// Config represents the application-wide configuration.
var Config *Cfg

// Cfg represents logging, steam-related, and API-related options.
type Cfg struct {
	LogConfig   CfgLog   `json:"logConfig"`
	SteamConfig CfgSteam `json:"steamConfig"`
	WebConfig   CfgWeb   `json:"webConfig"`
	DebugConfig CfgDebug `json:"debugConfig"`
}

func getNewLineForOS() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// InitConfig reads the configuration file from disk and if successful, sets the
// application wide-configuration. Otherwise, it will panic.
func InitConfig() {
	if Config != nil {
		return
	}

	f, err := os.Open(constants.GetCfgPath())
	if err != nil {
		panic(fmt.Sprintf(`
"Error reading config file. You might need to recreate it by using
the --config switch. Error: %s`, err))
	}
	defer f.Close()
	r := bufio.NewReader(f)
	d := json.NewDecoder(r)
	cfg := &Cfg{}
	if err := d.Decode(cfg); err != nil {
		panic(fmt.Sprintf(`
"Error decoding config file. You might need to recreate it by using
the --config switch. Error: %s`, err))
	}
	// Set the configuration which will live throughout the application's lifetime
	Config = cfg
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
	cfg := &Cfg{
		LogConfig:   CfgLog{},
		SteamConfig: CfgSteam{},
		WebConfig:   CfgWeb{},
		DebugConfig: CfgDebug{},
	}
	color.Set(color.FgHiYellow)
	fmt.Printf(`
%s - configuration file creation
Type a value and press 'ENTER'. Leave a value empty and press 'ENTER' to use the
default value.

`, constants.AppInfo)
	color.Unset()

	// Logging configuration
	// Determine if application, Steam, and/or web API logging should be enabled
	cfg.LogConfig.EnableAppLogging = configureLoggingEnable(reader, constants.LTypeApp)
	cfg.LogConfig.EnableSteamLogging = configureLoggingEnable(reader, constants.LTypeSteam)
	cfg.LogConfig.EnableWebLogging = configureLoggingEnable(reader, constants.LTypeWeb)
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

	// Web API configuration
	// Direct queries: whether users can query any host (not just those with IDs)
	cfg.WebConfig.AllowDirectUserQueries = configureDirectQueries(reader,
		cfg.SteamConfig.AutoQueryMaster)
	// Maximum number of servers to allow users to query via API
	cfg.WebConfig.MaximumHostsPerAPIQuery = configureMaxHostsPerAPIQuery(reader)
	// Time in seconds before HTTP requests time out
	cfg.WebConfig.APIWebTimeout = configureWebTimeout(reader)
	// Port that API's web server will listen on
	cfg.WebConfig.APIWebPort = configureWebServerPort(reader)
	// Enable or disable gzip compression of responses
	cfg.WebConfig.CompressResponses = configureResponseCompression(reader)

	// Debug configuration (not user-selectable. for debug/development purposes)
	// Print a few "debug" messages to stdout
	cfg.DebugConfig.EnableDebugMessages = defaultEnableDebugMessages
	// Dump the retrieved server information to a JSON file on disk
	cfg.DebugConfig.EnableServerDump = defaultEnableServerDump
	// Use a pre-defined JSON file as disk as the master server list for the API
	cfg.DebugConfig.ServerDumpFileAsMasterList = defaultServerDumpFileAsMasterList
	// Name of the pre-defined JSON file to use as the master server list for API
	cfg.DebugConfig.ServerDumpFilename = defaultServerDumpFile

	if err := util.WriteJSONConfig(cfg, constants.ConfigDirectory,
		constants.ConfigFilePath); err != nil {
		panic(err)
	}
}

// CreateDebugConfig creates the configuration file that is used when running the
// applciation in debug mode.
func CreateDebugConfig() {
	cfg := &Cfg{}
	cfg.LogConfig.EnableAppLogging = true
	cfg.LogConfig.EnableSteamLogging = false // even in debug mode; disable
	cfg.LogConfig.EnableWebLogging = true
	cfg.LogConfig.MaximumLogCount = defaultMaxLogCount
	cfg.LogConfig.MaximumLogSize = defaultMaxLogSize
	cfg.SteamConfig.AutoQueryMaster = false
	cfg.SteamConfig.AutoQueryGame = "QuakeLive"
	cfg.SteamConfig.TimeBetweenMasterQueries = defaultTimeBetweenMasterQueries
	cfg.SteamConfig.MaximumHostsToReceive = defaultMaxHostsToReceive
	cfg.WebConfig.AllowDirectUserQueries = true
	cfg.WebConfig.APIWebPort = defaultAPIWebPort
	cfg.WebConfig.APIWebTimeout = defaultAPIWebTimeout
	cfg.WebConfig.CompressResponses = defaultCompressResponses
	cfg.WebConfig.MaximumHostsPerAPIQuery = defaultMaxHostsPerAPIQuery
	cfg.DebugConfig.EnableDebugMessages = true
	cfg.DebugConfig.EnableServerDump = true
	cfg.DebugConfig.ServerDumpFileAsMasterList = true
	cfg.DebugConfig.ServerDumpFilename = defaultServerDumpFile
	if err := util.WriteJSONConfig(cfg, constants.ConfigDirectory,
		constants.DebugConfigFilePath); err != nil {
		panic(err)
	}
	// Set the configuration which will live throughout the application's lifetime
	Config = cfg
}

// CreateTestConfig creates the configuration that is used when running automated
// testing.
func CreateTestConfig() {
	// boolean values intentionally default to false and are omitted unless
	// otherwise specified, which is different from the normal configuration
	cfg := &Cfg{}
	cfg.LogConfig.MaximumLogCount = defaultMaxLogCount
	cfg.LogConfig.MaximumLogSize = defaultMaxLogSize
	cfg.SteamConfig.AutoQueryGame = "QuakeLive"
	cfg.SteamConfig.TimeBetweenMasterQueries = defaultTimeBetweenMasterQueries
	cfg.SteamConfig.MaximumHostsToReceive = defaultMaxHostsToReceive
	cfg.WebConfig.AllowDirectUserQueries = true
	cfg.WebConfig.APIWebPort = 40081
	cfg.WebConfig.APIWebTimeout = defaultAPIWebTimeout
	cfg.WebConfig.CompressResponses = defaultCompressResponses
	cfg.WebConfig.MaximumHostsPerAPIQuery = defaultMaxHostsPerAPIQuery
	cfg.DebugConfig.ServerDumpFileAsMasterList = true
	cfg.DebugConfig.ServerDumpFilename = "test-api-servers.json"
	if err := util.WriteJSONConfig(cfg, constants.TestTempDirectory,
		constants.TestConfigFilePath); err != nil {
		panic(err)
	}
	// Set the configuration which will live throughout the application's lifetime
	Config = cfg
}
