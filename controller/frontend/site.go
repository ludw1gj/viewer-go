// This file contains functions for rendering standard site pages.

package frontend

import (
	"errors"
	"log"
	"net/http"

	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
	"github.com/gorilla/mux"
)

// ViewerPage handles the viewer page. It uses the path variable in the route to determine which directory in the user's
// directory in the filesystem to display a directory list for.
func ViewerPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// urlPath should not contain a leading /
	urlPath := strings.TrimPrefix(mux.Vars(r)["path"], "/")
	data := struct {
		CurrentDir string
		User       db.User
	}{
		urlPath,
		user,
	}
	renderTemplate(w, r, viewerTpl, data)
}

// SendFile sends file to client.
func SendFile(w http.ResponseWriter, r *http.Request) {
	// get user from session
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// path to file
	filePath := path.Join(user.DirectoryRoot, mux.Vars(r)["path"])

	// get file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}
	if fileInfo.IsDir() {
		renderErrorPage(w, r, errors.New("Requested item is not a file."))
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}
	w.Header().Add("Content-Type", contentType(filePath))
	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(data))
}

// RedirectToViewer redirects users to the viewer page.
func RedirectToViewer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}

// LoginPage method renders the login page.
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	err := loginTpl.Execute(&buf, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Server error"))
		log.Println(err)
	}
	w.Write(buf.Bytes())
}

// AboutPage handles the about page.
func AboutPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, r, aboutTpl, userInfo{user})
}

// UserPage renders the user page.
func UserPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, r, userTpl, userInfo{user})
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
