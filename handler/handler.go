package handler

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type entity struct {
	URL  string
	Name string
}

type view struct {
	Index    bool
	Previous string
	Entities []entity
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles(
		path.Join("templates", "index.gohtml"), path.Join("templates", "view.gohtml")))
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	filePath := "./test/" + strings.TrimPrefix(r.URL.Path, "/viewer/")
	f, err := os.Open(filePath)
	if err != nil {
		log.Println("Viewer Hanlder:", err)
		return
	}
	defer f.Close()

	// not a directory
	fileInfo, _ := os.Stat(filePath)
	if !fileInfo.IsDir() {
		w.Header().Add("Content-Type", contentType(filePath))
		http.ServeContent(w, r, filePath, time.Now(), f)
		return
	}

	// is a directory
	dirs, err := f.Readdir(-1)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	// remove "/" from end of url
	if string(r.URL.Path[len(r.URL.Path)-1]) == "/" {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	}

	var entities []entity
	for _, d := range dirs {
		name := d.Name()
		url := r.URL.Path + "/" + name

		if d.IsDir() {
			name += "/"
		}
		entities = append(entities, entity{url, name})
	}

	// get previous link
	index := true
	var previous string
	if r.URL.Path != "/viewer" {
		index = false
		urlParts := strings.Split(r.URL.Path, "/")
		previous = strings.TrimSuffix(r.URL.Path, urlParts[len(urlParts)-1])
	}
	tpl.Execute(w, view{index, previous, entities})
}

func CreateFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	newFolder := "./test/" + r.FormValue("folder-name")

	err := os.MkdirAll(newFolder, os.ModePerm)
	if err != nil {
		log.Println("Could not create directory.")
		http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
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

	for _, fileHeaders := range r.MultipartForm.File {
		for _, hdr := range fileHeaders {
			// open uploaded
			inFile, err := hdr.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// open destination
			outFile, err := os.Create("./uploaded/" + hdr.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// 32K buffer copy
			written, err := io.Copy(outFile, inFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
		}
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	file := "./test/" + strings.TrimPrefix(r.FormValue("file"), "/viewer/")
	fmt.Println(file)

	err := os.RemoveAll(file)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	dir := "./test" + strings.TrimPrefix(r.FormValue("directory"), "/viewer/")
	d, err := os.Open(dir)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
		return
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
			return
		}
	}
	http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
}
