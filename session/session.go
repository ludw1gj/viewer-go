// Package session contains logic concerning session and cookies.
package session

import (
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/gorilla/sessions"
)

const cookieStoreAuthKey = "something-very-secret"

var store = sessions.NewCookieStore([]byte(cookieStoreAuthKey))

// CheckIfAuth checks if user is authenticated.
func CheckIfAuth(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		log.Println(err)
	}

	// check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// user is not auth, send to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	return true
}

// NewUserSession creates a new user session and authenticates the user.
func NewUserSession(w http.ResponseWriter, r *http.Request, username string, password string) error {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		return err
	}

	// validate and get user's id
	userID, err := db.CheckUserValidation(username, password)
	if err != nil {
		return err
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["id"] = userID
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUserAuthFromSession removes user's session authentication.
func RemoveUserAuthFromSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		return err
	}

	// revoke user's authentication
	session.Values["authenticated"] = false
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// GetUserFromSession returns a user identified by the user's session.
func GetUserFromSession(r *http.Request) (user db.User, err error) {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		log.Println(err)
	}

	id := session.Values["id"].(int)
	user, err = db.GetUser(id)
	if err != nil {
		return user, err
	}
	return user, nil
}