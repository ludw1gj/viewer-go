package handler

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type indexData struct {
	List       template.HTML
	CurrentDir string
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, baseURL, http.StatusMovedPermanently)
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method must be GET", http.StatusBadRequest)
		return
	}

	list, isFile, file, err := createDirectoryList(r.URL.Path)
	if err != nil {
		err = errors.New("There has been an error getting directory list: " + err.Error())
	}

	defer file.File.Close()
	if isFile {
		w.Header().Set("Content-Type", contentType(file.TrueFilePath))
		http.ServeContent(w, r, file.TrueFilePath, time.Now(), file.File)
		return
	}

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

	// parse request
	const _24K = (1 << 10) * 24
	err := r.ParseMultipartForm(_24K)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.URL.Query().Get("path")
	err = processMultipartFormFiles(path, r.MultipartForm.File)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Println("Could not create directory.", err)
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
		log.Println("File name cannot be empty")
		http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
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
