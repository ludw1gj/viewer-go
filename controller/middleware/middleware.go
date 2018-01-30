// Package middleware contains middleware for routes.
package middleware

import (
	"net/http"

	"github.com/robertjeffs/viewer-go/session"
)

// AuthenticateRoute is middleware that checks if users are authenticated.
func AuthenticateRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAuth := session.CheckUserAuth(r); !isAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// if user is authenticated, proceed to route
		h.ServeHTTP(w, r)
	})
}
