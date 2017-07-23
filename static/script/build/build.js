"use strict";
// initAdminPage function should be run at initialisation of admin page.
function initAdminPage() {
    var adminApiRoute = "/api/admin/";
    // handle create user form logic
    var adminCreateUserForm = document.getElementById("create-user-form");
    adminCreateUserForm.addEventListener('submit', function (event) {
        event.preventDefault();
        var url = adminApiRoute + "create-user";
        var data = {
            username: adminCreateUserForm.username.value,
            password: adminCreateUserForm.password.value,
            first_name: adminCreateUserForm.first_name.value,
            last_name: adminCreateUserForm.last_name.value,
            directory_root: adminCreateUserForm.directory_root.value,
            is_admin: adminCreateUserForm.is_admin.checked
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });
    // handle delete user form logic
    var adminDeleteUserForm = document.getElementById("delete-user-form");
    adminDeleteUserForm.addEventListener('submit', function (event) {
        event.preventDefault();
        var url = adminApiRoute + "delete-user";
        var data = {
            user_id: parseInt(adminDeleteUserForm.user_id.value)
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            adminDeleteUserForm.user_id.value = "";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });
}
// submitAjaxJson submits an AJAX POST request.
function submitAjaxJson(url, data, errFunc, okFunc) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        var DONE = 4;
        if (xhr.readyState === DONE) {
            var resp = JSON.parse(xhr.responseText);
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp);

            }
            else if ("data" in resp) {
                okFunc(resp);

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

            }
            else if ("data" in resp) {
                okFunc(resp);

            }
            else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(formData);
}
// initBasePage function should be run at initialisation of base page.
function initBasePage() {
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
    var baseLogoutButton = document.getElementById("logout-button");
    baseLogoutButton.addEventListener('click', function () {
        var url = "/api/user/logout";
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function () {
            window.location.href = "/login";
        };
        submitAjaxJson(url, undefined, errFunc, okFunc);
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
        initBasePage();
    }
    if (page.search("/viewer/") !== -1) {
        initViewerPage();
        return;
    }
    switch (page) {
        case "/user":
            initUserPage();
            break;
        case "/admin":
            initAdminPage();
            break;
    }
}
// run init.
loadAuthorizedPages();
// initLoginPage function should be run at initialisation of login page.
function initLoginPage() {
    // handle login user form logic
    var loginForm = document.getElementById("login-form");
    loginForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var url = "/api/user/login";
        var data = {
            username: loginForm.username.value,
            password: loginForm.password.value
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
        submitAjaxJson(url, data, errFunc, okFunc);
    });
}
// loadLoginPage loads login page script if at login page.
function loadLoginPage() {
    if (window.location.pathname === "/login") {
        initLoginPage();
    }
}
// run init.
loadLoginPage();
// initUserPage function should be run at initialisation of user page.
function initUserPage() {
    var userApiRoute = "/api/user/";
    // handle change password form logic
    var userChangePasswordForm = document.getElementById("change-password-form");
    userChangePasswordForm.addEventListener('submit', function (event) {
        event.preventDefault();
        var url = userApiRoute + "change-password";
        var oldPw = userChangePasswordForm.elements.item(0);
        var newPw = userChangePasswordForm.elements.item(1);
        var data = {
            old_password: oldPw.value,
            new_password: newPw.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function (resp) {
            displaySuccessNotification(resp.data.content);
            oldPw.value = "";
            newPw.value = "";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });
    // handle delete user form logic
    var userDeleteAccountForm = document.getElementById("delete-account-form");
    userDeleteAccountForm.addEventListener('submit', function (event) {
        event.preventDefault();
        var url = userApiRoute + "delete";
        var pw = userDeleteAccountForm.elements.item(0);
        var data = {
            password: pw.value
        };
        var errFunc = function (resp) {
            displayErrorNotification(resp.error.message);
        };
        var okFunc = function () {
            window.location.href = "/login";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });
}
// initViewerPage function should be run at initialisation of viewer page.
function initViewerPage() {
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
        var folderName = createFolderForm.folder_name.value;
        var data = {
            path: makePath(currentDir, folderName)
        };
        viewerAjaxHelper(apiRoute + "create", data);
    });
    // handle delete file/folder form logic
    var deleteFileFolderForm = document.getElementById("delete-file-folder-form");
    deleteFileFolderForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var fileName = deleteFileFolderForm.file_name.value;
        var data = {
            path: makePath(currentDir, fileName)
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
