// --- Nav-bar ---
var burgerIcon = document.getElementById("burger-icon");
var burgerMenu = document.getElementById("burger-menu");
var content = document.getElementById("content");

burgerIcon.addEventListener("click", function () {
    if (burgerIcon.classList.contains("is-active") || burgerMenu.classList.contains("is-active")) {
        burgerIcon.classList.remove("is-active");
        burgerMenu.classList.remove("is-active");
    } else {
        burgerIcon.classList.add("is-active");
        burgerMenu.classList.add("is-active");
    }
}, false);

content.addEventListener("click", function () {
    if (burgerIcon.classList.contains("is-active") || burgerMenu.classList.contains("is-active")) {
        burgerIcon.classList.remove("is-active");
        burgerMenu.classList.remove("is-active");
    }
}, false);


// --- User page ---
// Submit change password form logic
if (document.getElementById("change-password-form")) {
    var changePasswordForm = document.getElementById("change-password-form");
    changePasswordForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var oldPassword = document.getElementById("change-password-form-old-password").value;
        var newPassword = document.getElementById("change-password-form-new-password").value;

        var data = {old_password: oldPassword, new_password: newPassword};
        var url = "/api/user/change-password";
        submitAjax(data, url);
    });
}

/**
 * This function submits an ajax request of content-type json to a given url
 * @param {Object} data
 * * A Form Element which has been serialized by serializeFormObj function
 * @param {String} url
 * A url to send ajax request to
 */
function submitAjax(data, url) {
    var request = new XMLHttpRequest();
    request.open('POST', url, true);

    request.onload = function () {
        var resp = JSON.parse(this.response);
        var userMessage = document.getElementById("user-message");

        if (this.status === 200) {
            userMessage.classList.add("is-success");
            userMessage.style.display = "block";
            userMessage.getElementsByClassName("message-body")[0].innerText = resp.content;
        } else if (this.status === 401 || this.status === 500) {
            userMessage.classList.add("is-danger");
            userMessage.style.display = "block";
            userMessage.getElementsByClassName("message-body")[0].innerText = resp.error;
        }
    };
    request.onerror = function () {
        console.log("There was a connection issue. Check your internet connection or the sever might be down.");
    };
    request.setRequestHeader('Content-Type', 'application/json');
    request.send(JSON.stringify(data));
}
