package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robertjeffs/viewer-go/controllers"
	"github.com/robertjeffs/viewer-go/logic/session"
)

// LoadRoutes initialises routes and a assets file handler if dev is true.
func LoadRoutes() {
	protected := mux.NewRouter()

	// controllers
	sc := controllers.NewSiteController()
	vc := controllers.NewViewerController()
	uc := controllers.NewUserController()
	ac := controllers.NewAdminController()

	http.HandleFunc("/login", sc.GetLoginPage)
	http.HandleFunc("/api/user/login", uc.Login)
	http.Handle("/", authenticateRoute(protected))

	// site
	protected.HandleFunc("/", redirectToViewerPage).Methods("GET")
	protected.HandleFunc("/viewer/{path:.*}", sc.GetViewerPage).Methods("GET")
	protected.HandleFunc("/file/{path:.*}", sc.SendFile).Methods("GET")
	protected.HandleFunc("/about", sc.GetAboutPage).Methods("GET")
	protected.HandleFunc("/user", sc.GetUserPage).Methods("GET")
	protected.HandleFunc("/admin", sc.GetAdminPage).Methods("GET")
	protected.HandleFunc("/admin/users", sc.GetAdminDisplayAllUsers).Methods("GET")
	protected.NotFoundHandler = http.HandlerFunc(sc.GetNotFoundPage)

	// api
	protected.HandleFunc("/api/viewer/upload", vc.Upload).Methods("POST")
	protected.HandleFunc("/api/viewer/create", vc.CreateFolder).Methods("POST")
	protected.HandleFunc("/api/viewer/delete", vc.Delete).Methods("POST")
	protected.HandleFunc("/api/viewer/delete-all", vc.DeleteAll).Methods("POST")

	protected.HandleFunc("/api/user/logout", uc.Logout).Methods("POST")
	protected.HandleFunc("/api/user/delete", uc.DeleteAccount).Methods("POST")
	protected.HandleFunc("/api/user/change-password", uc.ChangePassword).Methods("POST")
	protected.HandleFunc("/api/user/change-name", uc.ChangeName).Methods("POST")

	protected.HandleFunc("/api/admin/change-username", ac.ChangeUserUsername).Methods("POST")
	protected.HandleFunc("/api/admin/create-user", ac.CreateUser).Methods("POST")
	protected.HandleFunc("/api/admin/delete-user", ac.DeleteUser).Methods("POST")
	protected.HandleFunc("/api/admin/change-dir-root", ac.ChangeDirRoot).Methods("POST")
	protected.HandleFunc("/api/admin/change-admin-status", ac.ChangeUserAdminStatus).Methods("POST")

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
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

func redirectToViewerPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}
