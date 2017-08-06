// This file contains functions for rendering standard site pages.

package frontend

import (
	"log"
	"net/http"

	"bytes"

	"fmt"

	"github.com/FriedPigeon/viewer-go/common"
)

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

	var tplBuf bytes.Buffer
	err := loginTpl.Execute(&tplBuf, nil)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500: Server error. %s", err.Error())))
		return
	}
	w.Write(tplBuf.Bytes())
}

// UserPage renders the user page.
func UserPage(w http.ResponseWriter, r *http.Request) {
	user, err := common.ValidateUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, r, userTpl, userInfo{user})
}

// AboutPage handles the about page.
func AboutPage(w http.ResponseWriter, r *http.Request) {
	user, err := common.ValidateUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, r, aboutTpl, userInfo{user})
}
