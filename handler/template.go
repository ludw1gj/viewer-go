// This files contains the template directory, and templates.

package handler

import (
	"html/template"
	"path"
)

const tplDir = "templates"

var (
	baseTplPath = path.Join(tplDir, "base.gohtml")

	viewerTpl,
	dirListTpl,
	aboutTpl,
	errorTpl,
	notFoundTpl *template.Template
)

func init() {
	viewerTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "viewer.gohtml")))
	dirListTpl = template.Must(template.ParseFiles(path.Join(tplDir, "directory_list.gohtml")))
	aboutTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "about.gohtml")))
	errorTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "error.gohtml")))
	notFoundTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "not_found.gohtml")))
}
