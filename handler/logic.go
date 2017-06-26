package handler

import (
	"bytes"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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

func createDirectoryList(pathURL string) (view template.HTML, isFile bool, file openedFile, err error) {
	trueFilePath := strings.Replace(pathURL, baseURL, wrkDir, -1)

	f, err := os.Open(trueFilePath)
	if err != nil {
		return
	}

	// check if path is a file
	fileInfo, _ := os.Stat(trueFilePath)
	if !fileInfo.IsDir() {
		return view, true, openedFile{f, trueFilePath}, nil
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return
	}
	f.Close()

	// get directory list
	var entities []entity
	for _, file := range files {
		fileName := file.Name()
		fileURL := pathURL + "/" + fileName
		entities = append(entities, entity{fileURL, fileName, file.IsDir()})
	}

	// get previous link if not at index of working directory
	index := true
	var previous string
	if trueFilePath != wrkDir {
		index = false
		urlParts := strings.Split(pathURL, "/")
		previous = strings.TrimSuffix(pathURL, "/"+urlParts[len(urlParts)-1])
	}

	var buf bytes.Buffer
	currentDir := strings.TrimPrefix(trueFilePath, wrkDir)
	err = dirListTpl.Execute(&buf, directoryList{index, previous, entities, currentDir})
	if err != nil {
		return
	}
	return template.HTML(buf.String()), false, file, nil
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
