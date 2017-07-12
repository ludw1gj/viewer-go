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

uploadForm.addEventListener("submit", function (event) {
    event.preventDefault();

    var formData = new FormData(uploadForm);
    var xhr = new XMLHttpRequest();
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            displayErrorNotification(resp.error);
        } else if (this.status === 200) {
            // success, reload the page to see changes
            location.reload(true);
        }
    };
    xhr.open("POST", apiRoute + "upload", true);
    xhr.send(formData);
});

createFolderForm.addEventListener("submit", function (event) {
    event.preventDefault();
    viewerAjaxHelper({path: currentDir + createFolderForm.elements[0].value}, "create");
});

deleteFileFolderForm.addEventListener("submit", function (event) {
    event.preventDefault();
    viewerAjaxHelper({path: currentDir + deleteFileFolderForm.elements[0].value}, "delete");
});

deleteAllForm.addEventListener("submit", function (event) {
    event.preventDefault();
    viewerAjaxHelper({path: currentDir}, "delete-all");
});

/**
 * This function is for sending ajax for create/delete forms on viewer page
 * @param data {Object}
 * The data to be sent to the url, parsed to json
 * @param url
 * The apiRoute specific url to send the ajax request to
 */
function viewerAjaxHelper(data, url) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", apiRoute + url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            displayErrorNotification(resp.error);
        } else if (this.status === 200) {
            // success, reload the page to see changes
            location.reload(true);
        }
    };
    xhr.send(JSON.stringify(data));
}