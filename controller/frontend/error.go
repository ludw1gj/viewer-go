// This file contains error page functions.

package frontend

import (
	"bytes"
	"log"
	"net/http"

	"fmt"

	"github.com/robertjeffs/viewer-go/logic/common"

	"github.com/robertjeffs/viewer-go/model/database"
)

// NotFound renders the not found page and sends status 404.
func NotFound(w http.ResponseWriter, r *http.Request) {
	user, err := common.ValidateUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var templateBuf bytes.Buffer
	if err := notFoundTemplate.ExecuteTemplate(&templateBuf, "base.gohtml", userInfo{user}); err != nil {
		renderErrorPage(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(templateBuf.Bytes())
}

// renderErrorPage renders the error page and sends status 500.
func renderErrorPage(w http.ResponseWriter, r *http.Request, pageErr error) {
	w.WriteHeader(http.StatusInternalServerError)
	user, err := common.ValidateUser(r)
	if err != nil {
		log.Printf("StatusInternalServerError failed to execute get user from session on error page: %s", err.Error())

		resp := fmt.Sprintf("500: Server error. Two errors have occured.<br>First Error: %s<br>Second Error: %s",
			pageErr.Error(), err.Error())
		w.Write([]byte(resp))
		return
	}

	data := struct {
		Error string
		User  database.User
	}{
		pageErr.Error(),
		user,
	}

	var templateBuf bytes.Buffer
	if err := errorTemplate.ExecuteTemplate(&templateBuf, "base.gohtml", data); err != nil {
		log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())

		resp := fmt.Sprintf("500: Server error. Two errors have occured.<br>First Error: %s<br>Second Error: %s",
			pageErr.Error(), err.Error())
		w.Write([]byte(resp))
		return
	}
	w.Write(templateBuf.Bytes())
}
