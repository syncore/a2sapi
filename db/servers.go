// servers.go - server identification database
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const serverDbFile = "servers.sqlite"

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createDb(dbfile string) error {
	if fileExists(dbfile) {
		return nil
	}

	_, err := os.Create(dbfile)
	if err != nil {
		return fmt.Errorf("Unable to create server DB: %s\n", err)
	}
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return fmt.Errorf("Unable to open server DB file for table creation: %s\n", err)
	}
	defer db.Close()
	q := `CREATE TABLE servers (
	server_id INTEGER NOT NULL,
	host TEXT NOT NULL,
	PRIMARY KEY(server_id)
	)`
	_, err = db.Exec(q)
	if err != nil {
		return fmt.Errorf("Unable to create servers table in servers DB: %s", err)
	}
	return nil
}

func serverExists(db *sql.DB, host string) (bool, error) {
	rows, err := db.Query("SELECT host FROM servers WHERE host =? LIMIT 1",
		host)
	if err != nil {
		// TODO: log this to disk
		fmt.Printf("Error querying database for host %s: %s\n", host, err)
		return false, err
	}

	defer rows.Close()
	var h string
	for rows.Next() {
		if err := rows.Scan(&h); err != nil {
			// TODO: log this to disk
			fmt.Printf("Error querying database for host %s: %s\n", host, err)
			return false, err
		}
	}
	if h != "" {
		return true, nil
	} else {
		return false, nil
	}
}

func OpenServerDB() (*sql.DB, error) {
	err := createDb(serverDbFile)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", serverDbFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddServersToDB(db *sql.DB, hosts []string) {
	for _, h := range hosts {
		exists, err := serverExists(db, h)
		if err != nil {
			// TODO: log the error to disk
			fmt.Printf("AddServersToDB exists error: %s\n", err)
			continue
		}
		if exists {
			continue
		}
		_, err = db.Exec("INSERT INTO servers (host) VALUES ($1)", h)
		if err != nil {
			// TODO: log the error to disk
			fmt.Printf("AddServersToDB exec error for host %s: %s\n", h, err)
			continue
		}
	}
}
