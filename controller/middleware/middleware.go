// Package middleware contains middleware for routes.
package middleware

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/session"
)

// AuthenticateRoute is middleware that checks if users are authenticated.
func AuthenticateRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuth := session.CheckUserAuth(r)
		// if user is authenticated, proceed to route
		if !isAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}
