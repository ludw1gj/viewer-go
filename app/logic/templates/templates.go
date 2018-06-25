// Package templates contains template logic.
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
		"login":      createSiteTemplate("login", false),
		"viewer":     createSiteTemplate("viewer", true),
		"about":      createSiteTemplate("about", true),
		"user":       createSiteTemplate("user", true),
		"admin":      createSiteTemplate("admin", true),
		"adminUsers": createSiteTemplate("admin_users", true),
		"error":      createSiteTemplate("error", true),
		"notFound":   createSiteTemplate("not_found", true),
	}

	directoryListTemplate = template.Must(template.ParseFiles(path.Join(getTemplateDir(), "viewer", "dir_list.gohtml")))

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

func getTemplateDir() string {
	return path.Join("app", "views")
}

// createSiteTemplate returns new templates.Template and parses files of templateName in the
// templates directory with base templates.
func createSiteTemplate(templateName string, baseTemplate bool) *template.Template {
	siteTemplateDir := path.Join(getTemplateDir(), "site")
	siteBaseTemplatePath := path.Join(siteTemplateDir, "base.gohtml")

	if baseTemplate {
		return template.Must(template.New(templateName).Funcs(funcMap).
			ParseFiles(siteBaseTemplatePath, path.Join(siteTemplateDir, templateName+".gohtml")))
	}
	return template.Must(template.ParseFiles(path.Join(siteTemplateDir, "login.gohtml")))
}

func renderError(w http.ResponseWriter, err error) {
	log.Printf("StatusInternalServerError template failed to execute: %s", err.Error())
	errMessage := fmt.Sprintf("500: Server error - %s", err.Error())

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errMessage))
}

// renderTemplate executes a templates and sends it to the client.
func renderTemplate(w http.ResponseWriter, name string, baseTemplate string, data interface{}) {
	// Ensure the template exists in the map.
	tpl, ok := siteTemplates[name]
	if !ok {
		renderError(w, errors.New("template does not exist"))
		return
	}

	executeTemplate := func(b *bytes.Buffer) error {
		if baseTemplate != "" {
			return tpl.ExecuteTemplate(b, baseTemplate, data)
		}
		return tpl.Execute(b, data)
	}

	var buf bytes.Buffer
	if err := executeTemplate(&buf); err != nil {
		renderError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

// RenderLoginTemplate renders the login template, acts as a wrapper for func renderTemplate.
func RenderLoginTemplate(w http.ResponseWriter) {
	renderTemplate(w, "login", "", nil)
}

// RenderSiteTemplate renders a site template, acts as a wrapper for func renderTemplate.
func RenderSiteTemplate(w http.ResponseWriter, name string, data interface{}) {
	renderTemplate(w, name, "base.gohtml", data)
}
