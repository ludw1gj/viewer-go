package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"fmt"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
	"golang.org/x/crypto/bcrypt"
)

// Login will process a user login.
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorJSON{"Method must be POST."})
		return
	}
	loginCredentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginCredentials)

	err = session.NewUserSession(w, r, loginCredentials.Username, loginCredentials.Password)
	if err == sql.ErrNoRows || err == bcrypt.ErrMismatchedHashAndPassword {
		err = errors.New("Invalid username or password.")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Login Successful."})
}

// Logout will logout the user by changing the session value "authenticated" to false.
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := session.RemoveUserAuthFromSession(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{fmt.Sprintf("Logout error: %s", err.Error())})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully logged out."})
}

// DeleteAccount will process the delete user form, if password is correct the user's account will be deleted and the
// user redirected to the login page.
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	password := struct {
		Password string `json:"password"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	err = db.DeleteUserPasswordValidated(user, password.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully logged out."})
}

// ChangePassword will process a json post request, comparing password sent with current password and if they match, the
// current password will be changed to the new password.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	passwords := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&passwords)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	err = db.ChangeUserPassword(user, passwords.OldPassword, passwords.NewPassword)
	if err != nil {
		if err.Error() == "Incorrect password." {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorJSON{err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Password changed successfully."})
}
