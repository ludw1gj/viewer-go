// Package tpl contains template logic.
package tpl

import (
	"html/template"
	"path"
)

const tplDir = "templates"

var (
	baseTplPath = path.Join(tplDir, "base.gohtml")

	LoginTpl,
	ViewerTpl,
	DirListTpl,
	AboutTpl,
	ErrorTpl,
	NotFoundTpl *template.Template
)

func init() {
	LoginTpl = template.Must(template.ParseFiles(path.Join(tplDir, "login.gohtml")))
	ViewerTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "viewer.gohtml")))
	DirListTpl = template.Must(template.ParseFiles(path.Join(tplDir, "directory_list.gohtml")))
	AboutTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "about.gohtml")))
	ErrorTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "error.gohtml")))
	NotFoundTpl = template.Must(template.ParseFiles(baseTplPath, path.Join(tplDir, "not_found.gohtml")))
}
