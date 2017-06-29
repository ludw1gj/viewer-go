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

func Load() (err error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	return nil
}
