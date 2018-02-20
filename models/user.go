package models

import (
	"errors"
	"os"

	"github.com/robertjeffs/viewer-go/logic/database"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidPassword type conforms to error type.
type ErrInvalidPassword struct {
	message string
}

// Error returns the error message.
func (e *ErrInvalidPassword) Error() string {
	return e.message
}

// NewErrInvalidPassword returns a pointer to a ErrInvalidPassword type instance.
func NewErrInvalidPassword() *ErrInvalidPassword {
	return &ErrInvalidPassword{
		message: "Password is invalid.",
	}
}

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
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return NewErrInvalidPassword()
	}

	if _, err := database.DB.Exec("DELETE FROM users WHERE id = $1", u.ID); err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates the user's password in the database, if the provided password is valid.
func (u User) UpdatePassword(password string, newPassword string) error {
	// check if oldPassword is valid
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return NewErrInvalidPassword()
	}

	// generate hash of new password
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// store new password
	if _, err := database.DB.Exec("UPDATE users SET password = $1 WHERE id = $2;", newHashPassword, u.ID); err != nil {
		return err
	}
	return nil
}

// UpdateName updates the user's first name and last name.
func (u User) UpdateName(firstName string, lastName string) error {
	_, err := database.DB.Exec("UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3;", firstName, lastName, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDirRoot updates the user's directory root. Must be Admin.
func (u User) UpdateDirRoot(dirRoot string) error {
	if !u.Admin {
		return errors.New("user must be admin")
	}

	if _, err := os.Stat(dirRoot); os.IsNotExist(err) {
		return errors.New("directory does not exist")
	}

	if _, err := database.DB.Exec("UPDATE users SET directory_root = $1 WHERE id = $2;", dirRoot, u.ID); err != nil {
		return err
	}
	return nil
}
