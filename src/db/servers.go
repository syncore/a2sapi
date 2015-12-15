package db

// servers.go - server identification database

import (
	"database/sql"
	"fmt"
	"os"
	"steamtest/src/constants"
	"steamtest/src/logger"
	"steamtest/src/models"
	"steamtest/src/steam/filters"
	"steamtest/src/util"
	// blank import for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func createServerDB(dbfile string) error {
	create := `CREATE TABLE servers (
	server_id INTEGER NOT NULL,
	host TEXT NOT NULL,
	game TEXT NOT NULL,
	PRIMARY KEY(server_id)
	)`

	if util.FileExists(dbfile) {
		// already exists, so verify integrity
		db, err := sql.Open("sqlite3", dbfile)
		if err != nil {
			return logger.LogAppErrorf(
				"Unable to open server DB file for verification: %s", err)
		}
		defer db.Close()
		var name string
		err = db.QueryRow(
			"SELECT name from sqlite_master where type='table' and name='servers'").Scan(&name)
		switch {
		case err == sql.ErrNoRows:
			if _, err = db.Exec(create); err != nil {
				return logger.LogAppErrorf("Unable to create servers table in DB: %s", err)
			}
		case err != nil:
			return logger.LogAppErrorf("Server DB table verification error: %s", err)
		}
		return nil
	}

	f, err := os.Create(dbfile)
	if err != nil {
		return logger.LogAppErrorf("Unable to create server DB: %s", err)
	}
	defer f.Close()
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return logger.LogAppErrorf(
			"Unable to open server DB file for table creation: %s", err)
	}
	defer db.Close()
	_, err = db.Exec(create)
	if err != nil {
		return logger.LogAppErrorf("Unable to create servers table in servers DB: %s",
			err)
	}
	return nil
}

func serverExists(db *sql.DB, host string, game string) (bool, error) {
	rows, err := db.Query(
		"SELECT host, game FROM servers WHERE host =? AND GAME =? LIMIT 1",
		host, game)
	if err != nil {
		return false, logger.LogAppErrorf(
			"Error querying database for host %s and game %s: %s", host, game, err)
	}

	defer rows.Close()
	var h string
	var g string
	for rows.Next() {
		if err := rows.Scan(&h, &g); err != nil {
			return false, logger.LogAppErrorf(
				"Error querying database for host %s and game %s: %s",
				host, game, err)
		}
	}
	if h != "" && g != "" {
		return true, nil
	}
	return false, nil
}

func getHostAndGame(db *sql.DB, id string) (host, game string, err error) {
	rows, err := db.Query("SELECT host, game FROM servers WHERE server_id =? LIMIT 1",
		id)
	if err != nil {
		return host, game, logger.LogAppErrorf("Error querying database for id %s: %s",
			id, err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&host, &game); err != nil {
			return host, game, logger.LogAppErrorf("Error querying database for id %s: %s",
				id, err)
		}
	}
	return host, game, nil
}

// OpenServerDB Opens a database connection to the server database file or if
// that file does not exists, creates it and then opens a database connection to it.
func OpenServerDB() (*sql.DB, error) {
	if err := verifyServerDbPath(); err != nil {
		// will panic if not verified
		return nil, logger.LogAppError(err)
	}
	db, err := sql.Open("sqlite3", constants.ServerDbFilePath)
	if err != nil {
		return nil, logger.LogAppError(err)
	}
	return db, nil
}

// AddServersToDB inserts a specified host and port with its game name into the
// server database.
func AddServersToDB(db *sql.DB, hostsgames map[string]string) {
	toInsert := make(map[string]string, len(hostsgames))
	for host, game := range hostsgames {
		// If direct queries are enabled, don't add 'Unspecified' game to server DB
		if game == filters.GameUnspecified.String() {
			continue
		}
		exists, err := serverExists(db, host, game)
		if err != nil {
			continue
		}
		if exists {
			continue
		}
		toInsert[host] = game
	}
	tx, err := db.Begin()
	if err != nil {
		logger.LogAppErrorf("AddServersToDB error creating tx: %s", err)
		return
	}
	var txexecerr error
	for host, game := range toInsert {
		_, txexecerr = tx.Exec("INSERT INTO servers (host, game) VALUES ($1, $2)",
			host, game)
		if txexecerr != nil {
			logger.LogAppErrorf(
				"AddServersToDB exec error for host %s and game %s: %s", host, game, err)
			break
		}
	}
	if txexecerr != nil {
		if err = tx.Rollback(); err != nil {
			logger.LogAppErrorf("AddServersToDB error rolling back tx: %s", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		logger.LogAppErrorf("AddServersToDB error committing tx: %s", err)
		return
	}
}

// GetIDsForServerList Retrieves the server ID numbers for a given set of hosts,
// from the server database file, in response to a request to build the master
// server detail list or the list of server details in response to a request
// coming in over the API. It sends its results over a map channel consisting of
// a host to id mapping.
func GetIDsForServerList(result chan map[string]int64, db *sql.DB,
	hosts map[string]string) {
	m := make(map[string]int64, len(hosts))
	for host, game := range hosts {
		rows, err := db.Query(
			"SELECT server_id FROM servers WHERE host =? AND game =? LIMIT 1",
			host, game)
		if err != nil {
			logger.LogAppErrorf(
				"Error querying database to retrieve ID for host %s and game %s: %s",
				host, game, err)
			return
		}
		defer rows.Close()
		var id int64
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				logger.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
					host, err)
				return
			}
		}
		m[host] = id
	}
	result <- m
}

// GetIDsAPIQuery Retrieves the server ID numbers, hosts, and game name for a given
// set of hosts (represented by query string values) from the server database
// file in response to a query from the API. Sends the results over a DbServerID
// channel for consumption.
func GetIDsAPIQuery(result chan *models.DbServerID, db *sql.DB, hosts []string) {
	m := &models.DbServerID{}
	for _, h := range hosts {
		logger.WriteDebug("DB: GetIDsAPIQuery, host: %s", h)
		rows, err := db.Query(
			"SELECT server_id, host, game FROM servers WHERE host LIKE ?",
			fmt.Sprintf("%%%s%%", h))
		if err != nil {
			logger.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
				h, err)
			return
		}
		defer rows.Close()
		var id int64
		var host string
		var game string

		for rows.Next() {
			sid := &models.DbServer{}
			if err := rows.Scan(&id, &host, &game); err != nil {
				logger.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
					h, err)
				return
			}
			sid.ID = id
			sid.Host = host
			sid.Game = game
			m.Servers = append(m.Servers, sid)
		}
	}
	m.ServerCount = len(m.Servers)
	result <- m
}

// GetHostsAndGameFromIDAPIQuery Retrieves the hosts and game names from the
// server database file in response to a user-specified API query for a given
// set of server ID numbers. Sends the results over a channel consisting of a
// host to game name string mapping.
func GetHostsAndGameFromIDAPIQuery(result chan map[string]string, db *sql.DB,
	ids []string) {
	hosts := make(map[string]string, len(ids))
	for _, id := range ids {
		host, game, err := getHostAndGame(db, id)
		if err != nil {
			logger.LogAppErrorf("Error getting host from ID for API query: %s", err)
			return
		}
		if host == "" && game == "" {
			continue
		}
		hosts[host] = game
	}
	result <- hosts
}
