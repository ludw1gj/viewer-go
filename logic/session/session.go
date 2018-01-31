// Package session contains logic concerning session and cookies.
package session

import (
	"net/http"

	"errors"
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/robertjeffs/viewer-go/models"
)

var store *sessions.CookieStore

// Load initialises CookieStore.
func Load(configJSONFile string) error {
	ck, err := loadCookieConfig(configJSONFile)
	if err != nil {
		return errors.New("Failed to initialise a CookieStore: " + err.Error())
	}
	store = sessions.NewCookieStore(ck.Cookie.AuthorisationKey, ck.Cookie.EncryptionKey)
	return nil
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
	if err := s.Save(r, w); err != nil {
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
	if err := s.Save(r, w); err != nil {
		return err
	}
	return nil
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

// ValidateUserSession checks if user's session is valid and then returns the user's information.
func ValidateUserSession(r *http.Request) (u models.User, err error) {
	userId, err := GetUserID(r)
	if err != nil {
		return u, err
	}

	u, err = models.GetUser(userId)
	if err != nil {
		return u, err
	}
	return u, nil
}

// ValidateAdminSession checks if the user is valid and is admin.
func ValidateAdminSession(r *http.Request) (u models.User, err error) {
	u, err = ValidateUserSession(r)
	if err != nil {
		return u, err
	}
	if !u.Admin {
		return u, errors.New("user is not an admin")
	}
	return u, nil
}
