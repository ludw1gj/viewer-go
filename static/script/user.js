"use strict";
window.onload = function () {
    var userApiRoute = "/api/user/",
        userChangePasswordForm = document.getElementById("change-password-form"),
        userDeleteAccountForm = document.getElementById("delete-account-form");

    // handle change password form logic
    userChangePasswordForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var url = userApiRoute + "change-password";
        var data = serializeFormObj(userChangePasswordForm);
        var errFunc = function (resp) {
            displayErrorNotification(resp.error);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.content);
            userChangePasswordForm.elements[0].value = "";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });

    // handle delete user form logic
    userDeleteAccountForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var url = userApiRoute + "delete";
        var data = serializeFormObj(userDeleteAccountForm);
        var errFunc = function (resp) {
            displayErrorNotification(resp.error);
        };
        var okFunc = function () {
            window.location = "/login";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });

};
