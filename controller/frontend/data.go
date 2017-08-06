package frontend

import "github.com/FriedPigeon/viewer-go/database"

// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User database.User
}
