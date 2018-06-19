// This file contains handlers for user specific api routes.

package controllers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/robertjeffs/viewer-go/app/logic/session"
	"github.com/robertjeffs/viewer-go/app/models"
)

type UserAPIController struct {
	*session.SessionManager
}

func NewUserAPIController() *UserAPIController {
	return &UserAPIController{
		session.NewSessionManager(),
	}
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

	userManager := models.NewUserManager()
	userID, err := userManager.ValidateUser(loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := uc.NewUserSession(w, r, userID); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	sendSuccessResponse(w, "Login successful.")
}

// Logout will logout the user by changing the session value "authenticated" to false.
func (uc UserAPIController) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := uc.RemoveUserAuthFromSession(w, r); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Logout error: %s", err.Error()))
		return
	}
	sendSuccessResponse(w, "Logout successful.")
}

// DeleteAccount will process the delete user form, if password is correct the user's account will be deleted and
// the user redirected to the login page.
func (uc UserAPIController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := uc.ValidateUserSession(r)
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
	if err := validateJSONInput(password); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := user.Delete(password.Password); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Account deletion successful.")
}

// ChangePassword will process a json post request, comparing password sent with current password and if they match,
// the current password will be changed to the new password.
func (uc UserAPIController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := uc.ValidateUserSession(r)
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

	if err := user.UpdatePassword(passwords.OldPassword, passwords.NewPassword); err != nil {
		switch err.(type) {
		case *models.ErrInvalidPassword:
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

	user, err := uc.ValidateUserSession(r)
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

	if err := user.UpdateName(data.FirstName, data.LastName); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Name changed successfully.")
}
