// Package database contains functions that load a sqlite3 database.
package database

import (
	"database/sql"

	"errors"

	"os"

	_ "github.com/mattn/go-sqlite3" // only the driver is needed
	"golang.org/x/crypto/bcrypt"
)

// CreateDBConn initialises connection to sqlite3 database.
func CreateDBConn(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, errors.New("failed to initialise a connection to sqlite3 database: " + err.Error())
	}

	var count int
	row := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='users';")
	row.Scan(&count)
	if count != 1 {
		// initialise the user table if not already present and creates a default admin.
		sqlCreateUsersTable := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			password TEXT NOT NULL,
			directory_root TEXT NOT NULL,
			admin BOOLEAN NOT NULL
		);`
		if _, err := db.Exec(sqlCreateUsersTable); err != nil {
			return nil, errors.New("failed to create user's table in sqlite3 database: " + err.Error())
		}

		// create a default admin user
		if err := os.MkdirAll("./data/users/admin", os.ModePerm); err != nil {
			return nil, err
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		if _, err := db.Exec("INSERT INTO users (username, first_name, last_name, password, directory_root, admin) "+
			"VALUES ($1, $2, $3, $4, $5, $6)", "admin", "John", "Smith", string(hashPassword), "./admin", true); err != nil {
			return nil, err
		}
	}
	return db, nil
}
