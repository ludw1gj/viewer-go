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
	templateDir      = path.Join("view", "frontend")
	baseTemplatePath = path.Join(templateDir, "base.gohtml")

	// templates
	loginTemplate   = template.Must(template.ParseFiles(path.Join(templateDir, "login.gohtml")))
	dirListTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "dir_list.gohtml")))

	// view that require base template
	viewerTemplate     = initTemplate("viewer")
	aboutTemplate      = initTemplate("about")
	userTemplate       = initTemplate("user")
	adminTemplate      = initTemplate("admin")
	adminUsersTemplate = initTemplate("admin_users")
	errorTemplate      = initTemplate("error")
	notFoundTemplate   = initTemplate("not_found")

	// function map for use in templates
	funcMap = template.FuncMap{
		"generateDirectoryList": func(userDirRoot string, urlPath string) template.HTML {
			list, err := generateDirectoryList(userDirRoot, urlPath)
			if err != nil {
				errMsg := fmt.Sprintf("There has been an error getting directory list: %s", err.Error())
				return template.HTML(errMsg)
			}
			return list
		},
	}
)

// initTemplate returns new template.Template and parses files of templateName in the template directory with base
// template.
func initTemplate(templateName string) *template.Template {
	return template.Must(template.New(templateName).Funcs(funcMap).
		ParseFiles(baseTemplatePath, path.Join(templateDir, templateName+".gohtml")))
}

// renderTemplate executes a template and sends it to the client.
func renderTemplate(w http.ResponseWriter, r *http.Request, template *template.Template, data interface{}) {
	var templateBuf bytes.Buffer
	if err := template.ExecuteTemplate(&templateBuf, "base.gohtml", data); err != nil {
		log.Println(err)

		renderErrorPage(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(templateBuf.Bytes())
}
