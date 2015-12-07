package util

// config.go - configuration operations

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
)

var (
	newline = getNewLineForOS()
	// ConfigFullPath represents the OS-independent full path to the config file.
	ConfigFullPath = path.Join(ConfigDirectory, ConfigFileName)
)

const (
	// ConfigDirectory specifies the directory in which to store the config file.
	ConfigDirectory = "conf"
	// ConfigFileName specifies the name of the configuration file.
	ConfigFileName = "config.conf"
	// Version is the version number of the application.
	Version = "0.1"
)

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

// ReadConfig reads the configuration file from disk and returns a pointer to
// a struct that contains the various configuration values if successful, otherwise
// returns an error.
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

func getBoolString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// CreateConfig initiates the configuration creation process by collecting user
// input for various configuration values and then writes the configuration file
// to disk if successful, otherwise returns an error.
func CreateConfig() error {
	reader := bufio.NewReader(os.Stdin)
	cfg := &Config{
		LogConfig:   CfgLog{},
		SteamConfig: CfgSteam{},
		WebConfig:   CfgWeb{},
	}
	fmt.Printf("steamtest v%s - configuration file creation\n", Version)
	fmt.Print(
		"Type a value and press 'ENTER'. Leave a value empty and press 'ENTER' to use the default value.\n\n")

	// Logging configuration
	// Determine if application, Steam, and/or web API logging should be enabled
	cfg.LogConfig.EnableAppLogging = configureLoggingEnable(reader, App)
	cfg.LogConfig.EnableSteamLogging = configureLoggingEnable(reader, Steam)
	cfg.LogConfig.EnableWebLogging = configureLoggingEnable(reader, Web)
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
	// Maximum # of servers to retrieve from Steam Master server
	cfg.SteamConfig.MaximumHostsToReceive = configureMaxServersToRetrieve(reader)
	// # hours before bugged "stuck" players are filtered out from the results
	cfg.SteamConfig.SteamBugPlayerTime = configureSteamBugPlayerTime(reader)
	// Time between Steam Master server queries
	cfg.SteamConfig.TimeBetweenMasterQueries = configureTimeBetweenQueries(reader)

	// Web API configuration
	// Direct queries: whether users can query any host (not just those with IDs)
	cfg.WebConfig.AllowDirectUserQueries = configureDirectQueries(reader)
	// Maximum number of servers to allow users to query via API
	cfg.WebConfig.MaximumHostsPerAPIQuery = configureMaxHostsPerAPIQuery(reader)
	// Time in seconds before HTTP requests time out
	cfg.WebConfig.APIWebTimeout = configureWebTimeout(reader)
	// Port that API's web server will listen on
	cfg.WebConfig.APIWebPort = configureWebServerPort(reader)

	if err := writeConfig(cfg); err != nil {
		return err
	}
	return nil
}
