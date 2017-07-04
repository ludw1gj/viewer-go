package controller

import (
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
)

type userController struct{}

func NewUserController() *userController {
	return &userController{}
}

func (userController) UserPage(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = userTpl.Execute(w, userInfo{user})
	if err != nil {
		// TODO: handler error properly
		log.Println(err)
	}
}

func (userController) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderLogin(w, nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := newUserSession(w, r, username, password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			renderLogin(w, err)
			return
		}
		http.Redirect(w, r, viewerRootURL, http.StatusSeeOther)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func renderLogin(w http.ResponseWriter, err error) {
	type Err struct {
		Error error `json:"error"`
	}

	err = loginTpl.Execute(w, Err{err})
	if err != nil {
		// TODO: handler error properly
		log.Println(err.Error())
	}
}

func (userController) Logout(w http.ResponseWriter, r *http.Request) {
	err := removeUserAuthFromSession(w, r)
	if err != nil {
		// TODO: controller error properly
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	oldPw := r.FormValue("old-password")
	newPw := r.FormValue("new-password")
	err = db.ChangeUserPassword(user, oldPw, newPw)
	if err != nil {
		// TODO: handle error properly
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func (userController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pw := r.FormValue("password")
	err = db.DeleteUser(user, pw)
	if err != nil {
		// TODO: handle error properly
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
