// This file contains handlers for viewer functionality api routes.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/robertjeffs/viewer-go/controller/common"
)

// CreateFolder creates a folder on the disk of the name of the form value "folder-name", then redirects to the viewer
// page at path provided in the query string "path".
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	folderPath := struct {
		Path string `json:"path"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&folderPath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := common.ValidateJsonInput(folderPath); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	dirPath := path.Join(user.DirectoryRoot, folderPath.Path)
	if err := createFolder(dirPath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not create directory: %s", err.Error()))
		return
	}
	sendSuccessResponse(w, "Successfully created folder.")
}

// Delete deletes a folder from the disk of the name of the form value "file-name", then redirects to the viewer
// page at path provided in the query string "path".
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		Path string `json:"path"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	filePath := path.Join(user.DirectoryRoot, data.Path)
	if err := deleteFile(filePath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted file/folder.")
}

// DeleteAll deletes the contents of a path from the disk of the query string value "path", then redirects to the viewer
// page at that path.
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	data := struct {
		Path string `json:"path"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := common.ValidateJsonInput(data); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	dirPath := path.Join(user.DirectoryRoot, data.Path)
	if err := deleteAllFiles(dirPath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted all contents.")
}

// Upload parses a multipart form and saves uploaded files to the disk at the path from query string "path", then
// redirects to the viewer page at that path.
func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := common.ValidateUser(r)
	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized.")
		return
	}

	// parse request
	const twentyFourK = (1 << 10) * 24
	if err := r.ParseMultipartForm(twentyFourK); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	dirPath := path.Join(user.DirectoryRoot, r.FormValue("path"))
	if err = saveFiles(dirPath, r.MultipartForm.File); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "File upload success.")
}
