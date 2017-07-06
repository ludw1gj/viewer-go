// Package db contains functions to connect to an postgres database and useful functions for getting/manipulating the
// data.
package db

import (
	"database/sql"
	"fmt"

	"github.com/FriedPigeon/viewer-go/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Load initialises a connection to a postgres database using the values from config.Config type.
func Load(c config.Config) (err error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		c.Database.User, c.Database.Password, c.Database.Name)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	return db.Ping()
}
