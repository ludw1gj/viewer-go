// TODO: refactor js
// remove location.reload(true); and replace with reloading ajax dir list

"use strict";
var currentDir = window.location.pathname.replace("/viewer/", "") + "/",
    uploadForm = document.getElementById("upload-form"),
    createFolderForm = document.getElementById("create-folder-form"),
    deleteFileFolderForm = document.getElementById("delete-file-folder-form"),
    deleteAllForm = document.getElementById("delete-all-form"),
    apiRoute = "/api/viewer/";

if (currentDir[0] === "/") {
    // there must not be a leading slash
    currentDir = currentDir.slice(1);
}

// handle upload form logic
uploadForm.addEventListener("submit", function (event) {
    event.preventDefault();

    var url = apiRoute + "upload";
    var errFunc = function (resp) {
        displayErrorNotification(resp.error);
    };
    var okFunc = function () {
        location.reload(true);
    };
    submitAjaxFormData(url, uploadForm, errFunc, okFunc);
});

// handle create folder form logic
createFolderForm.addEventListener("submit", function (event) {
    event.preventDefault();
    viewerAjaxHelper({path: currentDir + createFolderForm.elements[0].value}, "create");
});

// handle delete file/folder form logic
deleteFileFolderForm.addEventListener("submit", function (event) {
    event.preventDefault();
    viewerAjaxHelper({path: currentDir + deleteFileFolderForm.elements[0].value}, "delete");
});

// handle delete all form logic
deleteAllForm.addEventListener("submit", function (event) {
    event.preventDefault();
    viewerAjaxHelper({path: currentDir}, "delete-all");
});

/**
 * This function is a wrapper for submitAjaxJson function
 * @param data {Object}
 * The data to be sent to the url, parsed to json
 * @param url
 * The apiRoute specific url to send the ajax request to
 */
function viewerAjaxHelper(data, url) {
    var errFunc = function (resp) {
        displayErrorNotification(resp.error);
    };
    var okFunc = function () {
        location.reload(true);
    };
    submitAjaxJson(apiRoute + url, data, errFunc, okFunc)
}

// TODO: doc
function travelDirTree(path) {
    var cntDir = document.getElementById("current-dir");
    var dirList = document.getElementById("dir-list");

    if (path === "") {
        cntDir.innerText = "/";
    } else {
        cntDir.innerText = path;
    }
    getDirList(cntDir, dirList);
}

// TODO: doc
function init() {
    var cntDir = document.getElementById("current-dir");
    var dirList = document.getElementById("dir-list");
    getDirList(cntDir, dirList);
}

init();

// TODO: doc
function getDirList(currentDir, dirList) {
    // clear dirList
    dirList.innerHTML = "";

    var errTestFunc = function (resp) {
        displayErrorNotification(resp.error);
    };
    var okTestFunc = function (resp) {
        console.log(resp);

        var backToIndex = document.getElementById("back-index");

        // if not at directory root, add back link and show back to index
        if (document.getElementById("current-dir").innerText !== "/") {
            var li = document.createElement("li");
            var a = document.createElement("a");
            a.id = "back-link";
            a.innerText = "../";

            a.addEventListener("click", function (event) {
                event.preventDefault();

                var previousPath = "";
                var array = currentDir.innerText.split("/");
                console.log("array 1", array);

                array.shift();
                console.log("array 2", array);

                for (var i = 0; i < array.length - 1; i++) {
                    previousPath += "/" + array[i];
                }

                console.log("previousPath", previousPath);
                travelDirTree(previousPath);
            });

            li.appendChild(a);
            dirList.appendChild(li);

            backToIndex.className = "";
        } else {
            backToIndex.className = "hidden";
        }

        // if items, create list
        if (resp.items !== null) {
            resp.items.forEach(function (item) {
                var li = document.createElement("li");
                var a = document.createElement("a");
                a.innerText = item.name;
                a.href = item.url;

                a.addEventListener("click", function (event) {
                    if (item.is_file) {
                        return
                    }
                    event.preventDefault();
                    travelDirTree(item.url);
                });

                li.appendChild(a);
                dirList.appendChild(li);
            });
        }
    };
    submitAjaxJson("/api/test", {path: currentDir.innerText}, errTestFunc, okTestFunc);
}