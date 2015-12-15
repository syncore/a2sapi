package constants

// misc_constants.go - Miscellaneous-related constants (and a few variables)

import (
	"fmt"
	"path"
)

const (
	// GameFile specifies the name of the Steam games file.
	GameFile = "games.conf"
)

var (
	// GameFileFullPath represents the OS-independent full path to the game file.
	GameFileFullPath = path.Join(ConfigDirectory, GameFile)
	// Version is the version number of the application.
	Version = "0.1"
	// AppInfo contains the application information.
	AppInfo = fmt.Sprintf("steamtest v%s by syncore <syncore@syncore.org>", Version)
)
