package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

// generateDirectoryList renders the directory list templates according the directory path and returns the HTML document
// fragment.
func generateDirectoryList(userDirRoot string, urlPath string) (list template.HTML, err error) {
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
		var rawURL string
		isDir := itemInfo[itemName]

		switch isDir {
		case true:
			rawURL = "/viewer/" + urlPath + "/" + itemName
		case false:
			rawURL = "/file/" + urlPath + "/" + itemName
		}
		itemURL, err := url.Parse(rawURL)
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

	// check if directory is empty
	isEmpty := false
	if len(items) == 0 {
		isEmpty = true
	}

	// directoryList holds information needed for executing the directory list templates.
	type directoryList struct {
		Index       bool
		PreviousURL string
		Entities    []entity
		CurrentDir  string
		IsEmpty     bool
	}

	// execute and return the templates
	var templateBuf bytes.Buffer
	if err := directoryListTemplate.Execute(&templateBuf, directoryList{index, previous.String(), entities, urlPath, isEmpty}); err != nil {
		return list, err
	}
	return template.HTML(templateBuf.String()), nil
}
