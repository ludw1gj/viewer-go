// This file contains template logic.

package templates

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path"

	"fmt"
)

var (
	templateDir      = "view"
	baseTemplatePath = path.Join(templateDir, "base.gohtml")

	directoryListTemplate = template.Must(template.ParseFiles(path.Join(templateDir, "dir_list.gohtml")))

	templates = map[string]*template.Template{
		"login":      template.Must(template.ParseFiles(path.Join(templateDir, "login.gohtml"))),
		"viewer":     initTemplate("viewer"),
		"about":      initTemplate("about"),
		"user":       initTemplate("user"),
		"admin":      initTemplate("admin"),
		"adminUsers": initTemplate("admin_users"),
		"error":      initTemplate("error"),
		"notFound":   initTemplate("not_found"),
	}

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

// initTemplate returns new templates.Template and parses files of templateName in the templates directory with base
// templates.
func initTemplate(templateName string) *template.Template {
	return template.Must(template.New(templateName).Funcs(funcMap).
		ParseFiles(baseTemplatePath, path.Join(templateDir, templateName+".gohtml")))
}

// RenderTemplate executes a templates and sends it to the client.
func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	runError := func(err error) {
		log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Server error"))
	}

	// Ensure the template exists in the map.
	tpl, ok := templates[name]
	if !ok {
		runError(errors.New("template does not exist"))
		return
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "base.gohtml", data); err != nil {
		runError(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
