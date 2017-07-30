// Package session contains logic concerning session and cookies.
package session

import (
	"net/http"

	"errors"
	"fmt"

	"log"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

// init initialises CookieStore.
func init() {
	ck, err := loadCookieConfig("config.json")
	if err != nil {
		log.Fatalln("Failed to initialise a CookieStore:", err.Error())
	}
	store = sessions.NewCookieStore(ck.Cookie.AuthorisationKey, ck.Cookie.EncryptionKey)
}

// CheckUserAuth checks if user is authenticated.
func CheckUserAuth(r *http.Request) bool {
	s, err := store.Get(r, "viewer-session")
	if err != nil {
		return false
	}

	// check if user is authenticated
	if auth, ok := s.Values["authenticated"].(bool); !ok || !auth {
		// user is not auth
		return false
	}
	return true
}

// NewUserSession creates a new user session and authenticates the user.
func NewUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	s, err := store.Get(r, "viewer-session")
	if err != nil {
		return errors.New(fmt.Sprintf("Cookie is invalid, clearing cookies may help. Error: \"%s\"", err.Error()))
	}

	// Set user as authenticated
	s.Values["id"] = userID
	s.Values["authenticated"] = true
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUserAuthFromSession removes user's session authentication.
func RemoveUserAuthFromSession(w http.ResponseWriter, r *http.Request) error {
	s, err := store.Get(r, "viewer-session")
	if err != nil {
		return err
	}

	// revoke user's authentication
	s.Values["authenticated"] = false
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// GetUserID returns a user's associated with a session.
func GetUserID(r *http.Request) (id int, err error) {
	s, err := store.Get(r, "viewer-session")
	if err != nil {
		return id, err
	}

	id, ok := s.Values["id"].(int)
	if !ok {
		return id, err
	}
	return id, nil
}
