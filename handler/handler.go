// Package handler contains handler functions for routes, and also functions which write to a http.Response.
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

	"strings"

	"github.com/FriedPigeon/viewer-go/config"
	"github.com/FriedPigeon/viewer-go/session"
	"github.com/FriedPigeon/viewer-go/tpl"
	"github.com/gorilla/mux"
)

// RedirectToViewer redirects users to the viewer page.
func RedirectToViewer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.ViewerRootURL, http.StatusMovedPermanently)
}

// Viewer handles the viewer page. It uses the path variable in the route to determine which directory of the filesystem
// to display a directory list for.
func Viewer(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	isFile, err := renderIfFile(w, r)
	if err != nil {
		renderErrorPage(w, errors.New("There has been an error: "+err.Error()))
		return
	} else if isFile {
		return
	}

	path := mux.Vars(r)["path"]
	list, err := user.GetDirectoryList(path)
	if err != nil {
		renderErrorPage(w, errors.New("There has been an error getting directory list: "+err.Error()))
	}
	data := struct {
		List       template.HTML
		CurrentDir string
		Username   string
	}{
		list,
		path,
		user.Username,
	}
	err = tpl.ViewerTpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// About handles the about page.
func About(w http.ResponseWriter, r *http.Request) {
	//user, err := session.GetUserFromSession(r)
	//if err != nil {
	//	renderErrorPage(w, err)
	//	return
	//}

	tpl.AboutTpl.Execute(w, nil)
}

// Upload parses a multipart form and saves uploaded files to the disk at the path from query string "path", then
// redirects to the viewer page at that path.
func Upload(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	path := r.URL.Query().Get("path")

	// parse request
	const _24K = (1 << 10) * 24
	err = r.ParseMultipartForm(_24K)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	err = user.ProcessMultipartFileHeaders(path, r.MultipartForm.File)
	if err != nil {
		renderErrorPage(w, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, config.ViewerRootURL+path, http.StatusSeeOther)
}

// CreateFolder creates a folder on the disk of the name of the form value "folder-name", then redirects to the viewer
// page at path provided in the query string "path".
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	path := r.URL.Query().Get("path")
	folderName := r.FormValue("folder-name")
	folderPath := path + "/" + folderName

	err = user.CreateFolder(folderPath)
	if err != nil {
		renderErrorPage(w, errors.New("Could not create directory: "+err.Error()))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, config.ViewerRootURL+path, http.StatusSeeOther)
}

// Delete deletes a folder from the disk of the name of the form value "file-name", then redirects to the viewer
// page at path provided in the query string "path".
func Delete(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	path := r.URL.Query().Get("path")

	fileName := r.FormValue("file-name")
	if fileName == "" {
		renderErrorPage(w, errors.New("File name cannot be empty."))
	}

	err = user.DeleteFile(path, fileName)
	if err != nil {
		renderErrorPage(w, err)
	}
	http.Redirect(w, r, config.ViewerRootURL+path, http.StatusSeeOther)
}

// DeleteAll deletes the contents of a path from the disk of the query string value "path", then redirects to the viewer
// page at that path.
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	path := r.URL.Query().Get("path")
	err = user.DeleteAllFiles(path)
	if err != nil {
		renderErrorPage(w, err)
	}
	http.Redirect(w, r, config.ViewerRootURL+path, http.StatusSeeOther)
}

// NotFound renders the not found page and sends status 404.
func NotFound(w http.ResponseWriter, r *http.Request) {
	//user, err := session.GetUserFromSession(r)
	//if err != nil {
	//	renderErrorPage(w, err)
	//	return
	//}

	var buf bytes.Buffer
	err := tpl.NotFoundTpl.Execute(&buf, nil)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(buf.Bytes())
}

// NotFound renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, err error) {
	var buf bytes.Buffer
	err = tpl.ErrorTpl.Execute(&buf, err)
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
	filePath := config.WrkDir + path

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

// contentType determines the content-type by the file extension of the file at the path.
func contentType(path string) (contentType string) {
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
