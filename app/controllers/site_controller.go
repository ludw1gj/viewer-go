// This file contains functions for rendering standard site pages.

package controllers

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ludw1gj/viewer-go/app/users"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/ludw1gj/viewer-go/app/logic/session"
	"github.com/ludw1gj/viewer-go/app/logic/templates"
)

// SiteController contains methods for general site route responses.
type SiteController struct {
	db      *sql.DB
	session *session.Manager
}

// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User users.User
}

// getErrorPage renders the error page and sends status 500.
func (sc SiteController) getErrorPage(w http.ResponseWriter, r *http.Request, pageErr error) {
	w.WriteHeader(http.StatusInternalServerError)
	user, err := sc.session.ValidateUserSession(r)
	if err != nil {
		log.Printf("StatusInternalServerError failed to execute get user from session on error page: %s", err.Error())

		resp := fmt.Sprintf("500: Server error. Two errors have occured.<br>First Error: %s<br>Second Error: %s",
			pageErr.Error(), err.Error())
		w.Write([]byte(resp))
		return
	}

	data := struct {
		Error string
		User  users.User
	}{
		pageErr.Error(),
		user,
	}
	templates.RenderSiteTemplate(w, "error", data)
}

// GetViewerPage handles the viewer page. It uses the path variable in the route to determine which
// directory in the user's directory in the filesystem to display a directory list for.
func (sc SiteController) GetViewerPage(w http.ResponseWriter, r *http.Request) {
	user, err := sc.session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// urlPath should not contain a leading /
	urlPath := strings.TrimPrefix(mux.Vars(r)["path"], "/")
	data := struct {
		CurrentDir string
		User       users.User
	}{
		urlPath,
		user,
	}
	templates.RenderSiteTemplate(w, "viewer", data)
}

// GetLoginPage method renders the login page.
func (SiteController) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	templates.RenderLoginTemplate(w)
}

// GetUserPage renders the user page.
func (sc SiteController) GetUserPage(w http.ResponseWriter, r *http.Request) {
	user, err := sc.session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderSiteTemplate(w, "user", userInfo{user})
}

// GetAboutPage handles the about page.
func (sc SiteController) GetAboutPage(w http.ResponseWriter, r *http.Request) {
	user, err := sc.session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderSiteTemplate(w, "about", userInfo{user})
}

// GetNotFoundPage renders the not found page and sends status 404.
func (sc SiteController) GetNotFoundPage(w http.ResponseWriter, r *http.Request) {
	user, err := sc.session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderSiteTemplate(w, "notFound", userInfo{user})
}

// SendFile sends file to client.
func (sc SiteController) SendFile(w http.ResponseWriter, r *http.Request) {
	// get user from session
	user, err := sc.session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// path to file
	filePath := cleanPath(user.DirectoryRoot, mux.Vars(r)["path"])

	// get file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		sc.getErrorPage(w, r, err)
		return
	}
	if fileInfo.IsDir() {
		sc.getErrorPage(w, r, errors.New("requested item is not a file"))
		return
	}

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		sc.getErrorPage(w, r, errors.New("requested file not found"))
		return
	}

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)

	fileStat, _ := file.Stat()
	FileSize := strconv.FormatInt(fileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+fileStat.Name())
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", FileSize)
	file.Seek(0, 0)
	io.Copy(w, file)
}

// GetAdminPage renders the Administration page. Client must be admin.
func (sc SiteController) GetAdminPage(w http.ResponseWriter, r *http.Request) {
	user, err := sc.session.ValidateAdminSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderSiteTemplate(w, "admin", userInfo{user})
}

// GetAdminDisplayAllUsers render a sub administration page which displays all users in database.
// Client must be admin.
func (sc SiteController) GetAdminDisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := sc.session.ValidateAdminSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	allUsers, err := users.GetAllUsers(sc.db)
	if err != nil {
		sc.getErrorPage(w, r, err)
		return
	}

	data := struct {
		User  users.User
		Users []users.User
	}{user, allUsers}
	templates.RenderSiteTemplate(w, "adminUsers", data)
}
