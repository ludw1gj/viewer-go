// Package route contains the Load function for initialising a router instance, which has registered routes and a static
// file handler for development purposes.
package route

import (
	"flag"
	"net/http"

	"github.com/FriedPigeon/viewer-go/handler"
	"github.com/gorilla/mux"
)

// Load initialises routes and a static file handler if dev flag is used.
func Load() {
	protected := mux.NewRouter()

	// protected routes
	protected.HandleFunc("/", handler.RedirectToViewer).Methods("GET")
	protected.HandleFunc("/viewer/{path:.*}", handler.Viewer).Methods("GET")
	protected.HandleFunc("/about", handler.About).Methods("GET")
	protected.HandleFunc("/user", handler.User).Methods("GET")

	protected.HandleFunc("/logout", handler.Logout).Methods("GET")
	protected.HandleFunc("/upload", handler.Upload).Methods("POST")
	protected.HandleFunc("/create-folder", handler.CreateFolder).Methods("POST")
	protected.HandleFunc("/delete", handler.Delete).Methods("POST")
	protected.HandleFunc("/delete-all", handler.DeleteAll).Methods("POST")
	protected.NotFoundHandler = http.HandlerFunc(handler.NotFound)

	http.Handle("/", authenticateRoute(protected))

	// open routes
	http.HandleFunc("/login", handler.Login)

	// static file handler in dev mode
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
		isAuth := handler.CheckIfAuth(w, r)

		// if user is authenticated, proceed to route
		if isAuth {
			h.ServeHTTP(w, r)
		}
	})
}
