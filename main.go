package main

import (
	"log"
	"net/http"

	"flag"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/robertjeffs/viewer-go/controller/api"
	"github.com/robertjeffs/viewer-go/controller/site"
	"github.com/robertjeffs/viewer-go/logic/session"
	"github.com/robertjeffs/viewer-go/model/database"
)

func main() {
	port := flag.Int("port", 3000, "Port number")
	dbFile := flag.String("dbFile", "viewer.db", "Database File")
	sessionConfigFile := flag.String("configFile", "config.json", "Session config json file")
	flag.Parse()

	// load database, session, and routes
	err := database.Load(*dbFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = session.Load(*sessionConfigFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	loadRoutes()

	// listen and serve
	log.Printf("viewer-go listening on port %d...", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}

// loadRoutes initialises routes and a static file handler if dev is true.
func loadRoutes() {
	protected := mux.NewRouter()

	http.HandleFunc("/login", site.GetLoginPage)
	http.HandleFunc("/api/user/login", api.UserLogin)
	http.Handle("/", authenticateRoute(protected))

	// redirectToViewer redirects users to the viewer page.
	redirectToViewer := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
	}

	// site
	protected.HandleFunc("/", redirectToViewer).Methods("GET")
	protected.HandleFunc("/viewer/{path:.*}", site.GetViewerPage).Methods("GET")
	protected.HandleFunc("/file/{path:.*}", site.SendFile).Methods("GET")
	protected.HandleFunc("/about", site.GetAboutPage).Methods("GET")
	protected.HandleFunc("/user", site.GetUserPage).Methods("GET")
	protected.HandleFunc("/admin", site.GetAdminPage).Methods("GET")
	protected.HandleFunc("/admin/users", site.GetAdminDisplayAllUsers).Methods("GET")
	protected.NotFoundHandler = http.HandlerFunc(site.GetNotFoundPage)

	// api
	protected.HandleFunc("/api/viewer/upload", api.ViewerUpload).Methods("POST")
	protected.HandleFunc("/api/viewer/create", api.ViewerCreateFolder).Methods("POST")
	protected.HandleFunc("/api/viewer/delete", api.ViewerDelete).Methods("POST")
	protected.HandleFunc("/api/viewer/delete-all", api.ViewerDeleteAll).Methods("POST")

	protected.HandleFunc("/api/user/logout", api.UserLogout).Methods("POST")
	protected.HandleFunc("/api/user/delete", api.UserDeleteAccount).Methods("POST")
	protected.HandleFunc("/api/user/change-password", api.UserChangePassword).Methods("POST")
	protected.HandleFunc("/api/user/change-name", api.UserChangeName).Methods("POST")

	protected.HandleFunc("/api/admin/change-username", api.AdminChangeUserUsername).Methods("POST")
	protected.HandleFunc("/api/admin/create-user", api.AdminCreateUser).Methods("POST")
	protected.HandleFunc("/api/admin/delete-user", api.AdminDeleteUser).Methods("POST")
	protected.HandleFunc("/api/admin/change-dir-root", api.AdminChangeDirRoot).Methods("POST")
	protected.HandleFunc("/api/admin/change-admin-status", api.AdminChangeUserAdminStatus).Methods("POST")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
}

// AuthenticateRoute is middleware that checks if users are authenticated.
func authenticateRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAuth := session.CheckUserAuth(r); !isAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// if user is authenticated, proceed to route
		h.ServeHTTP(w, r)
	})
}
