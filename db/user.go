package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int
	Username      string
	HashPassword  string
	DirectoryRoot string
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

func GetUser(id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.ID, &user.Username, &user.HashPassword, &user.DirectoryRoot)
	if err != nil {
		return user, err

	}
	return user, nil
}

func ValidateUser(username string, password string) (user User, auth bool) {
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.HashPassword, &user.DirectoryRoot)
	if err != nil {
		// TODO: Properly handle error
		log.Println(err)
	}

	// comparing password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password)); err != nil {
		return user, false
	}
	return user, true
}
