// This file contains template logic.

package frontend

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path"

	"fmt"

	"github.com/FriedPigeon/viewer-go/controller/api"
)

var (
	tplDir      = path.Join("templates", "frontend")
	baseTplPath = path.Join(tplDir, "base.gohtml")

	// frontend templates
	loginTpl,
	viewerTpl,
	aboutTpl,
	userTpl,
	adminTpl,
	adminUsersTpl,
	errorTpl,
	notFoundTpl *template.Template
)

func init() {
	loginTpl = template.Must(template.ParseFiles(path.Join(tplDir, "login.gohtml")))

	viewerTpl = initTemplate("viewer", true)
	aboutTpl = initTemplate("about", true)
	userTpl = initTemplate("user", true)
	adminTpl = initTemplate("admin", true)
	adminUsersTpl = initTemplate("admin_users", true)
	errorTpl = initTemplate("error", true)
	notFoundTpl = initTemplate("not_found", true)
}

// initTemplate creates new template.Template and parses files of tplName in the template directory (tplDir).
func initTemplate(tplName string, withBase bool) *template.Template {
	if withBase {
		return template.Must(
			template.New(tplName).Funcs(funcMap).ParseFiles(baseTplPath, path.Join(tplDir, tplName+".gohtml")))
	}
	return template.Must(template.New(tplName).Funcs(funcMap).ParseFiles(path.Join(tplDir, tplName+".gohtml")))
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tpl *template.Template, data interface{}) {
	var buf bytes.Buffer
	err := tpl.ExecuteTemplate(&buf, "base.gohtml", data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		renderErrorPage(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

var funcMap = template.FuncMap{
	"genDirectoryList": func(userDirRoot string, urlPath string) template.HTML {
		list, err := api.GenDirectoryList(userDirRoot, urlPath)
		if err != nil {
			errMsg := fmt.Sprintf("There has been an error getting directory list: %s", err.Error())
			return template.HTML(errMsg)
		}
		return list
	},
}
