package db

import (
	"log"

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
	IsAdmin       bool   `json:"is_admin"`
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

		err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot, &user.IsAdmin)
		if err != nil {
			log.Println(err)
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// CreateUser inserts a new user into the database.
func CreateUser(u User) error {
	// generate hash of user password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	// store user in db
	_, err = db.Exec("INSERT INTO users (username, first_name, last_name, hash_password, directory_root, is_admin) VALUES ($1, $2, $3, $4, $5, $6)",
		u.Username, u.FirstName, u.LastName, string(hashPassword), u.DirectoryRoot, u.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database that matches the provided id.
func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserPasswordValidated deletes a user from the database if the provided password is valid.
func DeleteUserPasswordValidated(user User, password string) error {
	// check if password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return err
	}

	_, err := db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns a single user from the database that matches the provided id.
func GetUser(id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot, &user.IsAdmin)
	if err != nil {
		return user, err

	}
	return user, nil
}

// ChangeUserPassword changes a user's password in the database, if the provided password is valid.
func ChangeUserPassword(user User, oldPassword string, newPassword string) error {
	// check if oldPassword is valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("Incorrect password.")
	}

	// generate hash of new password
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	// store new password
	_, err = db.Exec("UPDATE users SET hash_password = $1 WHERE id = $2;", newHashPassword, user.ID)
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
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.DirectoryRoot, &user.IsAdmin)
	if err != nil {
		return userID, errors.New("Invalid username.")
	}

	// comparing password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return userID, errors.New("Invalid password.")
	}
	return user.ID, err
}
