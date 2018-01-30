package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/robertjeffs/viewer-go/controller/common"
	"github.com/robertjeffs/viewer-go/database"
)

// ViewerPage handles the viewer page. It uses the path variable in the route to determine which directory in the user's
// directory in the filesystem to display a directory list for.
func ViewerPage(w http.ResponseWriter, r *http.Request) {
	user, err := common.ValidateUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// urlPath should not contain a leading /
	urlPath := strings.TrimPrefix(mux.Vars(r)["path"], "/")
	data := struct {
		CurrentDir string
		User       database.User
	}{
		urlPath,
		user,
	}
	renderTemplate(w, r, viewerTpl, data)
}

// genDirectoryList renders the directory list template according the directory path and returns the HTML document
// fragment.
func genDirectoryList(userDirRoot string, urlPath string) (list template.HTML, err error) {
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

	// check if directory is empty
	isEmpty := false
	if len(items) == 0 {
		isEmpty = true
	}

	// directoryList holds information needed for executing the directory list template.
	type directoryList struct {
		Index       bool
		PreviousURL string
		Entities    []entity
		CurrentDir  string
		IsEmpty     bool
	}

	// execute and return the template
	var tplBuf bytes.Buffer
	if err := dirListTpl.Execute(&tplBuf, directoryList{index, previous.String(), entities, urlPath, isEmpty}); err != nil {
		return list, err
	}
	return template.HTML(tplBuf.String()), nil
}

// SendFile sends file to client.
func SendFile(w http.ResponseWriter, r *http.Request) {
	// get user from session
	user, err := common.ValidateUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// path to file
	filePath := path.Join(user.DirectoryRoot, mux.Vars(r)["path"])

	// get file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}
	if fileInfo.IsDir() {
		renderErrorPage(w, r, errors.New("requested item is not a file"))
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		renderErrorPage(w, r, err)
		return
	}
	w.Header().Add("Content-Type", contentType(filePath))
	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(data))
}

// contentType determines the content-type by the file extension of the file at the path.
func contentType(path string) (contentType string) {
	hasSuffix := func(suffix string) bool {
		return strings.HasSuffix(path, suffix)
	}

	if hasSuffix(".css") {
		return "text/css"
	} else if hasSuffix(".js") {
		return "application/javascript"
	} else if hasSuffix(".png") {
		return "image/png"
	} else if hasSuffix(".jpg") {
		return "image/jpeg"
	} else if hasSuffix(".jpeg") {
		return "image/jpeg"
	} else if hasSuffix(".mp4") {
		return "video/mp4"
	}
	return "text/plain"
}
