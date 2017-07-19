"use strict";
window.onload = function () {
    var apiRoute = "/api/viewer/";
    var currentDir = document.getElementById("current-dir").innerText.slice(1);
    var notification = document.getElementById("notification");
    var uploadForm = document.getElementById("upload-form"),
        createFolderForm = document.getElementById("create-folder-form"),
        deleteFileFolderForm = document.getElementById("delete-file-folder-form"),
        deleteAllForm = document.getElementById("delete-all-form");

    // handle upload form logic
    uploadForm.addEventListener("submit", function (event) {
        event.preventDefault();

        var errFunc = function (resp) {
            displayErrorNotification(notification, resp.error);
        };
        var okFunc = function () {
            location.reload(true);
        };
        submitAjaxFormData(apiRoute + "upload", uploadForm, errFunc, okFunc);
    });

    // handle create folder form logic
    createFolderForm.addEventListener("submit", function (event) {
        event.preventDefault();
        viewerAjaxHelper(apiRoute + "create", {path: makePath(currentDir, createFolderForm.elements[0].value)}, notification);
    });

    // handle delete file/folder form logic
    deleteFileFolderForm.addEventListener("submit", function (event) {
        event.preventDefault();
        viewerAjaxHelper(apiRoute + "delete", {path: makePath(currentDir, deleteFileFolderForm.elements[0].value)}, notification);
    });

    // handle delete all form logic
    deleteAllForm.addEventListener("submit", function (event) {
        event.preventDefault();
        viewerAjaxHelper(apiRoute + "delete-all", {path: currentDir}, notification);
    });
};

/**
 * This function is a wrapper for submitAjaxJson function
 * @param url
 * The apiRoute specific url to send the ajax request to
 * @param data {Object}
 * The data to be sent to the url, parsed to json
 * @param {Object} notifElm
 * The notification element
 */
function viewerAjaxHelper(url, data, notifElm) {
    var errFunc = function (resp) {
        displayErrorNotification(notifElm, resp.error);
    };
    var okFunc = function () {
        location.reload(true);
    };
    submitAjaxJson(url, data, errFunc, okFunc)
}

/**
 * This function returns path
 * @param currentDir {String}
 * The current directory
 * @param fileName {String}
 * File name to append to directory
 * @return {String}
 */
function makePath(currentDir, fileName) {
    var index = currentDir === "";
    if (index) {
        return fileName;
    } else {
        console.log(currentDir + "/" + fileName);
        return currentDir + "/" + fileName;
    }
}