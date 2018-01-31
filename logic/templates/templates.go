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
	templateDir      = "views"
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

func renderError(w http.ResponseWriter, err error) {
	log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500: Server error"))
}

// RenderLoginTemplate renders the login template.
func RenderLoginTemplate(w http.ResponseWriter) {
	// Ensure the template exists in the map.
	tpl, ok := templates["login"]
	if !ok {
		renderError(w, errors.New("template does not exist"))
		return
	}

	var buf bytes.Buffer
	err := tpl.Execute(&buf, nil)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500: Server error. %s", err.Error())))
		return
	}
	w.Write(buf.Bytes())
}

// RenderTemplate executes a templates and sends it to the client.
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// Ensure the template exists in the map.
	tpl, ok := templates[name]
	if !ok {
		renderError(w, errors.New("template does not exist"))
		return
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "base.gohtml", data); err != nil {
		renderError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
