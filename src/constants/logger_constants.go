package constants

// logger_constants.go - Logger-related constants (and a few variables)

import "path"

const (
	// LogDirectory specifies the directory in which to store the log files.
	LogDirectory = "logs"
	// AppLogFilename specifies the name of the application log file.
	AppLogFilename = "app.log"
	// SteamLogFilename specifies the name of the Steam log file.
	SteamLogFilename = "steam.log"
	// WebLogFilename specifies the name of the web log file.
	WebLogFilename = "web.log"
)

// LogType represents the type of log.
type LogType int

const (
	// LTypeApp represents the Application-related log type.
	LTypeApp LogType = iota
	// LTypeDebug represents the Debug-related log type.
	LTypeDebug
	// LTypeSteam represents the Steam-related log type.
	LTypeSteam
	// LTypeWeb represents the Web-related log type.
	LTypeWeb
)

var (
	// AppLogFilePath represents the OS-independent full path to app log file.
	AppLogFilePath = path.Join(LogDirectory, AppLogFilename)
	// SteamLogFilePath represents the OS-independent full path to Steam log file.
	SteamLogFilePath = path.Join(LogDirectory, SteamLogFilename)
	// WebLogFilePath represents the OS-independent full path to web log file.
	WebLogFilePath = path.Join(LogDirectory, WebLogFilename)
)
