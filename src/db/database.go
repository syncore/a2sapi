package db

// database.go - Database initilization.

import (
	"a2sapi/src/constants"
	"a2sapi/src/logger"
	"a2sapi/src/util"
	"fmt"
)

// CountryDB is a package-level variable that contains a country
// geoelocation database connection. It is initialized once for re-usability
// when building server lists.
var CountryDB *CDB

// ServerDB is a package-level variable that contains a server information
// database connection. It is initialized once for re-usability when building
// server lists.
var ServerDB *SDB

// InitDBs initializes the geolocation and server information databases for
// re-use across server list builds. Panics on failure to initialize.
func InitDBs() {
	if CountryDB != nil && ServerDB != nil {
		return
	}

	cdb, err := OpenCountryDB()
	if err != nil {
		panic(fmt.Sprintf("Unable to initialize country database connection: %s",
			err))
	}
	sdb, err := OpenServerDB()
	if err != nil {
		panic(fmt.Sprintf(
			"Unable to initialize server information database connection: %s", err))
	}
	// Set package-level variables
	CountryDB = cdb
	ServerDB = sdb
}

func verifyServerDbPath() error {
	if err := util.CreateDirectory(constants.DbDirectory); err != nil {
		logger.LogAppError(err)
		panic(fmt.Sprintf("Unable to create database directory %s: %s",
			constants.DbDirectory, err))
	}
	if err := createServerDBtable(constants.GetServerDBPath()); err != nil {
		logger.LogAppErrorf("Unable to verify database path: %s", err)
		panic("Unable to verify database path")
	}

	return nil
}
