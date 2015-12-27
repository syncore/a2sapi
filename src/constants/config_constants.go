package constants

// config_constants.go - Configuration-related constants (and a few variables)

import "path"

const (
	// ConfigDirectory specifies the directory in which to store the config file.
	ConfigDirectory = "conf"
	// ConfigFilename specifies the name of the configuration file.
	ConfigFilename = "config.conf"
	// DebugConfigFilename specifies the name of the configuration file to use when
	// debug mode is set
	DebugConfigFilename = "debug.conf"
)

var (
	// IsDebug will determine whether the debug configuration is used. This is
	// set on application startup.
	IsDebug = false
	// IsTest will determine whether the test configuration is used when running
	// tests. This variable is only set when running tests.
	IsTest = false
	// ConfigFilePath represents the OS-independent full path to the config file.
	ConfigFilePath = path.Join(ConfigDirectory, ConfigFilename)
	// DebugConfigFilePath represents the OS-independent full path to the debug
	// configuration file.
	DebugConfigFilePath = path.Join(ConfigDirectory, DebugConfigFilename)
)

// GetCfgPath returns the full OS-independent path to the configuration file.
func GetCfgPath() string {
	if IsTest {
		return path.Join(TestTempDirectory, TestConfigFilename)
	}
	if IsDebug {
		return path.Join(ConfigDirectory, DebugConfigFilename)
	}
	return path.Join(ConfigDirectory, ConfigFilename)
}
