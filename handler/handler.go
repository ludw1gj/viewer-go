package handler

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"

	"io/ioutil"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, baseURL, http.StatusMovedPermanently)
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	isFile, err := renderIfFile(w, r)
	if err != nil {
		renderErrorPage(w, "", errors.New("There has been an error: "+err.Error()))
		return
	} else if isFile {
		return
	}

	path := mux.Vars(r)["path"]
	list, err := getDirectoryList(path)
	if err != nil {
		renderErrorPage(w, "", errors.New("There has been an error getting directory list: "+err.Error()))
	}
	data := struct {
		List       template.HTML
		CurrentDir string
	}{
		list,
		path,
	}
	err = viewerTpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func About(w http.ResponseWriter, _ *http.Request) {
	aboutTpl.Execute(w, nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {
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
	path := r.URL.Query().Get("path")
	err := deleteAllEntities(path)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
	}
	http.Redirect(w, r, baseURL+path, http.StatusMovedPermanently)
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	var buf bytes.Buffer
	err := notFoundTpl.Execute(&buf, nil)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(buf.Bytes())
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

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(buf.Bytes())
}

func renderIfFile(w http.ResponseWriter, r *http.Request) (isFile bool, err error) {
	path := mux.Vars(r)["path"]
	filePath := wrkDir + path

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if !fileInfo.IsDir() {
		isFile = true

		data, _ := ioutil.ReadFile(filePath)
		w.Header().Add("Content-Type", contentType(filePath))
		http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(data))
		return
	}
	return
}
