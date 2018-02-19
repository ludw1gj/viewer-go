package models

import (
	"golang.org/x/crypto/bcrypt"
)

// comparePasswords compares the hash string with a string to determine if it is equivalent.
func comparePasswords(hashPW string, pw string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPW), []byte(pw)); err != nil {
		return NewErrInvalidPassword("Password is invalid.")
	}
	return nil
}
