// This file contains functions for rendering standard site pages.

package controllers

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"bytes"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/robertjeffs/viewer-go/logic/session"
	"github.com/robertjeffs/viewer-go/logic/templates"
	"github.com/robertjeffs/viewer-go/models"
)

type SiteController struct{}

func NewSiteController() *SiteController {
	return &SiteController{}
}

// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User models.User
}

// GetErrorPage renders the error page and sends status 500.
func (SiteController) GetErrorPage(w http.ResponseWriter, r *http.Request, pageErr error) {
	w.WriteHeader(http.StatusInternalServerError)
	user, err := session.ValidateUserSession(r)
	if err != nil {
		log.Printf("StatusInternalServerError failed to execute get user from session on error page: %s", err.Error())

		resp := fmt.Sprintf("500: Server error. Two errors have occured.<br>First Error: %s<br>Second Error: %s",
			pageErr.Error(), err.Error())
		w.Write([]byte(resp))
		return
	}

	data := struct {
		Error string
		User  models.User
	}{
		pageErr.Error(),
		user,
	}
	templates.RenderTemplate(w, "error", data)
}

// GetViewerPage handles the viewer page. It uses the path variable in the route to determine which directory in the
// user's directory in the filesystem to display a directory list for.
func (SiteController) GetViewerPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// urlPath should not contain a leading /
	urlPath := strings.TrimPrefix(mux.Vars(r)["path"], "/")
	data := struct {
		CurrentDir string
		User       models.User
	}{
		urlPath,
		user,
	}
	templates.RenderTemplate(w, "viewer", data)
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
func (SiteController) GetUserPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderTemplate(w, "user", userInfo{user})
}

// GetAboutPage handles the about page.
func (SiteController) GetAboutPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderTemplate(w, "about", userInfo{user})
}

// GetNotFoundPage renders the not found page and sends status 404.
func (SiteController) GetNotFoundPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderTemplate(w, "notFound", userInfo{user})
}

// SendFile sends file to client.
func (sc SiteController) SendFile(w http.ResponseWriter, r *http.Request) {
	// get user from session
	user, err := session.ValidateUserSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// path to file
	filePath := path.Join(user.DirectoryRoot, mux.Vars(r)["path"])

	// get file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		sc.GetErrorPage(w, r, err)
		return
	}
	if fileInfo.IsDir() {
		sc.GetErrorPage(w, r, errors.New("requested item is not a file"))
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		sc.GetErrorPage(w, r, err)
		return
	}

	// contentType determines the content-type by the file extension of the file at the path.
	contentType := func(path string) (contentType string) {
		hasSuffix := func(suffix string) bool {
			return strings.HasSuffix(path, suffix)
		}

		if hasSuffix(".css") {
			return "text/css"
		} else if hasSuffix(".js") {
			return "application/javascript"
		} else if hasSuffix(".png") {
			return "images/png"
		} else if hasSuffix(".jpg") {
			return "images/jpeg"
		} else if hasSuffix(".jpeg") {
			return "images/jpeg"
		} else if hasSuffix(".mp4") {
			return "video/mp4"
		}
		return "text/plain"
	}

	w.Header().Add("Content-Type", contentType(filePath))
	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(data))
}

// GetAdminPage renders the Administration page. Client must be admin.
func (SiteController) GetAdminPage(w http.ResponseWriter, r *http.Request) {
	u, err := session.ValidateAdminSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderTemplate(w, "admin", userInfo{u})
}

// GetAdminDisplayAllUsers render a sub administration page which displays all users in models. Client must be admin.
func (sc SiteController) GetAdminDisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := session.ValidateAdminSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		sc.GetErrorPage(w, r, err)
		return
	}

	data := struct {
		User  models.User
		Users []models.User
	}{u, users}
	templates.RenderTemplate(w, "adminUsers", data)
}
