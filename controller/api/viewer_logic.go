// This file contains logic for viewer functions.

package api

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var dirListTpl = template.Must(template.ParseFiles(path.Join("templates", "api", "dir_list.gohtml")))

// GenDirectoryList renders the directory list template according the directory path and returns the HTML document
// fragment.
func GenDirectoryList(userDirRoot string, urlPath string) (list template.HTML, err error) {
	// get items in directory
	f, err := os.Open(path.Join(userDirRoot, urlPath))
	if err != nil {
		return list, err
	}
	defer f.Close()
	items, err := f.Readdir(-1)
	if err != nil {
		return list, err
	}

	// sort items by name
	itemInfo := make(map[string]bool)
	for _, item := range items {
		itemInfo[item.Name()] = item.IsDir()
	}
	itemNamesSorted := make([]string, len(itemInfo))
	i := 0
	for itemName := range itemInfo {
		itemNamesSorted[i] = itemName
		i++
	}
	sort.Strings(itemNamesSorted)

	// entity holds information of either a file or directory.
	type entity struct {
		URL   string
		Name  string
		IsDir bool
	}

	// get directory list
	var entities []entity
	for _, itemName := range itemNamesSorted {
		var rawUrl string
		isDir := itemInfo[itemName]

		switch isDir {
		case true:
			rawUrl = "/viewer/" + urlPath + "/" + itemName
		case false:
			rawUrl = "/file/" + urlPath + "/" + itemName
		}
		itemURL, err := url.Parse(rawUrl)
		if err != nil {
			return list, err
		}
		entities = append(entities, entity{itemURL.String(), itemName, isDir})
	}

	// previous link
	index := true
	var previous bytes.Buffer
	fmt.Fprint(&previous, "/viewer/")
	// previous link if not at index
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
		// if previous url is not the index remove trailing slash
		if previous.String() != "/viewer/" {
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

	// execute and return the template
	var tplBuf bytes.Buffer
	err = dirListTpl.Execute(&tplBuf, directoryList{index, previous.String(), entities, urlPath})
	if err != nil {
		return list, err
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
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("File or Folder does not exist.")
	}

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
