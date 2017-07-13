"use strict";
var notification = document.getElementById("notification"),
    loginForm = document.getElementById("login-form");

// handle login user form logic
loginForm.addEventListener("submit", function (event) {
    event.preventDefault();

    var url = "/api/user/login";
    var data = {
        username: loginForm.elements[0].value,
        password: loginForm.elements[1].value
    };
    var errFunc = function (resp) {
        notification.classList.remove("is-success", "hidden");
        notification.classList.add("is-danger");
        notification.innerText = resp.error;
    };
    var okFunc = function () {
        window.location = "/viewer/";
    };
    submitAjaxJson(url, data, errFunc, okFunc);
});
