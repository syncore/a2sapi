package db

// dbfile.go - database file operations

import (
	"fmt"
	"os"
	"path"
	"steamtest/src/util"
)

const (
	serverDbFile    = "servers.sqlite"
	countryMMDbFile = "GeoLite2-City.mmdb"
	dbDirectory     = "db"
)

var (
	countryDbFilePath = path.Join(dbDirectory, countryMMDbFile)
	serverDbFilepath  = path.Join(dbDirectory, serverDbFile)
)

func verifyServerDbPath() error {
	if err := createDbDir(); err != nil {
		util.LogAppError(err)
		panic(fmt.Sprintf("Unable to create database directory %s: %s", dbDirectory,
			err))
	}
	if err := createServerDB(serverDbFilepath); err != nil {
		util.LogAppErrorf("Unable to verify database path: %s", err)
		panic("Unable to verify database path")
	}

	return nil
}

func createDbDir() error {
	if util.DirExists(dbDirectory) {
		return nil
	}
	if err := os.Mkdir(dbDirectory, os.ModePerm); err != nil {
		return err
	}
	return nil
}
