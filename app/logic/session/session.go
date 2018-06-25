// Package session contains logic concerning session and cookies.
package session

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/robertjeffs/viewer-go/app/users"

	"github.com/gorilla/sessions"
)

// Manager contains the cookie store and methods useless for managing the user session.
type Manager struct {
	store  *sessions.CookieStore
	loaded bool
}

// NewManager loads the cookie store, and returns a Manager instance and an error if one had
// occured.
func NewManager(configJSONFile string) (sm Manager, err error) {
	if sm.loaded {
		return sm, errors.New("CookieStore is already loaded")
	}
	store, err := generateCookieStore(configJSONFile)
	if err != nil {
		return sm, err
	}

	sm.store = store
	return sm, nil
}

// NewUserSession creates a new user session and authenticates the users.
func (sm Manager) NewUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	s, err := sm.store.Get(r, "viewer-session")
	if err != nil {
		if err.Error() == "securecookie: the value is not valid" && s != nil {
			sm.store.Save(r, w, s)
		} else {
			return fmt.Errorf("Cookie error. Error: \"%s\"", err.Error())
		}
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
func (sm Manager) RemoveUserAuthFromSession(w http.ResponseWriter, r *http.Request) error {
	s, err := sm.store.Get(r, "viewer-session")
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
func (sm Manager) CheckUserAuth(r *http.Request) bool {
	s, err := sm.store.Get(r, "viewer-session")
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

// ValidateUserSession checks if user's session is valid and then returns the user's information.
func (sm Manager) ValidateUserSession(r *http.Request) (user users.User, err error) {
	// getUserID returns a user's associated with a session.
	getUserID := func(r *http.Request) (id int, err error) {
		s, err := sm.store.Get(r, "viewer-session")
		if err != nil {
			return id, err
		}

		id, ok := s.Values["id"].(int)
		if !ok {
			return id, err
		}
		return id, nil
	}

	userID, err := getUserID(r)
	if err != nil {
		return user, err
	}

	user, err = users.GetUser(userID)
	if err != nil {
		return user, err
	}
	return user, nil
}

// ValidateAdminSession checks if the user is valid and is admin.
func (sm Manager) ValidateAdminSession(r *http.Request) (user users.User, err error) {
	user, err = sm.ValidateUserSession(r)
	if err != nil {
		return user, err
	}
	if !user.Admin {
		return user, errors.New("user is not an admin")
	}
	return user, nil
}
