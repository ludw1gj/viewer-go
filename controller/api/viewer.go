// This file contains handlers for viewer functionality api routes.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/robertjeffs/viewer-go/logic/common"
)

// ViewerCreateFolder creates a folder on the disk of the name of the form value "folder-name", then redirects to the
// viewer page at path provided in the query string "path".
func ViewerCreateFolder(w http.ResponseWriter, r *http.Request) {
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

	// createFolder creates a folder in the directory path.
	createFolder := func(dirPath string) error {
		if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
			return errors.New("folder already exists")
		}

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	dirPath := path.Join(user.DirectoryRoot, folderPath.Path)
	if err := createFolder(dirPath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not create directory: %s", err.Error()))
		return
	}
	sendSuccessResponse(w, "Successfully created folder.")
}

// ViewerDelete deletes a folder from the disk of the name of the form value "file-name", then redirects to the viewer
// page at path provided in the query string "path".
func ViewerDelete(w http.ResponseWriter, r *http.Request) {
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

	// deleteFile deletes the file at file path.
	deleteFile := func(filePath string) (err error) {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return errors.New("file or Folder does not exist")
		}

		if err := os.RemoveAll(filePath); err != nil {
			return err
		}
		return nil
	}

	filePath := path.Join(user.DirectoryRoot, data.Path)
	if err := deleteFile(filePath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted file/folder.")
}

// ViewerDeleteAll deletes the contents of a path from the disk of the query string value "path", then redirects to the
// viewer page at that path.
func ViewerDeleteAll(w http.ResponseWriter, r *http.Request) {
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

	// deleteAllFiles deletes all files in the directory oath.
	deleteAllFiles := func(dirPath string) (err error) {
		d, err := os.Open(dirPath)
		if err != nil {
			return err
		}
		defer d.Close()

		names, err := d.Readdirnames(-1)
		if err != nil {
			return err
		}
		for _, name := range names {
			if err := os.RemoveAll(filepath.Join(dirPath, name)); err != nil {
				return err
			}
		}
		return nil
	}

	dirPath := path.Join(user.DirectoryRoot, data.Path)
	if err := deleteAllFiles(dirPath); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "Successfully deleted all contents.")
}

// ViewerUpload parses a multipart form and saves uploaded files to the disk at the path from query string "path", then
// redirects to the viewer page at that path.
func ViewerUpload(w http.ResponseWriter, r *http.Request) {
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

	// saveFiles opens the FileHeader's associated Files, creates the destinations at the directory path and saves the
	// files.
	saveFiles := func(dirPath string, file map[string][]*multipart.FileHeader) error {
		for _, fileHeaders := range file {
			for _, hdr := range fileHeaders {
				// open files
				inFile, err := hdr.Open()
				if err != nil {
					return err
				}

				// open destination
				outFile, err := os.Create(dirPath + "/" + hdr.Filename)
				if err != nil {
					return err
				}

				// 32K buffer copy
				if _, err = io.Copy(outFile, inFile); err != nil {
					return err
				}
			}
		}
		return nil
	}

	dirPath := path.Join(user.DirectoryRoot, r.FormValue("path"))
	if err = saveFiles(dirPath, r.MultipartForm.File); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendSuccessResponse(w, "File upload success.")
}
