"use strict";
var loginForm = document.getElementById("login-form"),
    notification = document.getElementById("notification");

loginForm.addEventListener("submit", function (event) {
    event.preventDefault();
    var data = {
        username: loginForm.elements[0].value,
        password: loginForm.elements[1].value
    };

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/user/login", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            // error, display error in notification
            notification.classList.remove("is-success", "hidden");
            notification.classList.add("is-danger");
            notification.innerText = resp.error;
        } else if (this.status === 200) {
            // success, redirect to the viewer page
            window.location = "/viewer/";
        }
    };
    xhr.send(JSON.stringify(data));
});
