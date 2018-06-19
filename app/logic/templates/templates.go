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
	siteTemplates = map[string]*template.Template{
		"login":      initSiteTemplate("login", false),
		"viewer":     initSiteTemplate("viewer", true),
		"about":      initSiteTemplate("about", true),
		"user":       initSiteTemplate("user", true),
		"admin":      initSiteTemplate("admin", true),
		"adminUsers": initSiteTemplate("admin_users", true),
		"error":      initSiteTemplate("error", true),
		"notFound":   initSiteTemplate("not_found", true),
	}

	directoryListTemplate = template.Must(template.ParseFiles(path.Join("views", "viewer", "dir_list.gohtml")))

	// function map for use in templates
	funcMap = template.FuncMap{
		"generateDirectoryList": func(userDirRoot string, urlPath string) template.HTML {
			list, err := generateDirectoryList(userDirRoot, urlPath, directoryListTemplate)
			if err != nil {
				errMsg := fmt.Sprintf("There has been an error getting directory list: %s", err.Error())
				return template.HTML(errMsg)
			}
			return list
		},
	}
)

// initSiteTemplate returns new templates.Template and parses files of templateName in the templates directory with base
// templates.
func initSiteTemplate(templateName string, baseTemplate bool) *template.Template {
	siteTemplateDir := path.Join("app", "views", "site")
	siteBaseTemplatePath := path.Join(siteTemplateDir, "base.gohtml")

	if baseTemplate {
		return template.Must(template.New(templateName).Funcs(funcMap).
			ParseFiles(siteBaseTemplatePath, path.Join(siteTemplateDir, templateName+".gohtml")))
	}
	return template.Must(template.ParseFiles(path.Join(siteTemplateDir, "login.gohtml")))
}

func renderError(w http.ResponseWriter, err error) {
	log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500: Server error"))
}

// RenderLoginTemplate renders the login template.
func RenderLoginTemplate(w http.ResponseWriter) {
	// Ensure the template exists in the map.
	tpl, ok := siteTemplates["login"]
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

// RenderSiteTemplate executes a templates and sends it to the client.
func RenderSiteTemplate(w http.ResponseWriter, name string, data interface{}) {
	// Ensure the template exists in the map.
	tpl, ok := siteTemplates[name]
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
