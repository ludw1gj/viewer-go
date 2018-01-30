// This file contains handlers for user specific api routes.

package api

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/robertjeffs/viewer-go/logic/session"
	"github.com/robertjeffs/viewer-go/logic/validate"
	"github.com/robertjeffs/viewer-go/model/database"
)

// UserLogin will process a user login.
func UserLogin(w http.ResponseWriter, r *http.Request) {
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
	if err := decoder.Decode(&loginCredentials); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := validate.ValidateJsonInput(loginCredentials); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := database.ValidateUser(loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := session.NewUserSession(w, r, userID); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	sendSuccessResponse(w, "Login successful.")
}

// UserLogout will logout the user by changing the session value "authenticated" to false.
func UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := session.RemoveUserAuthFromSession(w, r); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Logout error: %s", err.Error()))
		return
	}
	sendSuccessResponse(w, "Logout successful.")
}

// UserDeleteAccount will process the delete user form, if password is correct the user's account will be deleted and
// the user redirected to the login page.
func UserDeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := validate.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	password := struct {
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := validate.ValidateJsonInput(password); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := user.Delete(password.Password); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Account deletion successful.")
}

// UserChangePassword will process a json post request, comparing password sent with current password and if they match,
// the current password will be changed to the new password.
func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := validate.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	passwords := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&passwords); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.ValidateJsonInput(passwords); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := user.UpdatePassword(passwords.OldPassword, passwords.NewPassword); err != nil {
		switch err.(type) {
		case *database.ErrInvalidPassword:
			sendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		default:
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	sendSuccessResponse(w, "Password changed successfully.")
}

// UserChangeName will change the user's first/last name.
func UserChangeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := validate.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := user.UpdateName(data.FirstName, data.LastName); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Name changed successfully.")
}
