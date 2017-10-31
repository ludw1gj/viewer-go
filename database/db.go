// Package database contains functions that load and manipulate an sqlite3 database.
package database

import (
	"database/sql"

	"errors"

	"os"

	"fmt"

	_ "github.com/mattn/go-sqlite3" // only the driver is needed
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Load initialises connection to sqlite3 database.
func Load(dbFile string) (err error) {
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return errors.New("failed to initialise a connection to sqlite3 database: " + err.Error())
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
			return errors.New("failed to create user's table in sqlite3 database: " + err.Error())
		}

		// create a default admin user.
		admin := User{
			1,
			"admin",
			"John",
			"Smith",
			"password",
			"./admin",
			true,
		}
		if err := CreateUser(admin); err != nil {
			return err
		}
	}
	return nil
}

// GetAllUsers returns all users in the database.
func GetAllUsers() (users []User, err error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}

		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot,
			&user.Admin); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUser returns a single user from the database that matches the provided id.
func GetUser(id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot,
		&user.Admin); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("there is no user by that ID")
		}
	}
	return user, nil
}

// CreateUser inserts a new user into the database.
func CreateUser(u User) error {
	// check if username is taken
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", u.Username)
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return errors.New("username is taken")
	}

	// create user root directory on disk
	if err := os.MkdirAll(u.DirectoryRoot, os.ModePerm); err != nil {
		return err
	}

	// generate hash of user password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// store user in database
	if _, err := db.Exec("INSERT INTO users (username, first_name, last_name, password, directory_root, admin) "+
		"VALUES ($1, $2, $3, $4, $5, $6)",
		u.Username, u.FirstName, u.LastName, string(hashPassword), u.DirectoryRoot, u.Admin); err != nil {
		return err
	}
	return nil
}

// ValidateUser validates a user with username and password. It will check if the username exists in the database
// and checks if the password is valid, then returning the user's id.
func ValidateUser(username string, password string) (userID int, err error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	if err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot,
		&user.Admin); err != nil {
		return userID, errors.New("there is no user by that username")
	}

	if err := comparePasswords(user.Password, password); err != nil {
		return userID, err
	}
	return user.ID, nil
}

// ChangeUserUsername updates the user's username.
func ChangeUserUsername(username string, newUsername string) error {
	// check if username exists
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username)
	row.Scan(&count)
	if count != 1 {
		return errors.New("username does not exist")
	}

	if _, err := db.Exec("UPDATE users SET username = $1 WHERE username = $2;", newUsername, username); err != nil {
		return err
	}
	return nil
}

// ChangeUserAdminStatus updates the user's admin status.
func ChangeUserAdminStatus(id int, isAdmin bool) error {
	if err := checkUserExists(id); err != nil {
		return err
	}
	if _, err := db.Exec("UPDATE users SET admin = $1 WHERE id = $2", id, isAdmin); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes the user from the database that corresponds to the given ID.
func DeleteUser(id int) error {
	if err := checkUserExists(id); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

// checkUserExists checks if user does exist.
func checkUserExists(id int) error {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id)
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("user by id %d does not exist", id)
	}
	return nil
}
