"use strict";
var changePasswordForm = document.getElementById("change-password-form"),
    deleteAccountForm = document.getElementById("delete-account-form"),
    apiRoute = "/api/user/";

// Submit change password form logic
changePasswordForm.addEventListener('submit', function (event) {
    event.preventDefault();

    var xhr = new XMLHttpRequest();
    xhr.open("POST", apiRoute + "change-password", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            displayErrorNotification(resp.error);
        } else if (this.status === 200) {
            // success
            displaySuccessNotification(resp.content);
            changePasswordForm.elements[0].value = "";
        }
    };
    xhr.send(JSON.stringify(serializeFormObj(changePasswordForm)));
});

// Submit delete user form logic
deleteAccountForm.addEventListener('submit', function (event) {
    event.preventDefault();

    var xhr = new XMLHttpRequest();
    xhr.open("POST", apiRoute + "delete", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            displayErrorNotification(resp.error);
        } else if (this.status === 200) {
            // success, redirect to login page
            window.location = "/login";
        }
    };
    xhr.send(JSON.stringify(serializeFormObj(deleteUserForm)));
});
