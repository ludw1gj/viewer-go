package site

import (
	"net/http"

	"github.com/robertjeffs/viewer-go/controller/templates"
	"github.com/robertjeffs/viewer-go/logic/common"
	"github.com/robertjeffs/viewer-go/model/database"
)

// GetAdminPage renders the Administration page. Client must be admin.
func GetAdminPage(w http.ResponseWriter, r *http.Request) {
	u, err := common.ValidateAdmin(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.RenderTemplate(w, r, "admin", userInfo{u})
}

// GetAdminDisplayAllUsers render a sub administration page which displays all users in database. Client must be admin.
func GetAdminDisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := common.ValidateAdmin(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		GetErrorPage(w, r, err)
		return
	}

	data := struct {
		User  database.User
		Users []database.User
	}{u, users}
	templates.RenderTemplate(w, r, "adminUsers", data)
}
