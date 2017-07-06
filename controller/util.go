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
func renderErrorPage(w http.ResponseWriter, r *http.Request, err error) {
	user, ok := session.GetUserFromSession(r)
	if ok != nil {
		log.Println(ok)
		// TODO: handle error properly
		return
	}

	data := struct {
		Error string
		User  db.User
	}{
		err.Error(),
		user,
	}

	var buf bytes.Buffer
	err = errorTpl.Execute(&buf, data)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(buf.Bytes())
}
