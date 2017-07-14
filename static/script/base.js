"use strict";
var notification = document.getElementById("notification"),
    baseMobileMenu = document.getElementById("mobile-menu"),
    baseMobileMenuButton = document.getElementById("mobile-menu-button"),
    baseLogoutButton = document.getElementById("logout-button");

// extend and collapse navigation menu for mobile
baseMobileMenuButton.addEventListener("click", function () {
    if (baseMobileMenuButton.classList.contains("is-active") || baseMobileMenuButton.classList.contains("is-active")) {
        baseMobileMenu.classList.remove("is-active");
        baseMobileMenuButton.classList.remove("is-active");
    } else {
        baseMobileMenuButton.classList.add("is-active");
        baseMobileMenu.classList.add("is-active");
    }
});

// handle logout user
baseLogoutButton.addEventListener('click', function () {
    var url = "/api/user/logout";
    var data = undefined;

    var errFunc = function (resp) {
        displayErrorNotification(resp.error);
    };

    var okFunc = function () {
        window.location = "/login";
    };
    submitAjaxJson(url, data, errFunc, okFunc);
});

/**
 * This function displays error notification.
 * @param {String} msg
 * A message to display in the notification
 */
function displayErrorNotification(msg) {
    notification.classList.remove("is-success", "hidden");
    notification.classList.add("is-danger");
    notification.innerText = msg;
}

/**
 * This function displays success notification.
 * @param {String} msg
 * A message to display in the notification
 */
function displaySuccessNotification(msg) {
    notification.classList.remove("is-danger", "hidden");
    notification.classList.add("is-success");
    notification.innerText = msg;
}
