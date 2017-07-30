// Package db contains functions to connect to an postgres database and useful functions for getting/manipulating the
// data.
package db

import (
	"database/sql"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

var sqlDB *sql.DB

// init initialises connection to sqlite3 database.
func init() {
	sqlDB, err := sql.Open("sqlite3", "viewer.db")
	if err != nil {
		log.Fatalln("Failed to initialise a connection to sqlite3 database:", err.Error())
	}

	var count int
	row := sqlDB.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND name='users';")
	row.Scan(&count)
	if count != 1 {
		if err := createUsersTable(); err != nil {
			log.Fatalln("Failed to create user's table in sqlite3 database:", err.Error())
		}
	}
}

// createUsersTable initialises the user table if not already present and creates a default admin.
func createUsersTable() error {
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
	_, err := sqlDB.Exec(sqlCreateUsersTable)
	if err != nil {
		return err
	}

	if err := createDefaultAdmin(); err != nil {
		return err
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
