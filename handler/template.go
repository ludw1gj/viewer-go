package handler

import (
	"html/template"
	"path"
)

var (
	baseTplPath = path.Join("templates", "base.gohtml")

	viewerTpl   *template.Template
	dirListTpl  *template.Template
	errorTpl    *template.Template
)

func init() {
	viewerTpl = template.Must(template.ParseFiles(baseTplPath, path.Join("templates", "viewer.gohtml")))
	dirListTpl = template.Must(template.ParseFiles(path.Join("templates", "view.gohtml")))
	errorTpl = template.Must(template.ParseFiles(baseTplPath, path.Join("templates", "error.gohtml")))
}
