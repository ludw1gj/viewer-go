define("Handler/NotificationHandler", ["require", "exports"], function (require, exports) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    function displayError(msg) {
        var notification = document.getElementById("notification");
        notification.classList.remove("is-success", "hidden");
        notification.classList.add("is-danger");
        notification.innerText = msg;
    }
    exports.displayError = displayError;
    function displaySuccess(msg) {
        var notification = document.getElementById("notification");
        notification.classList.remove("is-danger", "hidden");
        notification.classList.add("is-success");
        notification.innerText = msg;
    }
    exports.displaySuccess = displaySuccess;
});
define("Handler/AjaxHandler", ["require", "exports", "Handler/NotificationHandler"], function (require, exports, NotificationHandler_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    function ajaxSubmitJSON(url, data, errFunc, okFunc) {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
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
                    NotificationHandler_1.displayError("There has been an error.");
                }
            }
        };
        xhr.send(JSON.stringify(data));
    }
    exports.ajaxSubmitJSON = ajaxSubmitJSON;
    function ajaxSubmitFormData(url, uploadForm, errFunc, okFunc) {
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
                    NotificationHandler_1.displayError("There has been an error.");
                }
            }
        };
        xhr.send(formData);
    }
    exports.ajaxSubmitFormData = ajaxSubmitFormData;
});
define("Logic/LoginPageLogic", ["require", "exports", "Handler/AjaxHandler"], function (require, exports, AjaxHandler_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    function addEventListenerLoginForm() {
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
            AjaxHandler_1.ajaxSubmitJSON("/api/user/login", data, errFunc, okFunc);
        });
    }
    function initiateLoginPage() {
        addEventListenerLoginForm();
    }
    exports.initiateLoginPage = initiateLoginPage;
});
define("Logic/NavbarLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_2, AjaxHandler_2) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    function addEventListenerToMobileMenuButton() {
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
    }
    function addEventListenerToLogoutButton() {
        var logoutButton = document.getElementById("logout-button");
        logoutButton.addEventListener('click', function () {
            var errFunc = function (resp) {
                NotificationHandler_2.displayError(resp.error.message);
            };
            var okFunc = function () {
                window.location.href = "/login";
            };
            AjaxHandler_2.ajaxSubmitJSON("/api/user/logout", undefined, errFunc, okFunc);
        });
    }
    function initiateNavbar() {
        addEventListenerToMobileMenuButton();
        addEventListenerToLogoutButton();
    }
    exports.initiateNavbar = initiateNavbar;
});
define("Logic/ViewerPageLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_3, AjaxHandler_3) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var apiRoute = "/api/viewer/";
    var currentDir = document.getElementById("current-dir").innerText.slice(1);
    function addEventListenerToUploadFileForm() {
        var uploadForm = document.getElementById("upload-form");
        uploadForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var errFunc = function (resp) {
                NotificationHandler_3.displayError(resp.error.message);
            };
            var okFunc = function () {
                location.reload(true);
            };
            AjaxHandler_3.ajaxSubmitFormData(apiRoute + "upload/" + currentDir, uploadForm, errFunc, okFunc);
        });
    }
    function addEventListenerToCreateFolderForm() {
        var createFolderForm = document.getElementById("create-folder-form");
        createFolderForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var folderName = createFolderForm.folder_name;
            var data = {
                path: makePath(currentDir, folderName.value)
            };
            ajaxHelper(apiRoute + "create", data);
        });
    }
    function addEventListenerToDeleteFileFolderForm() {
        var deleteFileFolderForm = document.getElementById("delete-file-folder-form");
        deleteFileFolderForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var fileName = deleteFileFolderForm.file_name;
            var data = {
                path: makePath(currentDir, fileName.value)
            };
            ajaxHelper(apiRoute + "delete", data);
        });
    }
    function addEventListenerToDeleteAllForm() {
        var deleteAllForm = document.getElementById("delete-all-form");
        deleteAllForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var path;
            if (currentDir === "") {
                path = "/";
            }
            else {
                path = currentDir;
            }
            var data = {
                path: path
            };
            ajaxHelper(apiRoute + "delete-all", data);
        });
    }
    function ajaxHelper(url, data) {
        var errFunc = function (resp) {
            NotificationHandler_3.displayError(resp.error.message);
        };
        var okFunc = function () {
            location.reload(true);
        };
        AjaxHandler_3.ajaxSubmitJSON(url, data, errFunc, okFunc);
    }
    function makePath(currentDir, fileName) {
        var index = (currentDir === "");
        if (index) {
            return fileName;
        }
        return currentDir + "/" + fileName;
    }
    function initiateViewerPage() {
        addEventListenerToUploadFileForm();
        addEventListenerToCreateFolderForm();
        addEventListenerToDeleteFileFolderForm();
        addEventListenerToDeleteAllForm();
    }
    exports.initiateViewerPage = initiateViewerPage;
});
define("Logic/UserPageLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_4, AjaxHandler_4) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var userApiRoute = "/api/user/";
    function addEventListenerToChangeNameForm() {
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
                NotificationHandler_4.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                NotificationHandler_4.displaySuccess(resp.data.content);
                location.reload(true);
            };
            AjaxHandler_4.ajaxSubmitJSON(userApiRoute + "change-name", data, errFunc, okFunc);
        });
    }
    function addEventListenerToChangePasswordForm() {
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
                NotificationHandler_4.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                NotificationHandler_4.displaySuccess(resp.data.content);
                changePasswordForm.reset();
            };
            AjaxHandler_4.ajaxSubmitJSON(userApiRoute + "change-password", data, errFunc, okFunc);
        });
    }
    function addEventListenerToDeleteAccountForm() {
        var deleteAccountForm = document.getElementById("delete-account-form");
        deleteAccountForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var password = deleteAccountForm.password;
            var data = {
                password: password.value
            };
            var errFunc = function (resp) {
                NotificationHandler_4.displayError(resp.error.message);
                deleteAccountForm.reset();
            };
            var okFunc = function () {
                window.location.href = "/login";
            };
            AjaxHandler_4.ajaxSubmitJSON(userApiRoute + "delete", data, errFunc, okFunc);
        });
    }
    function initiateUserPage() {
        addEventListenerToChangeNameForm();
        addEventListenerToChangePasswordForm();
        addEventListenerToDeleteAccountForm();
    }
    exports.initiateUserPage = initiateUserPage;
});
define("Logic/AdminPageLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_5, AjaxHandler_5) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var adminApiRoute = "/api/admin/";
    function addEventListenerToChangeUsernameForm() {
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
                NotificationHandler_5.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                var username = document.getElementById("username");
                if (data.current_username === username.innerText) {
                    location.reload(true);
                    return;
                }
                NotificationHandler_5.displaySuccess(resp.data.content);
            };
            AjaxHandler_5.ajaxSubmitJSON(adminApiRoute + "change-username", data, errFunc, okFunc);
        });
    }
    function addEventListenerToChangeDirectoryRootForm() {
        var changeDirForm = document.getElementById("change-dir-root-form");
        changeDirForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var dirRoot = changeDirForm.dir_root;
            var data = {
                dir_root: dirRoot.value
            };
            var errFunc = function (resp) {
                NotificationHandler_5.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                NotificationHandler_5.displaySuccess(resp.data.content);
                changeDirForm.reset();
            };
            AjaxHandler_5.ajaxSubmitJSON(adminApiRoute + "change-dir-root", data, errFunc, okFunc);
        });
    }
    function addEventListenerToChangeAdminStatusForm() {
        var changeAdminStatusForm = document.getElementById("change-admin-status-form");
        changeAdminStatusForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var userID = changeAdminStatusForm.user_id;
            var isAdmin = changeAdminStatusForm.is_admin;
            var data = {
                user_id: parseInt(userID.value),
                is_admin: isAdmin.checked
            };
            var errFunc = function (resp) {
                NotificationHandler_5.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                NotificationHandler_5.displaySuccess(resp.data.content);
                changeAdminStatusForm.reset();
            };
            AjaxHandler_5.ajaxSubmitJSON(adminApiRoute + "change-admin-status", data, errFunc, okFunc);
        });
    }
    function addEventListenerToCreateUserForm() {
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
                NotificationHandler_5.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                NotificationHandler_5.displaySuccess(resp.data.content);
                createUserForm.reset();
            };
            AjaxHandler_5.ajaxSubmitJSON(adminApiRoute + "create-user", data, errFunc, okFunc);
        });
    }
    function addEventListenerToDeleteUserForm() {
        var deleteUserForm = document.getElementById("delete-user-form");
        deleteUserForm.addEventListener("submit", function (event) {
            event.preventDefault();
            var userID = deleteUserForm.user_id;
            var data = {
                user_id: parseInt(userID.value)
            };
            var errFunc = function (resp) {
                NotificationHandler_5.displayError(resp.error.message);
            };
            var okFunc = function (resp) {
                NotificationHandler_5.displaySuccess(resp.data.content);
                deleteUserForm.reset();
            };
            AjaxHandler_5.ajaxSubmitJSON(adminApiRoute + "delete-user", data, errFunc, okFunc);
        });
    }
    function initiateAdminPage() {
        addEventListenerToChangeUsernameForm();
        addEventListenerToChangeDirectoryRootForm();
        addEventListenerToChangeAdminStatusForm();
        addEventListenerToCreateUserForm();
        addEventListenerToDeleteUserForm();
    }
    exports.initiateAdminPage = initiateAdminPage;
});
define("LogicController", ["require", "exports", "Logic/LoginPageLogic", "Logic/NavbarLogic", "Logic/ViewerPageLogic", "Logic/UserPageLogic", "Logic/AdminPageLogic"], function (require, exports, LoginPageLogic_1, NavbarLogic_1, ViewerPageLogic_1, UserPageLogic_1, AdminPageLogic_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    function run() {
        var page = window.location.pathname;
        if (page === "/login") {
            LoginPageLogic_1.initiateLoginPage();
            return;
        }
        NavbarLogic_1.initiateNavbar();
        var isViewerPage = page.search("/viewer/") !== -1;
        if (isViewerPage) {
            ViewerPageLogic_1.initiateViewerPage();
            return;
        }
        switch (page) {
            case "/user":
                UserPageLogic_1.initiateUserPage();
                break;
            case "/admin":
                AdminPageLogic_1.initiateAdminPage();
                break;
        }
    }
    run();
});
