// Package users contains types and methods for interacting with users stored in the database.
package users

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/robertjeffs/viewer-go/app/logic/config"
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

// DeleteUser deletes a user from the database if the provided password is valid.
func DeleteUser(db *sql.DB, u User, password string) error {
	// check if password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return NewErrInvalidPassword()
	}

	if _, err := db.Exec("DELETE FROM users WHERE id = $1", u.ID); err != nil {
		return err
	}
	return nil
}

// UpdateUserPassword updates the user's password in the database, if the provided password is valid.
func UpdateUserPassword(db *sql.DB, u User, password string, newPassword string) error {
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
	if _, err := db.Exec("UPDATE users SET password = $1 WHERE id = $2;", newHashPassword, u.ID); err != nil {
		return err
	}
	return nil
}

// UpdateUserFullname updates the user's first name and last name.
func UpdateUserFullname(db *sql.DB, u User, firstName string, lastName string) error {
	_, err := db.Exec("UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3;",
		firstName, lastName, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserDirRoot updates the user's directory root. Must be Admin.
func UpdateUserDirRoot(db *sql.DB, u User, dirRoot string) error {
	if !u.Admin {
		return errors.New("user must be admin")
	}

	if _, err := os.Stat(dirRoot); os.IsNotExist(err) {
		return errors.New("directory does not exist")
	}

	if _, err := db.Exec("UPDATE users SET directory_root = $1 WHERE id = $2;",
		dirRoot, u.ID); err != nil {
		return err
	}
	return nil
}

// GetAllUsers returns all users in the database.
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		u := User{}

		if err := rows.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Password, &u.DirectoryRoot,
			&u.Admin); err != nil {
			return []User{}, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUser returns a single user from the database that matches the provided id.
func GetUser(db *sql.DB, id int) (User, error) {
	var u User

	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	if err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Password, &u.DirectoryRoot,
		&u.Admin); err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("there is no user by that ID")
		}
	}
	return u, nil
}

// CreateUser inserts a new user into the database.
func CreateUser(db *sql.DB, u User) error {
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
	userDirectory := filepath.Join(config.GetUsersDirectory(),
		filepath.FromSlash(path.Clean("/"+u.DirectoryRoot)))
	if err := os.MkdirAll(userDirectory, os.ModePerm); err != nil {
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
func ValidateUser(db *sql.DB, username string, password string) (int, error) {
	var u User
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	if err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Password, &u.DirectoryRoot,
		&u.Admin); err != nil {
		return -1, errors.New("there is no user by that username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return -1, NewErrInvalidPassword()
	}
	return u.ID, nil
}

// ChangeUserUsername updates the user's username.
func ChangeUserUsername(db *sql.DB, username string, newUsername string) error {
	// check if username exists
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username)
	row.Scan(&count)
	if count != 1 {
		return errors.New("username does not exist")
	}

	if _, err := db.Exec("UPDATE users SET username = $1 WHERE username = $2;",
		newUsername, username); err != nil {
		return err
	}
	return nil
}

// ChangeUserAdminStatus updates the user's admin status.
func ChangeUserAdminStatus(db *sql.DB, id int, isAdmin bool) error {
	if err := checkUserExists(db, id); err != nil {
		return err
	}
	if _, err := db.Exec("UPDATE users SET admin = $1 WHERE id = $2", id, isAdmin); err != nil {
		return err
	}
	return nil
}

// DeleteUserByID deletes the user from the database that corresponds to the given ID.
func DeleteUserByID(db *sql.DB, id int) error {
	if err := checkUserExists(db, id); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

// checkUserExists checks if user does exist.
func checkUserExists(db *sql.DB, id int) error {
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
