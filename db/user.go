package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int
	Username      string
	FirstName     string
	LastName      string
	HashPassword  string
	DirectoryRoot string
	IsAdmin       string
}

func CreateUser(username string, password string) {
	// generate hash of user password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}

	// store user in db
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashPassword))
	if err != nil {
		// TODO: Properly handle error
		log.Println(err)
	}
}

func DeleteUser(user User, password string) error {
	// check if password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return err
	}

	_, err := db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.HashPassword, &user.DirectoryRoot, &user.IsAdmin)
	if err != nil {
		return user, err

	}
	return user, nil
}

func ChangeUserPassword(user User, oldPassword string, newPassword string) error {
	// check if oldPassword is valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(oldPassword)); err != nil {
		return err
	}

	// generate hash of new password
	newHashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}

	// store new password
	_, err = db.Exec("UPDATE users SET hash_password = $1 WHERE id = $2;", newHashPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserValidation(username string, password string) (userID int, err error) {
	var user User
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err = row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.HashPassword, &user.DirectoryRoot, &user.IsAdmin)
	if err != nil {
		return userID, err
	}

	// comparing password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return userID, err
	}
	return user.ID, err
}
