// TODO: package doc

package database

import (
	"database/sql"

	"errors"

	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Load initialises connection to sqlite3 database.
func Load(dbFile string) (err error) {
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return errors.New("Failed to initialise a connection to sqlite3 database: " + err.Error())
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
		_, err := db.Exec(sqlCreateUsersTable)
		if err != nil {
			return errors.New("Failed to create user's table in sqlite3 database: " + err.Error())
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

		err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot,
			&user.Admin)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUser returns a single user from the database that matches the provided id.
func GetUser(id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot,
		&user.Admin)
	switch err {
	case sql.ErrNoRows:
		return user, errors.New("There is no user by that ID.")
	default:
		return user, err
	}
}

// CreateUser inserts a new user into the database.
func CreateUser(u User) error {
	// check if username is taken
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", u.Username)
	row.Scan(&count)
	if count == 1 {
		return errors.New("Username is taken.")
	}

	// create user root directory on disk
	err := os.MkdirAll(u.DirectoryRoot, os.ModePerm)
	if err != nil {
		return err
	}

	// generate hash of user password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// store user in database
	_, err = db.Exec("INSERT INTO users (username, first_name, last_name, password, directory_root, admin) "+
		"VALUES ($1, $2, $3, $4, $5, $6)",
		u.Username, u.FirstName, u.LastName, string(hashPassword), u.DirectoryRoot, u.Admin)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes the user from the database that corresponds to the given ID.
func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// ValidateUser validates a user with username and password. It will check if the username exists in the database
// and checks if the password is valid, then returning the user's id.
func ValidateUser(username string, password string) (userID int, err error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot,
		&user.Admin)
	if err != nil {
		return userID, errors.New("There is no user by that username.")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return userID, errors.New("Password is invalid.")
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
		return errors.New("Username does not exist.")
	}

	_, err := db.Exec("UPDATE users SET username = $1 WHERE username = $2;", newUsername, username)
	if err != nil {
		return err
	}
	return nil
}

// TODO: func to change admin status by user ID
