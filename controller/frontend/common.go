package frontend

import "github.com/FriedPigeon/viewer-go/db"

// TODO: may not be needed
// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User db.User
}
