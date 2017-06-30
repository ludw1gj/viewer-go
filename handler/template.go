// This file contains template logic.

package handler

import (
	"html/template"
	"path"
)

const tplDir = "templates"

var (
	baseTplPath = path.Join(tplDir, "base.gohtml")

	// independent templates
	loginTpl,
	dirListTpl *template.Template

	// site templates
	viewerTpl,
	aboutTpl,
	userTpl,
	errorTpl,
	notFoundTpl *template.Template
)

func init() {
	loginTpl = initTemplate("login.gohtml", false)
	dirListTpl = initTemplate("directory_list.gohtml", false)

	viewerTpl = initTemplate("viewer.gohtml", true)
	aboutTpl = initTemplate("about.gohtml", true)
	userTpl = initTemplate("user.gohtml", true)
	errorTpl = initTemplate("error.gohtml", true)
	notFoundTpl = initTemplate("not_found.gohtml", true)
}

// initTemplate creates new template.Template and parses files of tplName in the template directory (tplDir).
func initTemplate(tplName string, withBase bool) *template.Template {
	if withBase {
		return template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, tplName)))
	}
	return template.Must(template.ParseFiles(path.Join(tplDir, tplName)))
}
