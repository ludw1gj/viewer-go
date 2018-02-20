package session

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/robertjeffs/viewer-go/models"
)

type SessionManager struct{}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// NewUserSession creates a new user session and authenticates the user.
func (SessionManager) NewUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	s, err := store.Get(r, "viewer-session")
	if err != nil {
		if err.Error() == "securecookie: the value is not valid" && s != nil {
			store.Save(r, w, s)
		} else {
			return errors.New(fmt.Sprintf("Cookie error. Error: \"%s\"", err.Error()))
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
func (SessionManager) RemoveUserAuthFromSession(w http.ResponseWriter, r *http.Request) error {
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
func (SessionManager) CheckUserAuth(r *http.Request) bool {
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

// ValidateUserSession checks if user's session is valid and then returns the user's information.
func (SessionManager) ValidateUserSession(r *http.Request) (u models.User, err error) {
	// getUserID returns a user's associated with a session.
	getUserID := func(r *http.Request) (id int, err error) {
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

	userId, err := getUserID(r)
	if err != nil {
		return u, err
	}

	u, err = models.NewUserManager().GetUser(userId)
	if err != nil {
		return u, err
	}
	return u, nil
}

// ValidateAdminSession checks if the user is valid and is admin.
func (sm SessionManager) ValidateAdminSession(r *http.Request) (u models.User, err error) {
	u, err = sm.ValidateUserSession(r)
	if err != nil {
		return u, err
	}
	if !u.Admin {
		return u, errors.New("user is not an admin")
	}
	return u, nil
}
