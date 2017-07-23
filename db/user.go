// This file contains database functions for interacting with User data.

package db

import (
	"log"

	"os"

	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User type contains user information.
type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Password      string `json:"password"`
	DirectoryRoot string `json:"directory_root"`
	Admin         bool   `json:"is_admin"`
}

// GetAllUsers returns all users in the database.
func GetAllUsers() (users []User, err error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}

		err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot, &user.Admin)
		if err != nil {
			log.Println(err)
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUser returns a single user from the database that matches the provided id.
func GetUser(id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot, &user.Admin)
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
	if count > 0 {
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

	// store user in db
	_, err = db.Exec("INSERT INTO users (username, first_name, last_name, password, directory_root, admin) VALUES ($1, $2, $3, $4, $5, $6)",
		u.Username, u.FirstName, u.LastName, string(hashPassword), u.DirectoryRoot, u.Admin)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database that matches the provided id.
func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// DeleteUserPasswordValidated deletes a user from the database if the provided password is valid.
func DeleteUserPasswordValidated(user User, password string) error {
	// check if password is valid
	err := comparePasswords(user.Password, password)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	return err
}

// ChangeUserPassword changes a user's password in the database, if the provided password is valid.
func ChangeUserPassword(user User, oldPassword string, newPassword string) error {
	// check if oldPassword is valid
	err := comparePasswords(user.Password, oldPassword)
	if err != nil {
		return err
	}

	// generate hash of new password
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	// store new password
	_, err = db.Exec("UPDATE users SET password = $1 WHERE id = $2;", newHashPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserValidation validates a user with username and password. It will check if the username exists in the database
// and checks if the password is valid, then returning the user's id.
func CheckUserValidation(username string, password string) (userID int, err error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot, &user.Admin)
	if err != nil {
		return userID, errors.New("There is no user by that username.")
	}
	err = comparePasswords(user.Password, password)
	return user.ID, err
}

// comparePasswords compares the hash string with a string to determine if it is equivalent.
func comparePasswords(hashPW string, pw string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPW), []byte(pw))
	if err != nil {
		return errors.New("Password is invalid.")
	}
	return nil
}
