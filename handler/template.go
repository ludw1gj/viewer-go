package handler

import (
	"html/template"
	"path"
)

var (
	viewerTpl  *template.Template
	dirListTpl *template.Template
)

func init() {
	viewerTpl = template.Must(template.ParseFiles(path.Join("templates", "index.gohtml")))
	dirListTpl = template.Must(template.ParseFiles(path.Join("templates", "view.gohtml")))
}
