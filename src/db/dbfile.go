package db

// dbfile.go - database file operations

import (
	"fmt"
	"steamtest/src/constants"
	"steamtest/src/logger"
	"steamtest/src/util"
)

func verifyServerDbPath() error {
	if err := util.CreateDirectory(constants.DbDirectory); err != nil {
		logger.LogAppError(err)
		panic(fmt.Sprintf("Unable to create database directory %s: %s",
			constants.DbDirectory, err))
	}
	if err := createServerDB(constants.ServerDbFilePath); err != nil {
		logger.LogAppErrorf("Unable to verify database path: %s", err)
		panic("Unable to verify database path")
	}

	return nil
}
