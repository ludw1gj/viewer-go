package frontend

import (
	"github.com/FriedPigeon/viewer-go/db"
)

// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User db.User
}
