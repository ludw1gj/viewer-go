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