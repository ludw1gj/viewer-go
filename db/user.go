package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Username     string
	HashPassword string
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

func ValidateUser(username string, password string) bool {
	var user User

	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.HashPassword)
	if err != nil {
		// TODO: Properly handle error
		log.Println(err)
	}

	// comparing password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
