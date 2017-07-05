package controller

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/FriedPigeon/viewer-go/db"
)

type userController struct{}

func NewUserController() *userController {
	return &userController{}
}

func (userController) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := loginTpl.Execute(w, errType{nil})
		if err != nil {
			log.Println(err)
		}
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := newUserSession(w, r, username, password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			err = loginTpl.Execute(w, errType{err})
			if err != nil {
				log.Println(err)
			}
			return
		}
		http.Redirect(w, r, viewerRootURL, http.StatusSeeOther)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
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

func (userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	passwords := struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&passwords)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}

	err = db.ChangeUserPassword(user, passwords.OldPassword, passwords.NewPassword)
	if err != nil {
		if err.Error() == "Incorrect password." {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorJSON{err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Password changed successfully."})
}

func (userController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pw := r.FormValue("password")
	err = db.DeleteUserPasswordValidated(user, pw)
	if err != nil {
		// TODO: handle error properly
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
