package model

import (
	"bytes"
	"html/template"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/FriedPigeon/viewer-go/config"
	"github.com/FriedPigeon/viewer-go/tpl"
)

type User struct {
	ID           int
	Username     string
	HashPassword string
}

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

// GetDirectoryList renders the directory list template according the path provided.
func (u User) GetDirectoryList(path string) (list template.HTML, err error) {
	trueFilePath := config.WrkDir + path

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
		fileURL := config.ViewerRootURL + path + "/" + fileName
		entities = append(entities, entity{fileURL, fileName, file.IsDir()})
	}

	index := true
	var previous bytes.Buffer
	fmt.Fprint(&previous, config.ViewerRootURL)

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
		if previous.String() != config.ViewerRootURL {
			previous.Truncate(len(previous.String()) - 1)
		}
	}

	var tplBuf bytes.Buffer
	err = tpl.DirListTpl.Execute(&tplBuf, directoryList{index, previous.String(), entities, path})
	if err != nil {
		return
	}
	return template.HTML(tplBuf.String()), nil
}

// ProcessMultipartFileHeaders opens the FileHeader's associated Files, creates the destinations at the path provided
// and saves the files.
func (u User) ProcessMultipartFileHeaders(path string, file map[string][]*multipart.FileHeader) error {
	truePath := config.WrkDir + path

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

// CreateFolder creates a folder at provided path in the working directory.
func (u User) CreateFolder(path string) error {
	truePath := config.WrkDir + path
	err := os.MkdirAll(truePath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile deletes the file at path/fileName provided in the working directory.
func (u User) DeleteFile(path string, fileName string) (err error) {
	filePath := config.WrkDir + path + "/" + fileName
	err = os.RemoveAll(filePath)
	if err != nil {
		return
	}
	return
}

// DeleteAllFiles deletes all files in the path provided in the working directory.
func (u User) DeleteAllFiles(path string) (err error) {
	dir := config.WrkDir + path

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
