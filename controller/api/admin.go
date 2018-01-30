// This file contains handlers for admin api routes.

package api

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/robertjeffs/viewer-go/logic/common"
	"github.com/robertjeffs/viewer-go/model/database"
)

// AdminCreateUser receives new user information via json and creates the user. Client must be admin.
func AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := common.ValidateAdmin(r); err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	user := database.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// number cannot be 0 as validation will fail
	user.ID = 1
	if err := common.ValidateJsonInput(user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.CreateUser(user); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully created user.")
}

// AdminDeleteUser receives user information via json and deletes the user. Client must be admin.
func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := common.ValidateAdmin(r); err != nil {
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
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteUser(data.UserID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted user.")
}

// AdminChangeUserUsername changes a user's username. Client must be admin.
func AdminChangeUserUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := common.ValidateAdmin(r); err != nil {
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
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.ChangeUserUsername(data.CurrentUsername, data.NewUsername); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Changed user's username successfully.")
}

// AdminChangeUserAdminStatus changes a user's admin status via the provided ID and updates it in the database. Client must
// be admin.
func AdminChangeUserAdminStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := common.ValidateAdmin(r); err != nil {
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
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.ChangeUserAdminStatus(data.UserID, data.IsAdmin); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, fmt.Sprintf("Changed admin status of user of id %d to %t", data.UserID, data.IsAdmin))
}

// AdminChangeDirRoot changes the directory root of the client and updates it in the database. Client must be admin.
func AdminChangeDirRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, err := common.ValidateAdmin(r)
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
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := u.UpdateDirRoot(data.DirRoot); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Changed directory root successfully.")
}
