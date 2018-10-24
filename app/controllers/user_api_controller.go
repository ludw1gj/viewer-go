// This file contains handlers for user specific api routes.

package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ludw1gj/viewer-go/app/users"

	"fmt"

	"github.com/ludw1gj/viewer-go/app/logic/session"
)

// UserAPIController contains methods for user api route responses.
type UserAPIController struct {
	db      *sql.DB
	session *session.Manager
}

// Login will process a user login.
func (uc UserAPIController) Login(w http.ResponseWriter, r *http.Request) {
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
	if err := validateJSONInput(loginCredentials); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := users.ValidateUser(uc.db, loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := uc.session.NewUserSession(w, r, userID); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	sendSuccessResponse(w, "Login successful.")
}

// Logout will logout the user by changing the session value "authenticated" to false.
func (uc UserAPIController) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := uc.session.RemoveUserAuthFromSession(w, r); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Logout error: %s", err.Error()))
		return
	}
	sendSuccessResponse(w, "Logout successful.")
}

// DeleteAccount will process the delete user form, if password is correct the user's account will be deleted and
// the user redirected to the login page.
func (uc UserAPIController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := uc.session.ValidateUserSession(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := validateJSONInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := users.DeleteUser(uc.db, user, data.Password); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Account deletion successful.")
}

// ChangePassword will process a json post request, comparing password sent with current password and if they match,
// the current password will be changed to the new password.
func (uc UserAPIController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := uc.session.ValidateUserSession(r)
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
	if err := validateJSONInput(passwords); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := users.UpdateUserPassword(uc.db, user, passwords.OldPassword, passwords.NewPassword); err != nil {
		switch err.(type) {
		case *users.ErrInvalidPassword:
			sendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		default:
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	sendSuccessResponse(w, "Password changed successfully.")
}

// ChangeName will change the user's first/last name.
func (uc UserAPIController) ChangeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := uc.session.ValidateUserSession(r)
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
	if err := validateJSONInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := users.UpdateUserFullname(uc.db, user, data.FirstName, data.LastName); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Name changed successfully.")
}
