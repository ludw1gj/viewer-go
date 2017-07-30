// This file contains template logic.

package frontend

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path"

	"fmt"
)

var (
	tplDir      = path.Join("templates", "frontend")
	baseTplPath = path.Join(tplDir, "base.gohtml")

	// templates
	loginTpl   = template.Must(template.ParseFiles(path.Join(tplDir, "login.gohtml")))
	dirListTpl = template.Must(template.ParseFiles(path.Join(tplDir, "dir_list.gohtml")))

	// templates that require base template
	viewerTpl     = initTemplate("viewer", tplDir, baseTplPath)
	aboutTpl      = initTemplate("about", tplDir, baseTplPath)
	userTpl       = initTemplate("user", tplDir, baseTplPath)
	adminTpl      = initTemplate("admin", tplDir, baseTplPath)
	adminUsersTpl = initTemplate("admin_users", tplDir, baseTplPath)
	errorTpl      = initTemplate("error", tplDir, baseTplPath)
	notFoundTpl   = initTemplate("not_found", tplDir, baseTplPath)

	// function map for use in templates.
	funcMap = template.FuncMap{
		"genDirectoryList": func(userDirRoot string, urlPath string) template.HTML {
			list, err := genDirectoryList(userDirRoot, urlPath)
			if err != nil {
				errMsg := fmt.Sprintf("There has been an error getting directory list: %s", err.Error())
				return template.HTML(errMsg)
			}
			return list
		},
	}
)

// initTemplate creates new template.Template and parses files of tplName in the given template directory (tplDir).
func initTemplate(tplName string, tplDir string, baseTplPath string) *template.Template {
	return template.Must(template.New(tplName).Funcs(funcMap).
		ParseFiles(baseTplPath, path.Join(tplDir, tplName+".gohtml")))
}

// renderTemplate executes a templates and sends it to the client.
func renderTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	var tplBuf bytes.Buffer
	err := tpl.ExecuteTemplate(&tplBuf, "base.gohtml", data)
	if err != nil {
		log.Println(err)

		renderErrorPage(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(tplBuf.Bytes())
}
