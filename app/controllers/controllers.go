// Package controllers contains controllers of route responses.
package controllers

import (
	"database/sql"

	"github.com/ludw1gj/viewer-go/app/logic/session"
)

// NewSiteController returns a SiteController instance.
func NewSiteController(db *sql.DB, sm *session.Manager) *SiteController {
	return &SiteController{
		db,
		sm,
	}
}

// NewUserAPIController returns a UserAPIController instance.
func NewUserAPIController(db *sql.DB, sm *session.Manager) *UserAPIController {
	return &UserAPIController{
		db,
		sm,
	}
}

// NewAdminAPIController returns a AdminAPIController instance.
func NewAdminAPIController(db *sql.DB, sm *session.Manager) *AdminAPIController {
	return &AdminAPIController{
		db,
		sm,
	}
}

// NewViewerAPIController returns a ViewerAPIController instance.
func NewViewerAPIController(sm *session.Manager) *ViewerAPIController {
	return &ViewerAPIController{
		sm,
	}
}
