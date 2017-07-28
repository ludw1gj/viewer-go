package controller

import (
	"net/http"

	"github.com/FriedPigeon/viewer-go/db"
	"github.com/FriedPigeon/viewer-go/session"
)

func ValidateUser(r *http.Request) (user db.User, err error) {
	userId, err := session.GetUserIDFromSession(r)
	if err != nil {
		return
	}
	user, err = db.GetUser(userId)
	if err != nil {
		return
	}
	return
}
