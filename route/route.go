// Package route contains the Load function for initialising a router instance, which has registered routes and a static
// file handler for development purposes.
package route

import (
	"flag"
	"net/http"

	"github.com/FriedPigeon/viewer-go/controller"
	"github.com/gorilla/mux"
)

// Load initialises routes and a static file controller if dev flag is used.
func Load() {
	protected := mux.NewRouter()

	sc := controller.NewSiteController()
	vc := controller.NewViewerController()
	uc := controller.NewUserController()

	// -- protected routes --
	// site
	protected.HandleFunc("/", sc.RedirectToViewer).Methods("GET")
	protected.HandleFunc("/about", sc.About).Methods("GET")
	protected.NotFoundHandler = http.HandlerFunc(sc.NotFound)

	// viewer
	protected.HandleFunc("/viewer/{path:.*}", vc.Viewer).Methods("GET")

	protected.HandleFunc("/upload", vc.Upload).Methods("POST")
	protected.HandleFunc("/create-folder", vc.CreateFolder).Methods("POST")
	protected.HandleFunc("/delete", vc.Delete).Methods("POST")
	protected.HandleFunc("/delete-all", vc.DeleteAll).Methods("POST")

	// user
	protected.HandleFunc("/user", uc.UserPage).Methods("GET")
	protected.HandleFunc("/user/logout", uc.Logout).Methods("GET")
	protected.HandleFunc("/user/password", uc.ChangePassword).Methods("POST")
	protected.HandleFunc("/user/delete", uc.DeleteAccount).Methods("POST")
	//-- end --

	// -- open routes --
	http.HandleFunc("/login", uc.Login)
	http.Handle("/", authenticateRoute(protected))
	//-- end --

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
		isAuth := controller.CheckIfAuth(w, r)

		// if user is authenticated, proceed to route
		if isAuth {
			h.ServeHTTP(w, r)
		}
	})
}
