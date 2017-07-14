// This file contains handlers for viewer functionality api routes.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"net/url"
	"os"
	"sort"

	"github.com/FriedPigeon/viewer-go/session"
)

// TODO: doc here

func NewGetDirectoryList(w http.ResponseWriter, r *http.Request) {
	// get user from session
	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	// init dirList
	var list dirList

	// get current system dir and append to type dirList.CurrentDir
	data := struct {
		Path string `json:"path"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if data.Path == "/" {
		data.Path = ""
	}
	list.CurrentDir = data.Path
	currentSystemDir := path.Join(user.DirectoryRoot, data.Path)

	// get items in current system dir
	f, err := os.Open(currentSystemDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}
	defer f.Close()
	items, err := f.Readdir(-1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	// sort items by name
	itemInfo := make(map[string]bool)
	for _, item := range items {
		itemInfo[item.Name()] = !item.IsDir()
	}
	itemNamesSorted := make([]string, len(itemInfo))
	i := 0
	for itemName := range itemInfo {
		itemNamesSorted[i] = itemName
		i++
	}
	sort.Strings(itemNamesSorted)

	// append items to type dirList.Items, determining if item is a file or directory
	for _, itemName := range itemNamesSorted {
		isFile := itemInfo[itemName]

		var rawUrl string
		switch isFile {
		case true:
			rawUrl = "/file" + data.Path + "/" + itemName
		case false:
			rawUrl = data.Path + "/" + itemName
			itemName = itemName + "/"
		}
		urlLink, err := url.Parse(rawUrl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorJSON{err.Error()})
			return
		}
		list.Items = append(list.Items, item{itemName, isFile, urlLink.String()})
	}

	// send json response of type dirList
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// CreateFolder creates a folder on the disk of the name of the form value "folder-name", then redirects to the viewer
// page at path provided in the query string "path".
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	folderPath := struct {
		Path string `json:"path"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&folderPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	dirPath := path.Join(user.DirectoryRoot, folderPath.Path)
	err = createFolder(dirPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{fmt.Sprintf("Could not create directory: %s", err.Error())})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully created folder."})
}

// Delete deletes a folder from the disk of the name of the form value "file-name", then redirects to the viewer
// page at path provided in the query string "path".
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	data := struct {
		Path string `json:"path"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	filePath := path.Join(user.DirectoryRoot, data.Path)
	err = deleteFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully deleted file/folder."})
}

// DeleteAll deletes the contents of a path from the disk of the query string value "path", then redirects to the viewer
// page at that path.
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	data := struct {
		Path string `json:"path"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	dirPath := path.Join(user.DirectoryRoot, data.Path)
	err = deleteAllFiles(dirPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully deleted all contents."})
}

// Upload parses a multipart form and saves uploaded files to the disk at the path from query string "path", then
// redirects to the viewer page at that path.
func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	// parse request
	const _24K = (1 << 10) * 24
	err = r.ParseMultipartForm(_24K)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	dirPath := path.Join(user.DirectoryRoot, r.FormValue("path"))
	err = uploadFiles(dirPath, r.MultipartForm.File)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"File upload success."})
}
