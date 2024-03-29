// Package router initialises a mux.Router instance and registers routes including a file server.
package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ludw1gj/viewer-go/app/controllers"
	"github.com/ludw1gj/viewer-go/app/logic/session"
)

// LoadRoutes initialises routes and a file server.
func LoadRoutes(dbConn *sql.DB, sm *session.Manager) {
	protected := mux.NewRouter()

	// controllers
	siteController := controllers.NewSiteController(dbConn, sm)
	viewerAPIController := controllers.NewViewerAPIController(sm)
	userAPIController := controllers.NewUserAPIController(dbConn, sm)
	adminAPIController := controllers.NewAdminAPIController(dbConn, sm)

	http.HandleFunc("/login", siteController.GetLoginPage)
	http.HandleFunc("/api/user/login", userAPIController.Login)
	http.Handle("/", authenticateRoute(protected, sm))

	// site
	protected.HandleFunc("/", redirectToViewerPage).Methods("GET")
	protected.HandleFunc("/viewer/{path:.*}", siteController.GetViewerPage).Methods("GET")
	protected.HandleFunc("/file/{path:.*}", siteController.SendFile).Methods("GET")
	protected.HandleFunc("/about", siteController.GetAboutPage).Methods("GET")
	protected.HandleFunc("/user", siteController.GetUserPage).Methods("GET")
	protected.HandleFunc("/admin", siteController.GetAdminPage).Methods("GET")
	protected.HandleFunc("/admin/users", siteController.GetAdminDisplayAllUsers).Methods("GET")
	protected.NotFoundHandler = http.HandlerFunc(siteController.GetNotFoundPage)

	// api
	protected.HandleFunc("/api/viewer/upload/{path:.*}", viewerAPIController.Upload).Methods("POST")
	protected.HandleFunc("/api/viewer/create", viewerAPIController.CreateFolder).Methods("POST")
	protected.HandleFunc("/api/viewer/delete", viewerAPIController.Delete).Methods("POST")
	protected.HandleFunc("/api/viewer/delete-all", viewerAPIController.DeleteAll).Methods("POST")

	protected.HandleFunc("/api/user/logout", userAPIController.Logout).Methods("POST")
	protected.HandleFunc("/api/user/delete", userAPIController.DeleteAccount).Methods("POST")
	protected.HandleFunc("/api/user/change-password", userAPIController.ChangePassword).Methods("POST")
	protected.HandleFunc("/api/user/change-name", userAPIController.ChangeName).Methods("POST")

	protected.HandleFunc("/api/admin/change-username", adminAPIController.ChangeUserUsername).Methods("POST")
	protected.HandleFunc("/api/admin/create-user", adminAPIController.CreateUser).Methods("POST")
	protected.HandleFunc("/api/admin/delete-user", adminAPIController.DeleteUser).Methods("POST")
	protected.HandleFunc("/api/admin/change-dir-root", adminAPIController.ChangeDirRoot).Methods("POST")
	protected.HandleFunc("/api/admin/change-admin-status", adminAPIController.ChangeUserAdminStatus).Methods("POST")

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
}

func authenticateRoute(h http.Handler, sm *session.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user is authenticated
		if isAuth := sm.CheckUserAuth(r); !isAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// if user is authenticated, proceed to route
		h.ServeHTTP(w, r)
	})
}

func redirectToViewerPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}
