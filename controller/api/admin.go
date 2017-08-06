// This file contains handlers for admin api routes.

package api

import (
	"net/http"

	"encoding/json"

	"github.com/FriedPigeon/viewer-go/common"
	"github.com/FriedPigeon/viewer-go/database"
)

// CreateUser receives new user information via json and creates the user. Client must be admin.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}
	if !user.Admin {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	u := database.User{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// number cannot be 0 as validation will fail
	u.ID = 1
	if err := common.ValidateJsonInput(u); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = database.CreateUser(u)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully created user.")
}

// DeleteUser receives user information via json and deletes the user. Client must be admin.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}
	if !user.Admin {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		UserID int `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// check if user exists
	_, err = database.GetUser(data.UserID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = database.DeleteUser(data.UserID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted user.")
}

// ChangeDirRoot receives new directory root via json and updates it in the database.
func ChangeDirRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}
	if !user.Admin {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		DirRoot string `json:"dir_root"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = user.UpdateDirRoot(data.DirRoot)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Changed directory root successfully.")
}

// ChangeUsername receives new username root via json and updates it in the database. Client must be admin.
func ChangeUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}
	if !user.Admin {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		CurrentUsername string `json:"current_username"`
		NewUsername     string `json:"new_username"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = database.ChangeUserUsername(data.CurrentUsername, data.NewUsername)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Changed user's username successfully.")
}
