package handler

import (
	"bytes"
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

const (
	baseURL = "/viewer/"
	wrkDir  = "./test/"
)

type openedFile struct {
	File         *os.File
	TrueFilePath string
}

type entity struct {
	URL   string
	Name  string
	IsDir bool
}

type directoryList struct {
	Index       bool
	PreviousURL string
	Entities    []entity
	CurrentDir  string
}

var (
	viewerTpl  *template.Template
	dirListTpl *template.Template
)

func init() {
	viewerTpl = template.Must(template.ParseFiles(path.Join("templates", "index.gohtml")))
	dirListTpl = template.Must(template.ParseFiles(path.Join("templates", "view.gohtml")))
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/viewer/", http.StatusMovedPermanently)
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	list, isFile, file := createDirectoryList(r.URL.Path)
	defer file.File.Close()
	if isFile {
		w.Header().Add("Content-Type", contentType(file.TrueFilePath))
		http.ServeContent(w, r, file.TrueFilePath, time.Now(), file.File)
		return
	}

	err := viewerTpl.Execute(w, list)
	if err != nil {
		log.Println(err)
	}
}

func createDirectoryList(pathURL string) (view template.HTML, isFile bool, file openedFile) {
	filePath := strings.TrimPrefix(pathURL, "/viewer/")
	trueFilePath := wrkDir + filePath

	f, err := os.Open(trueFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	// check if path is a file
	fileInfo, _ := os.Stat(trueFilePath)
	if !fileInfo.IsDir() {
		return view, true, openedFile{f, trueFilePath}
	}

	files, err := f.Readdir(-1)
	if err != nil {
		log.Println(err)
		return
	}
	f.Close()

	// get directory list
	var entities []entity
	for _, file := range files {
		fileName := file.Name()
		fileURL := baseURL + filePath + "/" + fileName
		entities = append(entities, entity{fileURL, fileName, file.IsDir()})
	}

	// get previous link if not at index of working directory
	index := true
	var previous string
	if filePath != "" {
		index = false
		urlParts := strings.Split(pathURL, "/")
		previous = strings.TrimSuffix(pathURL, "/"+urlParts[len(urlParts)-1])
	}

	var buf bytes.Buffer
	err = dirListTpl.Execute(&buf, directoryList{index, previous, entities, filePath})
	if err != nil {
		log.Println(err)
	}
	return template.HTML(buf.String()), false, file
}

// fix
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method must be POST", http.StatusBadRequest)
		return
	}

	folderName := r.FormValue("folder-name")
	err := os.MkdirAll(folderName, os.ModePerm)
	if err != nil {
		log.Println("Could not create directory.", err)
		return
	}

}

// fix
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
	http.Redirect(w, r, "/viewer/", http.StatusTemporaryRedirect)
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
