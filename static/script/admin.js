"use strict";
var createUserForm = document.getElementById("create-user-form"),
    deleteUserForm = document.getElementById("delete-user-form");

// Submit create user form logic
createUserForm.addEventListener('submit', function (event) {
    event.preventDefault();

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/admin/create-user", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            displayErrorNotification(resp.error);
        } else if (this.status === 200) {
            // success
            displaySuccessNotification(resp.content);
            createUserForm.elements[0].value = "";
        }
    };
    console.log(serializeFormObj(createUserForm));
    xhr.send(JSON.stringify(serializeFormObj(createUserForm)));
});

// Submit delete user form logic
deleteUserForm.addEventListener('submit', function (event) {
    event.preventDefault();

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/admin/delete-user", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            displayErrorNotification(resp.error);
        } else if (this.status === 200) {
            // success
            displaySuccessNotification(resp.content);
            deleteUserForm.elements[0].value = "";
        }
    };
    xhr.send(JSON.stringify(serializeFormObj(deleteUserForm)));
});
