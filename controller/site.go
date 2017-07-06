// This files contains the site controller, which contains methods for rendering standard pages.

package controller

import (
	"bytes"
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/session"
)

type siteController struct{}

func NewSiteController() *siteController {
	return &siteController{}
}

// RedirectToViewer redirects users to the viewer page.
func (siteController) RedirectToViewer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, viewerRootURL, http.StatusMovedPermanently)
}

// About handles the about page.
func (siteController) About(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = aboutTpl.Execute(w, userInfo{user})
	if err != nil {
		log.Println(err)
	}
}

// NotFound renders the not found page and sends status 404.
func (siteController) NotFound(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var buf bytes.Buffer
	err = notFoundTpl.Execute(&buf, userInfo{user})
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(buf.Bytes())
}
