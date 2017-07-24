// Package db contains functions to connect to an postgres database and useful functions for getting/manipulating the
// data.
package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Load initialises connection to sqlite3 database.
func Load() (err error) {
	db, err = sql.Open("sqlite3", "viewer.db")
	if err != nil {
		return err
	}

	if err := initUsersTable(); err != nil {
		return err
	}
	return nil
}

// initUsersTable initialises the user table if not already present and creates a default admin.
func initUsersTable() error {
	var count int
	row := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type=$1 AND name=$2;", "table", "users")
	row.Scan(&count)
	if count != 1 {
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
		_, err := db.Exec(sqlCreateUsersTable)
		if err != nil {
			return err
		}

		if err := createDefaultAdmin(); err != nil {
			return err
		}
	}
	return nil
}

// createDefaultAdmin creates a default admin user.
func createDefaultAdmin() error {
	admin := User{
		1,
		"admin",
		"John",
		"Smith",
		"password",
		"./admin",
		true,
	}
	err := CreateUser(admin)
	if err != nil {
		return err
	}
	return nil
}
