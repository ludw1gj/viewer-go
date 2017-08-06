package frontend

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/common"

	"github.com/FriedPigeon/viewer-go/database"
)

// validateAdmin checks if the user is valid and is admin.
func validateAdmin(r *http.Request) (user database.User, err error) {
	user, err = common.ValidateUser(r)
	if err != nil {
		return
	}
	if !user.Admin {
		return
	}
	return
}

// AdminPage renders the Administration page. Client must be admin.
func AdminPage(w http.ResponseWriter, r *http.Request) {
	user, err := validateAdmin(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, r, adminTpl, userInfo{user})
}

// AdminDisplayAllUsers render a sub administration page which displays all users in database. Client must be admin.
func AdminDisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := validateAdmin(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}

	data := struct {
		User  database.User
		Users []database.User
	}{user, users}
	renderTemplate(w, r, adminUsersTpl, data)
}
