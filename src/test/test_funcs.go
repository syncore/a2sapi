package test

import (
	"a2sapi/src/config"
	"a2sapi/src/constants"
	"a2sapi/src/db"
	"fmt"
	"os"
)

// SetupEnvironment sets up the environment for tests. This should only be
// called once per package and only in the first _test.go file of the package
// that needs it.
func SetupEnvironment() {
	fmt.Println("Setting up test environment...")
	// Need base directory for config and other files
	err := os.Chdir("../../bin")
	if err != nil {
		panic("Unable to change directory for tests")
	}
	// Remove old test files
	deleteFiles(constants.TestTempDirectory)

	// Use testing configuration
	config.CreateTestConfig()
	constants.IsTest = true

	// Dump is not in test directory and needs config access
	deleteFiles(constants.DumpFileFullPath(
		config.Config.DebugConfig.ServerDumpFilename))

	// Initialize database connections
	db.InitDBs()
}

func deleteFiles(filepaths ...string) {
	fmt.Println("Running pre-test cleanup...")
	for _, fps := range filepaths {
		err := os.RemoveAll(fps)
		if err != nil {
			fmt.Printf("Error running test cleanup; unable to remove %s: %s", fps, err)
		}
	}
}
