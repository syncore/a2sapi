package config

// debugconfig.go - Options for debugging/development; not user-selectable

const (
	defaultEnableDebugMessages        = false
	defaultEnableServerDump           = false
	defaultServerDumpFileAsMasterList = false
	defaultServerDumpFile             = "api-test-servers.json"
)

// CfgDebug represents options for debugging and development.
type CfgDebug struct {
	// stdout "debug" msgs
	EnableDebugMessages bool `json:"debugMessages"`
	// dump server JSON to disk
	EnableServerDump bool `json:"dumpServers"`
	// use a pre-defined server JSON file as the master list for API (for testing)
	ServerDumpFileAsMasterList bool `json:"useServerDumpAsMaster"`
	// name of the pre-defined server JSON file to use as master list
	ServerDumpFilename string `json:"serverDumpFilename"`
}
