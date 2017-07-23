// This file contains handlers for user specific api routes.

package api

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

// Login will process a user login.
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		sendErrorResponse(w, http.StatusBadRequest, "Method must be POST.")
		return
	}
	loginCredentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginCredentials)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	err = session.NewUserSession(w, r, loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error()) //"Unauthorized."
		return
	}
	sendSuccessResponse(w, "Login successful.")
}

// Logout will logout the user by changing the session value "authenticated" to false.
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := session.RemoveUserAuthFromSession(w, r)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Logout error: %s", err.Error()))
		return
	}
	sendSuccessResponse(w, "Logout successful.")
}

// DeleteAccount will process the delete user form, if password is correct the user's account will be deleted and the
// user redirected to the login page.
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	password := struct {
		Password string `json:"password"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = db.DeleteUserPasswordValidated(user, password.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Account deletion successful.")
}

// ChangePassword will process a json post request, comparing password sent with current password and if they match, the
// current password will be changed to the new password.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	passwords := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&passwords)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.ChangeUserPassword(user, passwords.OldPassword, passwords.NewPassword)
	if err != nil {
		if err.Error() == "Incorrect password." {
			sendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Password changed successfully.")
}
