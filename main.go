package main

import (
	"log"
	"net/http"

	"flag"

	"fmt"

	"github.com/FriedPigeon/viewer-go/controller/api"
	"github.com/FriedPigeon/viewer-go/controller/frontend"
	"github.com/FriedPigeon/viewer-go/controller/middleware"
	"github.com/FriedPigeon/viewer-go/database"
	"github.com/FriedPigeon/viewer-go/session"
	"github.com/gorilla/mux"
)

func main() {
	dev := flag.Bool("dev", true, "Use in development")
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
	loadRoutes(*dev)

	// listen and serve
	log.Printf("viewer-go listening on port %d...", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatalln(err)
	}
}

// loadRoutes initialises routes and a static file handler if dev is true.
func loadRoutes(dev bool) {
	protected := mux.NewRouter()

	http.HandleFunc("/login", frontend.LoginPage)
	http.HandleFunc("/api/user/login", api.Login)
	http.Handle("/", middleware.AuthenticateRoute(protected))

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

	// static file handler in dev mode
	if dev {
		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static", fs))
	}
}
