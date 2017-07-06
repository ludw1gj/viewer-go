// This file contains utility functions and types used by controllers.

package controller

import (
	"bytes"
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User db.User
}

// errType is used for data object of error for rendering templates.
type errType struct {
	Error error
}

// contentJSON is used for json response with generic key of "content".
type contentJSON struct {
	Content string `json:"content"`
}

// errorJSON is used for json response with the key of "error".
type errorJSON struct {
	Error string `json:"error"`
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
	err = errorTpl.Execute(&buf, data)
	if err != nil {
		w.Write([]byte("500: Server error"))
		log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())
		return
	}

	w.Write(buf.Bytes())
}
