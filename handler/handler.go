package handler

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, baseURL, http.StatusMovedPermanently)
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method must be GET", http.StatusBadRequest)
		return
	}

	list, err := getDirectoryList(w, r)
	data := struct {
		List       template.HTML
		CurrentDir string
		Error      error
	}{
		list,
		strings.TrimPrefix(r.URL.Path, baseURL),
		err,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = viewerTpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	path := r.URL.Query().Get("path")

	// parse request
	const _24K = (1 << 10) * 24
	err := r.ParseMultipartForm(_24K)
	if err != nil {
		renderErrorPage(w, path, err)
		return
	}

	err = processMultipartFormFiles(path, r.MultipartForm.File)
	if err != nil {
		renderErrorPage(w, path, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
}

func CreateFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	path := r.URL.Query().Get("path")
	folderName := r.FormValue("folder-name")
	folderPath := path + "/" + folderName

	err := createFolder(folderPath)
	if err != nil {
		renderErrorPage(w, path, errors.New("Could not create directory: "+err.Error()))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	path := r.URL.Query().Get("path")
	fileName := r.FormValue("file-name")
	if fileName == "" {
		renderErrorPage(w, path, errors.New("File name cannot be empty."))
	}
	file := wrkDir + path + "/" + fileName

	err := deleteEntity(file)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
	}
	http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	path := r.URL.Query().Get("path")
	err := deleteAllEntities(path)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
	}
	http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
}

func renderErrorPage(w http.ResponseWriter, path string, err error) {
	page := baseURL + path

	data := struct {
		Error error
		Page  string
	}{
		err,
		page,
	}
	var buf bytes.Buffer
	err = errorTpl.Execute(&buf, data)
	if err != nil {
		log.Println(err)
	}
	w.Write(buf.Bytes())
}
