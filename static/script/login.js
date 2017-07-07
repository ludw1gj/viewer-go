"use strict";

var loginForm = document.getElementById("login-form");

loginForm.addEventListener("submit", function (event) {
    event.preventDefault();

    var data = serializeFormObj(loginForm);
    submitLoginAjax(data);
});

/**
 * This function submits an ajax post request of content-type json to the login route
 * @param {Object} data
 * A Form Element Object
 */
function submitLoginAjax(data) {
    var request = new XMLHttpRequest();
    request.open("POST", "/login", true);

    request.onload = function () {
        var userMessage = document.getElementById("user-message");

        if (this.status === 200) {
            window.location = "/viewer/";
        } else if (this.status === 401 || this.status === 500) {
            userMessage.classList.add("is-danger");
            userMessage.style.display = "block";
            userMessage.getElementsByClassName("message-body")[0].innerText = JSON.parse(this.response).error;
        }
    };
    request.onerror = function () {
        console.log("There was a connection issue. Check your internet connection or the sever might be down.");
    };
    request.setRequestHeader('Content-Type', 'application/json');
    request.send(JSON.stringify(data));
}