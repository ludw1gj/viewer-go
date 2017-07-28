// Package session contains logic concerning session and cookies.
package session

import (
	"net/http"

	"errors"
	"fmt"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

// Load returns a new CookieStore.
func Load(configFile string) error {
	ck, err := loadCookieConfig(configFile)
	if err != nil {
		return err
	}
	store = sessions.NewCookieStore(ck.Cookie.AuthorisationKey, ck.Cookie.EncryptionKey)
	return nil
}

// CheckIfAuth checks if user is authenticated.
func CheckIfAuth(r *http.Request) bool {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		return false
	}

	// check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// user is not auth
		return false
	}
	return true
}

// NewUserSession creates a new user session and authenticates the user.
func NewUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		return errors.New(fmt.Sprintf("Cookie is invalid, clearing cookies may help. Error: \"%s\"", err.Error()))
	}

	// Set user as authenticated
	session.Values["id"] = userID
	session.Values["authenticated"] = true
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

// GetUserIDFromSession returns a user's associated with a session.
func GetUserIDFromSession(r *http.Request) (id int, err error) {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		return id, err
	}

	id, ok := session.Values["id"].(int)
	if !ok {
		return id, err
	}
	return id, nil
}
