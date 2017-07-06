package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// getDirectoryList renders the directory list template according the directory path and returns the HTML document
// fragment.
func getDirectoryList(dirPath string, urlPath string) (list template.HTML, err error) {
	f, err := os.Open(dirPath)
	defer f.Close()
	if err != nil {
		return
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return
	}

	// entity holds information of either a file or directory.
	type entity struct {
		URL   string
		Name  string
		IsDir bool
	}

	// get directory list
	var entities []entity
	for _, file := range files {
		fileName := file.Name()
		fileURL := viewerRootURL + urlPath + "/" + fileName
		entities = append(entities, entity{fileURL, fileName, file.IsDir()})
	}

	index := true
	var previous bytes.Buffer
	fmt.Fprint(&previous, viewerRootURL)

	// get previous link if not at index
	if urlPath != "" {
		index = false

		urlSegments := strings.Split(urlPath, "/")
		count := len(urlSegments) - 1
		for i, segment := range urlSegments {
			if i == count {
				break
			}
			fmt.Fprintf(&previous, "%s/", segment)
		}
		if previous.String() != viewerRootURL {
			previous.Truncate(len(previous.String()) - 1)
		}
	}

	// directoryList holds information needed for executing the directory list template.
	type directoryList struct {
		Index       bool
		PreviousURL string
		Entities    []entity
		CurrentDir  string
	}

	var tplBuf bytes.Buffer
	err = dirListTpl.Execute(&tplBuf, directoryList{index, previous.String(), entities, urlPath})
	if err != nil {
		return
	}
	return template.HTML(tplBuf.String()), nil
}

// uploadFiles opens the FileHeader's associated Files, creates the destinations at the directory path and saves the
// files.
func uploadFiles(dirPath string, file map[string][]*multipart.FileHeader) error {
	for _, fileHeaders := range file {
		for _, hdr := range fileHeaders {
			// open uploaded files
			inFile, err := hdr.Open()
			if err != nil {
				return err
			}

			// open destination
			outFile, err := os.Create(dirPath + "/" + hdr.Filename)
			if err != nil {
				return err
			}

			// 32K buffer copy
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// createFolder creates a folder in the directory path.
func createFolder(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// deleteFile deletes the file at file path.
func deleteFile(filePath string) (err error) {
	err = os.RemoveAll(filePath)
	if err != nil {
		return
	}
	return
}

// deleteAllFiles deletes all files in the directory oath.
func deleteAllFiles(dirPath string) (err error) {
	d, err := os.Open(dirPath)
	if err != nil {
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dirPath, name))
		if err != nil {
			return
		}
	}
	return
}
