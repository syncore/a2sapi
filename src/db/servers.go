package db

// servers.go - server identification database

import (
	"database/sql"
	"fmt"
	"os"
	"steamtest/src/util"
	// blank import for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const serverDbFile = "servers.sqlite"

func createDb(dbfile string) error {
	if util.FileExists(dbfile) {
		return nil
	}

	f, err := os.Create(dbfile)
	if err != nil {
		return util.LogAppErrorf("Unable to create server DB: %s", err)
	}
	defer f.Close()
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return util.LogAppErrorf(
			"Unable to open server DB file for table creation: %s", err)
	}
	defer db.Close()
	q := `CREATE TABLE servers (
	server_id INTEGER NOT NULL,
	host TEXT NOT NULL,
	PRIMARY KEY(server_id)
	)`
	_, err = db.Exec(q)
	if err != nil {
		return util.LogAppErrorf("Unable to create servers table in servers DB: %s",
			err)
	}
	return nil
}

func serverExists(db *sql.DB, host string) (bool, error) {
	rows, err := db.Query("SELECT host FROM servers WHERE host =? LIMIT 1",
		host)
	if err != nil {
		return false, util.LogAppErrorf("Error querying database for host %s: %s",
			host, err)
	}

	defer rows.Close()
	var h string
	for rows.Next() {
		if err := rows.Scan(&h); err != nil {
			return false, util.LogAppErrorf("Error querying database for host %s: %s",
				host, err)
		}
	}
	if h != "" {
		return true, nil
	}
	return false, nil
}

func OpenServerDB() (*sql.DB, error) {
	if err := createDb(serverDbFile); err != nil {
		return nil, util.LogAppError(err)
	}
	db, err := sql.Open("sqlite3", serverDbFile)
	if err != nil {
		return nil, util.LogAppError(err)
	}
	return db, nil
}

func AddServersToDB(db *sql.DB, hosts []string) {
	var toInsert []string
	for _, h := range hosts {
		exists, err := serverExists(db, h)
		if err != nil {
			continue
		}
		if exists {
			continue
		}
		toInsert = append(toInsert, h)
	}
	tx, err := db.Begin()
	if err != nil {
		util.LogAppErrorf("AddServersToDB error creating tx: %s", err)
		return
	}
	var txexecerr error
	for _, i := range toInsert {
		_, txexecerr = tx.Exec("INSERT INTO servers (host) VALUES ($1)", i)
		if txexecerr != nil {
			util.LogAppErrorf("AddServersToDB exec error for host %s: %s", i, err)
			break
		}
	}
	if txexecerr != nil {
		if err = tx.Rollback(); err != nil {
			util.LogAppErrorf("AddServersToDB error rolling back tx: %s", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		util.LogAppErrorf("AddServersToDB error committing tx: %s", err)
		return
	}
}

func GetIDsForServerList(result chan map[string]int64, db *sql.DB,
	hosts []string) {
	m := make(map[string]int64, len(hosts))
	for _, host := range hosts {
		rows, err := db.Query("SELECT server_id FROM servers WHERE host =? LIMIT 1",
			host)
		if err != nil {
			util.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
				host, err)
			return
		}

		defer rows.Close()
		var id int64
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				util.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
					host, err)
				return
			}
		}
		m[host] = id
	}
	result <- m
}

func GetIDsForAPIQuery(result chan map[int64]string, db *sql.DB, hosts []string) {
	m := make(map[int64]string, len(hosts))
	for _, h := range hosts {
		fmt.Printf("GetIDsForAPIQuery DB function, host: %s\n", h)
		rows, err := db.Query("SELECT server_id, host FROM servers WHERE host LIKE ?",
			fmt.Sprintf("%%%s%%", h))
		if err != nil {
			util.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
				h, err)
			return
		}

		defer rows.Close()
		var id int64
		var host string
		for rows.Next() {
			if err := rows.Scan(&id, &host); err != nil {
				util.LogAppErrorf("Error querying database to retrieve ID for host %s: %s",
					h, err)
				return
			}
			m[id] = host
		}
	}
	result <- m
}
