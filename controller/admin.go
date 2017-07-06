// This file contains the admin controller, which has methods concerning administrations pages and operations.

package controller

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

type adminController struct{}

func NewAdminController() *adminController {
	return &adminController{}
}

// AdminPage renders the Administration page. Client must be admin.
func (adminController) AdminPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.IsAdmin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}

	err = adminTpl.Execute(w, userInfo{user})
	if err != nil {
		// TODO: handle error properly
		log.Println(err)
	}
}

// DisplayAllUsers render a sub administration page which displays all users in database. Client must be admin.
func (adminController) DisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.IsAdmin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}

	users, err := db.GetAllUsers()
	if err != nil {
		renderErrorPage(w, r, err)
	}
	err = adminUsersTpl.Execute(w, struct {
		User  db.User
		Users []db.User
	}{user, users})
	if err != nil {
		log.Println(err)
	}
}

// CreateUser receives new user information via json and creates the user. Client must be admin.
func (adminController) CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.IsAdmin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	u := db.User{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	if err != nil {
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}

	err = db.CreateUser(u)
	if err != nil {
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}
	json.NewEncoder(w).Encode(contentJSON{"Successfully created user."})
}

// DeleteUser receives user information via json and deletes the user. Client must be admin.
func (adminController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.IsAdmin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		ID int `json:"id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}

	err = db.DeleteUser(data.ID)
	if err != nil {
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}
	json.NewEncoder(w).Encode(contentJSON{"Successfully deleted user."})
}
