package controller

import (
	"bytes"
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
)

type userInfo struct {
	User db.User
}

type errType struct {
	Error error
}

type contentJSON struct {
	Content string `json:"content"`
}

type errorJSON struct {
	Error string `json:"error"`
}

// renderErrorPage renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, r *http.Request, err error) {
	user, ok := getUserFromSession(r)
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
