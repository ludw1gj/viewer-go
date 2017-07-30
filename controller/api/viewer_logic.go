// This file contains logic for viewer functions.

package api

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// createFolder creates a folder in the directory path.
func createFolder(dirPath string) error {
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return errors.New("Folder already exists.")
	}

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
		return err
	}
	return nil
}

// deleteAllFiles deletes all files in the directory oath.
func deleteAllFiles(dirPath string) (err error) {
	d, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dirPath, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// saveFiles opens the FileHeader's associated Files, creates the destinations at the directory path and saves the
// files.
func saveFiles(dirPath string, file map[string][]*multipart.FileHeader) error {
	for _, fileHeaders := range file {
		for _, hdr := range fileHeaders {
			// open files
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
