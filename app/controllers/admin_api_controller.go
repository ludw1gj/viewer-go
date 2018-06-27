// This file contains handlers for admin api routes.

package controllers

import (
	"database/sql"
	"net/http"

	"github.com/robertjeffs/viewer-go/app/users"

	"github.com/robertjeffs/viewer-go/app/logic/session"

	"encoding/json"

	"fmt"
)

// AdminAPIController contains methods for admin api route responses.
type AdminAPIController struct {
	db      *sql.DB
	session *session.Manager
}

// CreateUser receives new user information via json and creates the users. Client must be admin.
func (ac AdminAPIController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := ac.session.ValidateAdminSession(r); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	user := users.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// number cannot be 0 as validation will fail
	user.ID = 1
	if err := validateJSONInput(user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := users.CreateUser(ac.db, user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully created users.")
}

// DeleteUser receives user information via json and deletes the users. Client must be admin.
func (ac AdminAPIController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := ac.session.ValidateAdminSession(r); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		UserID int `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := validateJSONInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := users.DeleteUserByID(ac.db, data.UserID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted users.")
}

// ChangeUserUsername changes a user's username. Client must be admin.
func (ac AdminAPIController) ChangeUserUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := ac.session.ValidateAdminSession(r); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		CurrentUsername string `json:"current_username"`
		NewUsername     string `json:"new_username"`
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

	if err := users.ChangeUserUsername(ac.db, data.CurrentUsername, data.NewUsername); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Changed user's username successfully.")
}

// ChangeUserAdminStatus changes a user's admin status via the provided ID and updates it in the
// database. Client must be admin.
func (ac AdminAPIController) ChangeUserAdminStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := ac.session.ValidateAdminSession(r); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		UserID  int  `json:"user_id"`
		IsAdmin bool `json:"is_admin"`
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

	if err := users.ChangeUserAdminStatus(ac.db, data.UserID, data.IsAdmin); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, fmt.Sprintf("Changed admin status of user of id %d to %t", data.UserID, data.IsAdmin))
}

// ChangeDirRoot changes the directory root of the client and updates it in the database. Client must be admin.
func (ac AdminAPIController) ChangeDirRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := ac.session.ValidateAdminSession(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		DirRoot string `json:"dir_root"`
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

	if err := users.UpdateUserDirRoot(ac.db, user, data.DirRoot); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Changed directory root successfully.")
}
