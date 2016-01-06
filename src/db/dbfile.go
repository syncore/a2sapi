package db

// dbfile.go - database file operations

import (
	"a2sapi/src/constants"
	"a2sapi/src/logger"
	"a2sapi/src/util"
	"fmt"
)

func verifyServerDbPath() error {
	if err := util.CreateDirectory(constants.DbDirectory); err != nil {
		logger.LogAppError(err)
		panic(fmt.Sprintf("Unable to create database directory %s: %s",
			constants.DbDirectory, err))
	}
	if err := createServerDB(constants.GetServerDBPath()); err != nil {
		logger.LogAppErrorf("Unable to verify database path: %s", err)
		panic("Unable to verify database path")
	}

	return nil
}
