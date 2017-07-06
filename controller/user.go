// This files contains the user controller, which contains methods for logging in and logging out for users, user
// actions, and rendering the user page.

package controller

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

type userController struct{}

func NewUserController() *userController {
	return &userController{}
}

// Login method when accessed via a GET request renders the login page, and when when accessed via a POST request it
// will process the login form values and login the user and redirect the user the viewer page.
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

		err := session.NewUserSession(w, r, username, password)
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

// Logout will logout the user by changing the session value "authenticated" to false.
func (userController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := session.RemoveUserAuthFromSession(w, r)
	if err != nil {
		// TODO: controller error properly
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// UserPage renders the user page.
func (userController) UserPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
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

// ChangePassword will process a json post request, comparing password sent with current password and if they match, the
// current password will be changed to the new password.
func (userController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	passwords := struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
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

// DeleteAccount will process the delete user form, if password is correct the user's account will be deleted and the
// user redirected to the login page.
func (userController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
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
