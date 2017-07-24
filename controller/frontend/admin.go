package frontend

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

// AdminPage renders the Administration page. Client must be admin.
func AdminPage(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.Admin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}
	renderTemplate(w, r, adminTpl, userInfo{user})
}

// AdminDisplayAllUsers render a sub administration page which displays all users in database. Client must be admin.
func AdminDisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !user.Admin {
		http.Redirect(w, r, "/viewer", http.StatusForbidden)
		return
	}

	users, err := db.GetAllUsers()
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}

	data := struct {
		User  db.User
		Users []db.User
	}{user, users}
	renderTemplate(w, r, adminUsersTpl, data)
}
