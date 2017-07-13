"use strict";
var adminApiRoute = "/api/admin/",
    adminCreateUserForm = document.getElementById("create-user-form"),
    adminDeleteUserForm = document.getElementById("delete-user-form");

// handle create user form logic
adminCreateUserForm.addEventListener('submit', function (event) {
    event.preventDefault();

    var url = adminApiRoute + "create-user";
    var data = serializeFormObj(adminCreateUserForm);
    var errFunc = function (resp) {
        displayErrorNotification(resp.error);
    };
    var okFunc = function (resp) {
        displaySuccessNotification(resp.content);
        adminCreateUserForm.elements[0].value = "";
    };
    submitAjaxJson(url, data, errFunc, okFunc);
});

// handle delete user form logic
adminDeleteUserForm.addEventListener('submit', function (event) {
    event.preventDefault();

    var url = adminApiRoute + "delete-user";
    var data = serializeFormObj(adminDeleteUserForm);
    var errFunc = function (resp) {
        displayErrorNotification(resp.error);
    };
    var okFunc = function (resp) {
        displaySuccessNotification(resp.content);
        adminDeleteUserForm.elements[0].value = "";
    };
    submitAjaxJson(url, data, errFunc, okFunc);
});
