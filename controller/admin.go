package controller

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/FriedPigeon/viewer-go/db"
)

type adminController struct{}

func NewAdminController() *adminController {
	return &adminController{}
}

func (adminController) AdminPage(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
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

func (adminController) DisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
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

func (adminController) CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
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

func (adminController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.IsAdmin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	id := struct {
		ID int `json:"id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&id)
	if err != nil {
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}

	err = db.DeleteUser(id.ID)
	if err != nil {
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
	}
	json.NewEncoder(w).Encode(contentJSON{"Successfully deleted user."})
}
