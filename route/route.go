// Package route contains the Load function for initialising a router instance, which has registered routes and a static
// file handler for development purposes.
package route

import (
	"flag"
	"net/http"

	"github.com/FriedPigeon/viewer-go/controller/api"
	"github.com/FriedPigeon/viewer-go/controller/frontend"
	"github.com/FriedPigeon/viewer-go/session"
	"github.com/gorilla/mux"
)

// Load initialises routes and a static file controller if dev flag is used.
func Load() {
	protected := mux.NewRouter()

	http.HandleFunc("/login", frontend.LoginPage)
	http.HandleFunc("/api/user/login", api.Login)
	http.Handle("/", authenticateRoute(protected))

	// frontend
	protected.HandleFunc("/", frontend.RedirectToViewer).Methods("GET")
	protected.HandleFunc("/viewer/{path:.*}", frontend.ViewerPage).Methods("GET")
	protected.HandleFunc("/file/{path:.*}", frontend.SendFile).Methods("GET")
	protected.HandleFunc("/about", frontend.AboutPage).Methods("GET")
	protected.HandleFunc("/user", frontend.UserPage).Methods("GET")
	protected.HandleFunc("/admin", frontend.AdminPage).Methods("GET")
	protected.HandleFunc("/admin/users", frontend.AdminDisplayAllUsers).Methods("GET")
	protected.NotFoundHandler = http.HandlerFunc(frontend.NotFound)

	// api
	protected.HandleFunc("/api/viewer/upload", api.Upload).Methods("POST")
	protected.HandleFunc("/api/viewer/create", api.CreateFolder).Methods("POST")
	protected.HandleFunc("/api/viewer/delete", api.Delete).Methods("POST")
	protected.HandleFunc("/api/viewer/delete-all", api.DeleteAll).Methods("POST")

	protected.HandleFunc("/api/user/logout", api.Logout).Methods("POST")
	protected.HandleFunc("/api/user/delete", api.DeleteAccount).Methods("POST")
	protected.HandleFunc("/api/user/change-password", api.ChangePassword).Methods("POST")
	protected.HandleFunc("/api/user/change-name", api.ChangeName).Methods("POST")

	protected.HandleFunc("/api/admin/change-username", api.ChangeUsername).Methods("POST")
	protected.HandleFunc("/api/admin/create-user", api.CreateUser).Methods("POST")
	protected.HandleFunc("/api/admin/delete-user", api.DeleteUser).Methods("POST")
	protected.HandleFunc("/api/admin/change-dir-root", api.ChangeDirRoot).Methods("POST")

	// static file controller in dev mode
	dev := flag.Bool("dev", false, "Use in development")
	flag.Parse()
	if *dev {
		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static", fs))
	}
}

// authenticateRoute is middleware that checks if users are authenticated.
func authenticateRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuth := session.CheckIfAuth(r)
		// if user is authenticated, proceed to route
		if !isAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}
