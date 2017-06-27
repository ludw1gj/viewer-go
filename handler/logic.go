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
		fileURL := baseURL + path + "/" + fileName
		entities = append(entities, entity{fileURL, fileName, file.IsDir()})
	}

	index := true
	var previous bytes.Buffer
	fmt.Fprint(&previous, baseURL)

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
		if previous.String() != baseURL {
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

func processMultipartFormFiles(path string, file map[string][]*multipart.FileHeader) error {
	truePath := wrkDir + path

	for _, fileHeaders := range file {
		for _, hdr := range fileHeaders {
			// open uploaded
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

func createFolder(path string) error {
	truePath := wrkDir + path
	err := os.MkdirAll(truePath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func deleteEntity(filePath string) (err error) {
	err = os.RemoveAll(filePath)
	if err != nil {
		return
	}
	return
}

func deleteAllEntities(path string) (err error) {
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
