package handler

import (
	"bytes"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"fmt"
)

// entity holds information of either a file or directory.
type entity struct {
	URL   string
	Name  string
	IsDir bool
}

// directoryList holds information needed for executing the directory list template.
type directoryList struct {
	Index       bool
	PreviousURL string
	Entities    []entity
	CurrentDir  string
}

// getDirectoryList renders the directory list template according the path provided.
func getDirectoryList(path string) (list template.HTML, err error) {
	trueFilePath := wrkDir + path

	f, err := os.Open(trueFilePath)
	defer f.Close()
	if err != nil {
		return
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return
	}

	// get directory list
	var entities []entity
	for _, file := range files {
		fileName := file.Name()
		fileURL := viewerRootURL + path + "/" + fileName
		entities = append(entities, entity{fileURL, fileName, file.IsDir()})
	}

	index := true
	var previous bytes.Buffer
	fmt.Fprint(&previous, viewerRootURL)

	// get previous link if not at index
	if path != "" {
		index = false

		urlSegments := strings.Split(path, "/")
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

	var tplBuf bytes.Buffer
	err = dirListTpl.Execute(&tplBuf, directoryList{index, previous.String(), entities, path})
	if err != nil {
		return
	}
	return template.HTML(tplBuf.String()), nil
}

// processMultipartFileHeaders opens the FileHeader's associated Files, creates the destinations at the path provided
// and saves the files.
func processMultipartFileHeaders(path string, file map[string][]*multipart.FileHeader) error {
	truePath := wrkDir + path

	for _, fileHeaders := range file {
		for _, hdr := range fileHeaders {
			// open uploaded files
			inFile, err := hdr.Open()
			if err != nil {
				return err
			}

			// open destination
			outFile, err := os.Create(truePath + "/" + hdr.Filename)
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

// createFolder creates a folder at provided path in the working directory.
func createFolder(path string) error {
	truePath := wrkDir + path
	err := os.MkdirAll(truePath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// deleteFile deletes the file at path/fileName provided in the working directory.
func deleteFile(path string, fileName string) (err error) {
	filePath := wrkDir + path + "/" + fileName
	err = os.RemoveAll(filePath)
	if err != nil {
		return
	}
	return
}

// deleteAllFiles deletes all files in the path provided in the working directory.
func deleteAllFiles(path string) (err error) {
	dir := wrkDir + path

	d, err := os.Open(dir)
	if err != nil {
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return
		}
	}
	return
}
