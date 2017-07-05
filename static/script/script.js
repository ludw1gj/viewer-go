// --- Util functions ---
/**
 * This function serialises a Form Element Object into a general Javascript Object
 * @param {Object} form
 * A DOM Form Element Object
 * @returns {Object}
 * A general Javascript Object
 */
function serializeFormObj(form) {
    var elems = form.elements;
    var obj = {};

    for (var i = 0; i < elems.length; i += 1) {
        var element = elems[i];
        var type = element.type;
        var name = element.name;
        var value = element.value;

        switch (type) {
            case "hidden":
            case "text":
                obj[name] = value;
                break;
            case "password":
                obj[name] = value;
            case "checkbox":
                if (value === "on") {
                    obj[name] = true;
                } else {
                    obj[name] = false;
                }
            case "number":
                obj[name] = parseInt(value);
            default:
                break;
        }
    }
    return obj;
}

/**
 * This function submits an ajax request of content-type json to the change-password route
 * @param {Object} data
 * A Form Element Object
 * @param {String} url
 * A URL to submit ajax to
 * @param {String} bannerId
 * The id of the message element
 */
function submitAjax(data, url, bannerId) {
    var request = new XMLHttpRequest();
    request.open('POST', url, true);

    request.onload = function () {
        var resp = JSON.parse(this.response);
        var userMessage = document.getElementById(bannerId);
        if (this.status === 200) {
            userMessage.classList.add("is-success");
            userMessage.style.display = "block";
            userMessage.getElementsByClassName("message-body")[0].innerText = resp.content;
        } else if (this.status === 401 || this.status === 500) {
            userMessage.classList.add("is-danger");
            userMessage.style.display = "block";
            userMessage.getElementsByClassName("message-body")[0].innerText = resp.error;
        }
        return resp
    };
    request.onerror = function () {
        console.log("There was a connection issue. Check your internet connection or the sever might be down.");
    };
    request.setRequestHeader('Content-Type', 'application/json');
    request.send(JSON.stringify(data));
}

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

        var data = serializeFormObj(changePasswordForm);
        submitAjax(data, "/api/user/change-password", "user-message");
    });
}

// --- Admin page ---
// Submit create user form logic
if (document.getElementById("create-user-form")) {
    var createUserForm = document.getElementById("create-user-form");
    createUserForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var data = serializeFormObj(createUserForm);
        submitAjax(data, "/api/admin/create-user", "admin-message");
    });
}

// Submit delete user form logic
if (document.getElementById("delete-user-form")) {
    var deleteUserForm = document.getElementById("delete-user-form");
    deleteUserForm.addEventListener('submit', function (event) {
        event.preventDefault();

        var data = serializeFormObj(deleteUserForm);
        submitAjax(data, "/api/admin/delete-user", "admin-message");
    });
}

