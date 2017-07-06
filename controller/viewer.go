// Package controller contains handler functions for routes, and also functions which write to a http.Response.

package controller

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"

	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
	"github.com/gorilla/mux"
)

const viewerRootURL = "/viewer/"

type viewerController struct{}

func NewViewerController() *viewerController {
	return &viewerController{}
}

// Viewer handles the viewer page. It uses the path variable in the route to determine which directory in the user's
// directory in the filesystem to display a directory list for.
func (vc viewerController) Viewer(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isFile, err := vc.renderIfFile(w, r, user)
	if err != nil {
		renderErrorPage(w, r, errors.New("There has been an error: "+err.Error()))
		return
	} else if isFile {
		return
	}

	urlPath := mux.Vars(r)["path"]
	dirPath := path.Join(user.DirectoryRoot, urlPath)
	list, err := getDirectoryList(dirPath, urlPath)
	if err != nil {
		renderErrorPage(w, r, errors.New("There has been an error getting directory list: "+err.Error()))
	}
	data := struct {
		List       template.HTML
		CurrentDir string
		User       db.User
	}{
		list,
		urlPath,
		user,
	}
	err = viewerTpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// Upload parses a multipart form and saves uploaded files to the disk at the path from query string "path", then
// redirects to the viewer page at that path.
func (viewerController) Upload(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// parse request
	const _24K = (1 << 10) * 24
	err = r.ParseMultipartForm(_24K)
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}

	urlPath := r.URL.Query().Get("path")
	dirPath := path.Join(user.DirectoryRoot, urlPath)
	err = uploadFiles(dirPath, r.MultipartForm.File)
	if err != nil {
		renderErrorPage(w, r, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, viewerRootURL+urlPath, http.StatusSeeOther)
}

// CreateFolder creates a folder on the disk of the name of the form value "folder-name", then redirects to the viewer
// page at path provided in the query string "path".
func (viewerController) CreateFolder(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	urlPath := r.URL.Query().Get("path")
	dirPath := path.Join(user.DirectoryRoot, urlPath, r.FormValue("folder-name"))
	err = createFolder(dirPath)
	if err != nil {
		renderErrorPage(w, r, errors.New("Could not create directory: "+err.Error()))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, viewerRootURL+urlPath, http.StatusSeeOther)
}

// Delete deletes a folder from the disk of the name of the form value "file-name", then redirects to the viewer
// page at path provided in the query string "path".
func (viewerController) Delete(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fileName := r.FormValue("file-name")
	if fileName == "" {
		renderErrorPage(w, r, errors.New("File name cannot be empty."))
	}

	urlPath := r.URL.Query().Get("path")
	filePath := path.Join(user.DirectoryRoot, urlPath, fileName)
	err = deleteFile(filePath)
	if err != nil {
		renderErrorPage(w, r, err)
	}
	http.Redirect(w, r, viewerRootURL+urlPath, http.StatusSeeOther)
}

// DeleteAll deletes the contents of a path from the disk of the query string value "path", then redirects to the viewer
// page at that path.
func (viewerController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	urlPath := r.URL.Query().Get("path")
	dirPath := path.Join(user.DirectoryRoot, urlPath)
	err = deleteAllFiles(dirPath)
	if err != nil {
		renderErrorPage(w, r, err)
	}
	http.Redirect(w, r, viewerRootURL+urlPath, http.StatusSeeOther)
}

// renderIfFile uses the path variable in the route to determine if path on disk is a file or a directory. If it is a
// file it will write the file to the client, but if it is a directory it will return isFile is false.
func (vc viewerController) renderIfFile(w http.ResponseWriter, r *http.Request, user db.User) (isFile bool, err error) {
	filePath := path.Join(user.DirectoryRoot, mux.Vars(r)["path"])

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}

	if !fileInfo.IsDir() {
		isFile = true

		data, _ := ioutil.ReadFile(filePath)
		w.Header().Add("Content-Type", vc.contentType(filePath))
		http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(data))
		return
	}
	return
}

// contentType determines the content-type by the file extension of the file at the path.
func (viewerController) contentType(path string) (contentType string) {
	if strings.HasSuffix(path, ".css") {
		return "text/css"
	} else if strings.HasSuffix(path, ".html") {
		return "text/html"
	} else if strings.HasSuffix(path, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(path, ".png") {
		return "image/png"
	} else if strings.HasSuffix(path, ".jpg") {
		return "image/jpeg"
	} else if strings.HasSuffix(path, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(path, ".mp4") {
		return "video/mp4"
	}
	return "text/plain"
}
