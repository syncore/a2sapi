package constants

// db_constants.go - Database-related constants (and a few variables)

import "path"

const (
	// DbDirectory specifies the directory in which to store the database files.
	DbDirectory = "db"
	// ServerDbFilename specifies the name of the server database file.
	ServerDbFilename = "servers.sqlite"
	// CountryMMDbFilename specifies the name of geolocation database file.
	CountryMMDbFilename = "GeoLite2-City.mmdb"
)

var (
	// CountryDbFilePath represents the OS-independent full path to the geolocation DB file.
	CountryDbFilePath = path.Join(DbDirectory, CountryMMDbFilename)
	// ServerDbFilePath represents the OS-independent full path to the server DB file.
	ServerDbFilePath = path.Join(DbDirectory, ServerDbFilename)
)
