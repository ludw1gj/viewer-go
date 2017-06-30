package controller

import (
	"bytes"
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
)

// NotFound renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, r *http.Request, err error) {
	user, ok := getUserFromSession(r)
	if ok != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
