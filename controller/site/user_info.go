package site

import "github.com/robertjeffs/viewer-go/model/database"

// userInfo is used for data object of error for rendering templates.
type userInfo struct {
	User database.User
}
