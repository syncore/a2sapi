package test

import (
	"fmt"
	"os"
	"steamtest/src/constants"
)

// Cleanup removes testing files.
func Cleanup(filepaths ...string) {
	if !constants.IsTest {
		return
	}
	fmt.Println("Running test cleanup...")
	for _, fps := range filepaths {
		err := os.RemoveAll(fps)
		if err != nil {
			fmt.Printf("Error running test cleanup; unable to remove %s: %s", fps, err)
		}
	}
}
