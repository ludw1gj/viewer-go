package handler

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"

	"io/ioutil"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// RedirectToViewer redirects users to the
func RedirectToViewer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, viewerRootURL, http.StatusMovedPermanently)
}

// Viewer handles the viewer page. It uses the path variable in the route to determine which directory of the filesystem
// to display a directory list for.
func Viewer(w http.ResponseWriter, r *http.Request) {
	isFile, err := renderIfFile(w, r)
	if err != nil {
		renderErrorPage(w, "", errors.New("There has been an error: "+err.Error()))
		return
	} else if isFile {
		return
	}

	path := mux.Vars(r)["path"]
	list, err := getDirectoryList(path)
	if err != nil {
		renderErrorPage(w, "", errors.New("There has been an error getting directory list: "+err.Error()))
	}
	data := struct {
		List       template.HTML
		CurrentDir string
	}{
		list,
		path,
	}
	err = viewerTpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// About handles the about page.
func About(w http.ResponseWriter, _ *http.Request) {
	aboutTpl.Execute(w, nil)
}

// Upload parses a multipart form and saves uploaded files to the disk at the path from query string "path", then
// redirects to the viewer page at that path.
func Upload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	// parse request
	const _24K = (1 << 10) * 24
	err := r.ParseMultipartForm(_24K)
	if err != nil {
		renderErrorPage(w, path, err)
		return
	}

	err = processMultipartFileHeaders(path, r.MultipartForm.File)
	if err != nil {
		renderErrorPage(w, path, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, viewerRootURL+path, http.StatusMovedPermanently)
}

// CreateFolder creates a folder on the disk of the name of the form value "folder-name", then redirects to the viewer
// page at path provided in the query string "path".
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	folderName := r.FormValue("folder-name")
	folderPath := path + "/" + folderName

	err := createFolder(folderPath)
	if err != nil {
		renderErrorPage(w, path, errors.New("Could not create directory: "+err.Error()))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, viewerRootURL+path, http.StatusMovedPermanently)
}

// Delete deletes a folder from the disk of the name of the form value "file-name", then redirects to the viewer
// page at path provided in the query string "path".
func Delete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	fileName := r.FormValue("file-name")
	if fileName == "" {
		renderErrorPage(w, path, errors.New("File name cannot be empty."))
	}

	err := deleteFile(path, fileName)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, viewerRootURL+path, http.StatusMovedPermanently)
	}
	http.Redirect(w, r, viewerRootURL+path, http.StatusMovedPermanently)
}

// DeleteAll deletes the contents of a path from the disk of the query string value "path", then redirects to the viewer
// page at that path.
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	err := deleteAllFiles(path)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, viewerRootURL+path, http.StatusMovedPermanently)
	}
	http.Redirect(w, r, viewerRootURL+path, http.StatusMovedPermanently)
}

// NotFound renders the not found page and sends status 404.
func NotFound(w http.ResponseWriter, _ *http.Request) {
	var buf bytes.Buffer
	err := notFoundTpl.Execute(&buf, nil)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(buf.Bytes())
}

// NotFound renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, path string, err error) {
	page := viewerRootURL + path

	data := struct {
		Error error
		Page  string
	}{
		err,
		page,
	}
	var buf bytes.Buffer
	err = errorTpl.Execute(&buf, data)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(buf.Bytes())
}

// renderIfFile uses the path variable in the route to determine if path on disk is a file or a directory. If it is a
// file it will write the file to the client, but if it is a directory it will return isFile is false.
func renderIfFile(w http.ResponseWriter, r *http.Request) (isFile bool, err error) {
	path := mux.Vars(r)["path"]
	filePath := wrkDir + path

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if !fileInfo.IsDir() {
		isFile = true

		data, _ := ioutil.ReadFile(filePath)
		w.Header().Add("Content-Type", contentType(filePath))
		http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(data))
		return
	}
	return
}
