package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

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
func NewErrInvalidPassword(message string) *ErrInvalidPassword {
	return &ErrInvalidPassword{
		message: message,
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
	if err := comparePasswords(u.Password, password); err != nil {
		return err
	}

	if _, err := database.DB.Exec("DELETE FROM users WHERE id = $1", u.ID); err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates the user's password in the database, if the provided password is valid.
func (u User) UpdatePassword(password string, newPassword string) error {
	// check if oldPassword is valid
	if err := comparePasswords(u.Password, password); err != nil {
		return err
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

// UpdateDirRoot updates the user's directory root.
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

// GetAllUsers returns all users in the database.
func GetAllUsers() (users []User, err error) {
	rows, err := database.DB.Query("SELECT * FROM users")
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
	row := database.DB.QueryRow("SELECT * FROM users WHERE id = $1", id)
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
	row := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", u.Username)
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return errors.New("username is taken")
	}

	// create user root directory on disk
	userDirectory := filepath.Join("users", filepath.FromSlash(path.Clean("/"+u.DirectoryRoot)))
	if err := os.MkdirAll(userDirectory, os.ModePerm); err != nil {
		return err
	}

	// generate hash of user password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// store user in database
	if _, err := database.DB.Exec("INSERT INTO users (username, first_name, last_name, password, directory_root, admin) "+
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
	row := database.DB.QueryRow("SELECT * FROM users WHERE username = $1", username)
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
	row := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username)
	row.Scan(&count)
	if count != 1 {
		return errors.New("username does not exist")
	}

	if _, err := database.DB.Exec("UPDATE users SET username = $1 WHERE username = $2;", newUsername, username); err != nil {
		return err
	}
	return nil
}

// ChangeUserAdminStatus updates the user's admin status.
func ChangeUserAdminStatus(id int, isAdmin bool) error {
	if err := checkUserExists(id); err != nil {
		return err
	}
	if _, err := database.DB.Exec("UPDATE users SET admin = $1 WHERE id = $2", id, isAdmin); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes the user from the database that corresponds to the given ID.
func DeleteUser(id int) error {
	if err := checkUserExists(id); err != nil {
		return err
	}
	if _, err := database.DB.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

// comparePasswords compares the hash string with a string to determine if it is equivalent.
func comparePasswords(hashPW string, pw string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPW), []byte(pw)); err != nil {
		return NewErrInvalidPassword("password is invalid")
	}
	return nil
}

// checkUserExists checks if user does exist.
func checkUserExists(id int) error {
	var count int
	row := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id)
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("user by id %d does not exist", id)
	}
	return nil
}
