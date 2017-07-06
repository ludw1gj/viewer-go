// Package db contains functions to connect to an postgres database and useful functions for getting/manipulating the
// data.
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "test"
)

var db *sql.DB

// Load initialises a connection to a postgres database using the values of the constants in the package.
func Load() (err error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	return db.Ping()
}
