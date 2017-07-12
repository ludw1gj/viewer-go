"use strict";
var notification = document.getElementById("notification"),
    mobileMenu = document.getElementById("mobile-menu"),
    mobileMenuButton = document.getElementById("mobile-menu-button"),
    logoutButton = document.getElementById("logout-button");

mobileMenuButton.addEventListener("click", function () {
    if (mobileMenuButton.classList.contains("is-active") || mobileMenuButton.classList.contains("is-active")) {
        mobileMenu.classList.remove("is-active");
        mobileMenuButton.classList.remove("is-active");
    } else {
        mobileMenuButton.classList.add("is-active");
        mobileMenu.classList.add("is-active");
    }
});

logoutButton.addEventListener('click', function () {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/user/logout", true);
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
    xhr.send();
});
