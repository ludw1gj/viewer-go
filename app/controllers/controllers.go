// Package controllers contains controllers of route responses.
package controllers

import "github.com/robertjeffs/viewer-go/app/logic/session"

// NewSiteController returns a SiteController instance.
func NewSiteController(sm *session.Manager) *SiteController {
	return &SiteController{
		sm,
	}
}

// NewUserAPIController returns a UserAPIController instance.
func NewUserAPIController(sm *session.Manager) *UserAPIController {
	return &UserAPIController{
		sm,
	}
}

// NewAdminAPIController returns a AdminAPIController instance.
func NewAdminAPIController(sm *session.Manager) *AdminAPIController {
	return &AdminAPIController{
		sm,
	}
}

// NewViewerAPIController returns a ViewerAPIController instance.
func NewViewerAPIController(sm *session.Manager) *ViewerAPIController {
	return &ViewerAPIController{
		sm,
	}
}
