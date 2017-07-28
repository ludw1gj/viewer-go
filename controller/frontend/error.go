// This file contains error page functions.

package frontend

import (
	"bytes"
	"log"
	"net/http"

	"fmt"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/controller"
)

// NotFound renders the not found page and sends status 404.
func NotFound(w http.ResponseWriter, r *http.Request) {
	user, err := controller.ValidateUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var tplBuf bytes.Buffer
	err = notFoundTpl.ExecuteTemplate(&tplBuf, "base.gohtml", userInfo{user})
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(tplBuf.Bytes())
}

// renderErrorPage renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, r *http.Request, pageErr error) {
	w.WriteHeader(http.StatusInternalServerError)
	user, err := controller.ValidateUser(r)
	if err != nil {
		log.Printf("StatusInternalServerError failed to execute get user from session on error page: %s", err.Error())

		resp := fmt.Sprintf("500: Server error. Two errors have occured.<br>First Error: %s<br>Second Error: %s",
			pageErr.Error(), err.Error())
		w.Write([]byte(resp))
		return
	}

	data := struct {
		Error string
		User  db.User
	}{
		pageErr.Error(),
		user,
	}

	var tplBuf bytes.Buffer
	err = errorTpl.ExecuteTemplate(&tplBuf, "base.gohtml", data)
	if err != nil {
		log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())

		resp := fmt.Sprintf("500: Server error. Two errors have occured.<br>First Error: %s<br>Second Error: %s",
			pageErr.Error(), err.Error())
		w.Write([]byte(resp))
		return
	}
	w.Write(tplBuf.Bytes())
}
