"use strict";
// addEventListenersAdminForms function should be run at initialisation of admin page.
function addEventListenersAdminForms() {
    var adminApiRoute = "/api/admin/";
    // handle change directory root form logic
    var changeUsernameForm = document.getElementById("change-username-form");
    changeUsernameForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var currentUsername = changeUsernameForm.current_username;
        var newUsername = changeUsernameForm.new_username;
        var data = {
            current_username: currentUsername.value,
            new_username: newUsername.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            var username = document.getElementById("username");
            if (data.current_username === username.innerText) {
                location.reload(true);
                return;
            }
            displaySuccessNotification(resp.data.content);
        };
        submitAjaxJson(adminApiRoute + "change-username", data, errFunc, okFunc);
    });
    // handle change directory root form logic
    var changeDirForm = document.getElementById("change-dir-root-form");
    changeDirForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var dirRoot = changeDirForm.dir_root;
        var data = {
            dir_root: dirRoot.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            changeDirForm.reset();
        };
        submitAjaxJson(adminApiRoute + "change-dir-root", data, errFunc, okFunc);
    });
    // handle create user form logic
    var createUserForm = document.getElementById("create-user-form");
    createUserForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var username = createUserForm.username;
        var password = createUserForm.password;
        var firstName = createUserForm.first_name;
        var lastName = createUserForm.last_name;
        var DirRoot = createUserForm.directory_root;
        var isAdmin = createUserForm.is_admin;
        var data = {
            username: username.value,
            password: password.value,
            first_name: firstName.value,
            last_name: lastName.value,
            directory_root: DirRoot.value,
            is_admin: isAdmin.checked
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            createUserForm.reset();
        };
        submitAjaxJson(adminApiRoute + "create-user", data, errFunc, okFunc);
    });
    // handle delete user form logic
    var deleteUserForm = document.getElementById("delete-user-form");
    deleteUserForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var userID = deleteUserForm.user_id;
        var data = {
            user_id: parseInt(userID.value)
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            createUserForm.reset();
        };
        submitAjaxJson(adminApiRoute + "delete-user", data, errFunc, okFunc);
    });
}
// submitAjaxJson submits an AJAX POST request.
function submitAjaxJson(url, data, errFunc, okFunc) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        var DONE = 4;
        if (xhr.readyState === DONE) {
            var resp = JSON.parse(xhr.responseText);
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp);
                return;
            }
            else if ("data" in resp) {
                okFunc(resp);
                return;
            }
            else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(JSON.stringify(data));
}
// submitAjaxFormData uploads files via AJAX.
function submitAjaxFormData(url, uploadForm, errFunc, okFunc) {
    var formData = new FormData(uploadForm);
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.onreadystatechange = function () {
        var DONE = 4;
        if (xhr.readyState === DONE) {
            var resp = JSON.parse(xhr.responseText);
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp);
                return;
            }
            else if ("data" in resp) {
                okFunc(resp);
                return;
            }
            else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(formData);
}
// addEventListenersBaseNav function should be run at initialisation of base page.
function addEventListenersBaseNav() {
    // extend and collapse navigation menu for mobile
    var mobileMenuButton = document.getElementById("mobile-menu-button");
    mobileMenuButton.addEventListener("click", function () {
        var mobileMenu = document.getElementById("mobile-menu");
        if (mobileMenuButton.classList.contains("is-active") || mobileMenuButton.classList.contains("is-active")) {
            mobileMenu.classList.remove("is-active");
            mobileMenuButton.classList.remove("is-active");
        }
        else {
            mobileMenuButton.classList.add("is-active");
            mobileMenu.classList.add("is-active");
        }
    });
    // handle logout user
    var logoutButton = document.getElementById("logout-button");
    logoutButton.addEventListener('click', function () {
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function () {
            window.location.href = "/login";
        };
        submitAjaxJson("/api/user/logout", undefined, errFunc, okFunc);
    });
}
// displayErrorNotification displays error notification.
function displayErrorNotification(msg) {
    var notification = document.getElementById("notification");
    notification.classList.remove("is-success", "hidden");
    notification.classList.add("is-danger");
    notification.innerText = msg;
}
// displaySuccessNotification displays success notification.
function displaySuccessNotification(msg) {
    var notification = document.getElementById("notification");
    notification.classList.remove("is-danger", "hidden");
    notification.classList.add("is-success");
    notification.innerText = msg;
}
// load authorized page's script.
function loadAuthorizedPages() {
    var page = window.location.pathname;
    if (page !== "/login") {
        addEventListenersBaseNav();
    }
    if (page.search("/viewer/") !== -1) {
        addEventListenersViewerForms();
        return;
    }
    switch (page) {
        case "/user":
            addEventListenersUserForms();
            break;
        case "/admin":
            addEventListenersAdminForms();
            break;
    }
}
// run init.
loadAuthorizedPages();
// addEventListenersLoginForm function should be run at initialisation of login page.
function addEventListenerLoginForm() {
    // handle login user form logic
    var loginForm = document.getElementById("login-form");
    loginForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var username = loginForm.username;
        var password = loginForm.password;
        var data = {
            username: username.value,
            password: password.value
        };
        var errFunc = function (resp) {
            var notification = document.getElementById("login-error-notification");
            notification.classList.remove("hidden");
            notification.classList.add("is-danger");
            notification.innerText = resp.error.message;
        };
        var okFunc = function () {
            window.location.href = "/viewer/";
        };
        submitAjaxJson("/api/user/login", data, errFunc, okFunc);
    });
}
// loadLoginPage loads login page script if at login page.
function loadLoginPageScript() {
    if (window.location.pathname === "/login") {
        addEventListenerLoginForm();
    }
}
// run init.
loadLoginPageScript();
// addEventListenersUserForms function should be run at initialisation of user page.
function addEventListenersUserForms() {
    var userApiRoute = "/api/user/";
    // handle change name form logic
    var changeNameForm = document.getElementById("change-name-form");
    changeNameForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var firstName = changeNameForm.first_name;
        var lastName = changeNameForm.last_name;
        var data = {
            first_name: firstName.value,
            last_name: lastName.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            location.reload(true);
        };
        submitAjaxJson(userApiRoute + "change-name", data, errFunc, okFunc);
    });
    // handle change password form logic
    var changePasswordForm = document.getElementById("change-password-form");
    changePasswordForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var oldPassword = changePasswordForm.old_password;
        var newPassword = changePasswordForm.new_password;
        var data = {
            old_password: oldPassword.value,
            new_password: newPassword.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            changePasswordForm.reset();
        };
        submitAjaxJson(userApiRoute + "change-password", data, errFunc, okFunc);
    });
    // handle delete user form logic
    var deleteAccountForm = document.getElementById("delete-account-form");
    deleteAccountForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var password = deleteAccountForm.password;
        var data = {
            password: password.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
            deleteAccountForm.reset();
        };
        var okFunc = function () {
            window.location.href = "/login";
        };
        submitAjaxJson(userApiRoute + "delete", data, errFunc, okFunc);
    });
}
// addEventListenersViewerForms function should be run at initialisation of viewer page.
function addEventListenersViewerForms() {
    var apiRoute = "/api/viewer/";
    var currentDir = document.getElementById("current-dir").innerText.slice(1);
    // handle upload form logic
    var uploadForm = document.getElementById("upload-form");
    uploadForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function () {
            location.reload(true);
        };
        submitAjaxFormData(apiRoute + "upload", uploadForm, errFunc, okFunc);
    });
    // handle create folder form logic
    var createFolderForm = document.getElementById("create-folder-form");
    createFolderForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var folderName = createFolderForm.folder_name;
        var data = {
            path: makePath(currentDir, folderName.value)
        };
        viewerAjaxHelper(apiRoute + "create", data);
    });
    // handle delete file/folder form logic
    var deleteFileFolderForm = document.getElementById("delete-file-folder-form");
    deleteFileFolderForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var fileName = deleteFileFolderForm.file_name;
        var data = {
            path: makePath(currentDir, fileName.value)
        };
        viewerAjaxHelper(apiRoute + "delete", data);
    });
    // handle delete all form logic
    var deleteAllForm = document.getElementById("delete-all-form");
    deleteAllForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var data = {
            path: currentDir
        };
        viewerAjaxHelper(apiRoute + "delete-all", data);
    });
}
// viewerAjaxHelper is a wrapper for submitAjaxJson function.
function viewerAjaxHelper(url, data) {
    var errFunc = function (resp) {
        displayErrorNotification(resp.error.message);
    };
    var okFunc = function () {
        location.reload(true);
    };
    submitAjaxJson(url, data, errFunc, okFunc);
}
// makePath generates a path.
function makePath(currentDir, fileName) {
    var index = (currentDir === "");
    if (index) {
        return fileName;
    }
    else {
        return currentDir + "/" + fileName;
    }
}
