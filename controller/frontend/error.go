// This file contains utility functions and types used by controllers.

package frontend

import (
	"bytes"
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

// NotFound renders the not found page and sends status 404.
func NotFound(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var buf bytes.Buffer
	err = notFoundTpl.ExecuteTemplate(&buf, "base.gohtml", userInfo{user})
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(buf.Bytes())
}

// renderErrorPage renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, r *http.Request, pageErr error) {
	w.WriteHeader(http.StatusInternalServerError)
	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.Write([]byte("500: Server error"))
		log.Print("StatusInternalServerError failed to execute get user from session on error page.")
		return
	}

	data := struct {
		Error string
		User  db.User
	}{
		pageErr.Error(),
		user,
	}

	var buf bytes.Buffer
	err = errorTpl.ExecuteTemplate(&buf, "base.gohtml", data)
	if err != nil {
		w.Write([]byte("500: Server error"))
		log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())
		return
	}

	w.Write(buf.Bytes())
}
