package frontend

import (
	"net/http"

	"github.com/robertjeffs/viewer-go/logic/common"
	"github.com/robertjeffs/viewer-go/model/database"
)

// AdminPage renders the Administration page. Client must be admin.
func AdminPage(w http.ResponseWriter, r *http.Request) {
	u, err := common.ValidateAdmin(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, r, adminTemplate, userInfo{u})
}

// AdminDisplayAllUsers render a sub administration page which displays all users in database. Client must be admin.
func AdminDisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := common.ValidateAdmin(r)
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
	}{u, users}
	renderTemplate(w, r, adminUsersTemplate, data)
}
