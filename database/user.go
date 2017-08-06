// This file contains the user model and its methods for interacting with user data.

package database

import (
	"os"

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

// Delete deletes a user from the database if the provided password is valid.
func (u User) Delete(password string) error {
	// check if password is valid
	err := comparePasswords(u.Password, password)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", u.ID)
	return err
}

// UpdatePassword updates the user's password in the database, if the provided password is valid.
func (u User) UpdatePassword(password string, newPassword string) error {
	// check if oldPassword is valid
	err := comparePasswords(u.Password, password)
	if err != nil {
		return err
	}

	// generate hash of new password
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// store new password
	_, err = db.Exec("UPDATE users SET password = $1 WHERE id = $2;", newHashPassword, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateName updates the user's first name and last name.
func (u User) UpdateName(firstName string, lastName string) error {
	_, err := db.Exec("UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3;", firstName, lastName, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDirRoot updates the user's directory root.
func (u User) UpdateDirRoot(dirRoot string) error {
	if !u.Admin {
		return errors.New("User must be admin.")
	}

	if _, err := os.Stat(dirRoot); os.IsNotExist(err) {
		return errors.New("Directory does not exist.")
	}

	_, err := db.Exec("UPDATE users SET directory_root = $1 WHERE id = $2;", dirRoot, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// comparePasswords compares the hash string with a string to determine if it is equivalent.
func comparePasswords(hashPW string, pw string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPW), []byte(pw))
	if err != nil {
		return errors.New("Password is invalid.")
	}
	return nil
}
