package session

import (
	"log"
	"net/http"

	"github.com/FriedPigeon/viewer-go/config"
	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/model"
	"github.com/FriedPigeon/viewer-go/tpl"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(config.CookieStoreAuthKey))

func GetUserFromSession(r *http.Request) (user model.User, err error) {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		log.Println(err)
	}

	id := session.Values["id"].(int)
	user, err = db.GetUser(id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CheckIfAuth(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		log.Println(err)
	}

	// check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// user is not auth, send to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	return true
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusForbidden)
		tpl.LoginTpl.Execute(w, nil)
	case "POST":
		session, err := store.Get(r, "viewer-session")
		if err != nil {
			log.Println(err)
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// validate authentication
		user, validated := db.ValidateUser(username, password)
		if validated != true {
			http.Redirect(w, r, config.ViewerRootURL, http.StatusSeeOther)
			return
		}

		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Values["id"] = user.ID
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, config.ViewerRootURL, http.StatusSeeOther)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "viewer-session")
	if err != nil {
		log.Println(err)
	}

	// revoke user's authentication
	session.Values["authenticated"] = false
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
