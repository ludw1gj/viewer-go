package api

import (
	"net/http"

	"encoding/json"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

// CreateUser receives new user information via json and creates the user. Client must be admin.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}
	if !user.Admin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	u := db.User{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	err = db.CreateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully created user."})
}

// DeleteUser receives user information via json and deletes the user. Client must be admin.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := session.GetUserFromSession(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}
	if !user.Admin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorJSON{"Unauthorised."})
		return
	}

	data := struct {
		ID int `json:"id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	// check if use exists
	_, err = db.GetUser(data.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	err = db.DeleteUser(data.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJSON{err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contentJSON{"Successfully deleted user."})
}
