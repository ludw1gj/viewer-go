package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type indexData struct {
	List       template.HTML
	CurrentDir string
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method must be GET", http.StatusBadRequest)
		return
	}

	list, isFile, file, err := createDirectoryList(r.URL.Path)
	if err != nil {
		log.Println(err)
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
	}{
		list,
		strings.TrimPrefix(r.URL.Path, baseURL),
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

	path := r.URL.Query().Get("path") // can replace full file-name requirement in front-end
	err := createFolder(r.FormValue("folder-name"))
	if err != nil {
		log.Println("Could not create directory.", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
}

// fix
func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	file := "./test/" + strings.TrimPrefix(r.FormValue("file"), "/viewer/")
	fmt.Println(file)

	err := os.RemoveAll(file)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}

// no /viewer/ sent
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	dir := "./test" + strings.TrimPrefix(r.FormValue("directory"), "/viewer/")
	d, err := os.Open(dir)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
		return
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
			return
		}
	}
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}
