package constants

// config_constants.go - Configuration-related constants (and a few variables)

import "path"

const (
	// ConfigDirectory specifies the directory in which to store the config file.
	ConfigDirectory = "conf"
	// ConfigFilename specifies the name of the configuration file.
	ConfigFilename = "config.conf"
)

var (
	// ConfigFilePath represents the OS-independent full path to the config file.
	ConfigFilePath = path.Join(ConfigDirectory, ConfigFilename)
)
