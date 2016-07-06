package constants

// misc_constants.go - Miscellaneous-related constants (and a few variables)

import (
	"fmt"
	"path"
)

const (
	// DumpDirectory represents the directory name used for server dump JSON files.
	DumpDirectory = "dump"
	// GameFile specifies the name of the Steam games file.
	GameFile = "games.conf"
)

var (
	// DumpFileFullPath represents the OS-independent full path to the server dump
	// JSON file.
	DumpFileFullPath = func(dumpfile string) string {
		return path.Join(DumpDirectory, dumpfile)
	}
	// GameFileFullPath represents the OS-independent full path to the game file.
	GameFileFullPath = path.Join(ConfigDirectory, GameFile)
	// Version is the version number of the application.
	Version = "0.1.5"
	// AppInfo contains the application information.
	AppInfo = fmt.Sprintf("a2sapi v%s by syncore <syncore@syncore.org>", Version)
)
