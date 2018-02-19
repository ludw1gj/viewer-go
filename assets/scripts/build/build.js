define("Handler/NotificationHandler", ["require", "exports"], function (require, exports) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var NotificationHandler = (function () {
        function NotificationHandler() {
        }
        NotificationHandler.displayError = function (msg) {
            var notification = document.getElementById("notification");
            notification.classList.remove("is-success", "hidden");
            notification.classList.add("is-danger");
            notification.innerText = msg;
        };
        NotificationHandler.displaySuccess = function (msg) {
            var notification = document.getElementById("notification");
            notification.classList.remove("is-danger", "hidden");
            notification.classList.add("is-success");
            notification.innerText = msg;
        };
        return NotificationHandler;
    }());
    exports.NotificationHandler = NotificationHandler;
});
define("Handler/AjaxHandler", ["require", "exports", "Handler/NotificationHandler"], function (require, exports, NotificationHandler_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var AjaxHandler = (function () {
        function AjaxHandler() {
        }
        AjaxHandler.submitJSON = function (url, data, errFunc, okFunc) {
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
                        NotificationHandler_1.NotificationHandler.displayError("There has been an error.");
                    }
                }
            };
            xhr.send(JSON.stringify(data));
        };
        AjaxHandler.submitFormData = function (url, uploadForm, errFunc, okFunc) {
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
                        NotificationHandler_1.NotificationHandler.displayError("There has been an error.");
                    }
                }
            };
            xhr.send(formData);
        };
        return AjaxHandler;
    }());
    exports.AjaxHandler = AjaxHandler;
});
define("Logic/LoginPageLogic", ["require", "exports", "Handler/AjaxHandler"], function (require, exports, AjaxHandler_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var LoginPageLogic = (function () {
        function LoginPageLogic() {
            this.addEventListenerLoginForm();
        }
        LoginPageLogic.prototype.addEventListenerLoginForm = function () {
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
                AjaxHandler_1.AjaxHandler.submitJSON("/api/user/login", data, errFunc, okFunc);
            });
        };
        return LoginPageLogic;
    }());
    exports.LoginPageLogic = LoginPageLogic;
});
define("Logic/NavbarLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_2, AjaxHandler_2) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var NavbarLogic = (function () {
        function NavbarLogic() {
            this.addEventListenerToMobileMenuButton();
            this.addEventListenerToLogoutButton();
        }
        NavbarLogic.prototype.addEventListenerToMobileMenuButton = function () {
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
        };
        NavbarLogic.prototype.addEventListenerToLogoutButton = function () {
            var logoutButton = document.getElementById("logout-button");
            logoutButton.addEventListener('click', function () {
                var errFunc = function (resp) {
                    NotificationHandler_2.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function () {
                    window.location.href = "/login";
                };
                AjaxHandler_2.AjaxHandler.submitJSON("/api/user/logout", undefined, errFunc, okFunc);
            });
        };
        return NavbarLogic;
    }());
    exports.NavbarLogic = NavbarLogic;
});
define("Logic/ViewerPageLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_3, AjaxHandler_3) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var ViewerPageLogic = (function () {
        function ViewerPageLogic() {
            this.apiRoute = "/api/viewer/";
            this.currentDir = document.getElementById("current-dir").innerText.slice(1);
            this.addEventListenerToUploadFileForm();
            this.addEventListenerToCreateFolderForm();
            this.addEventListenerToDeleteFileFolderForm();
            this.addEventListenerToDeleteAllForm();
        }
        ViewerPageLogic.prototype.addEventListenerToUploadFileForm = function () {
            var _this = this;
            var uploadForm = document.getElementById("upload-form");
            uploadForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var errFunc = function (resp) {
                    NotificationHandler_3.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function () {
                    location.reload(true);
                };
                AjaxHandler_3.AjaxHandler.submitFormData(_this.apiRoute + "upload/" + _this.currentDir, uploadForm, errFunc, okFunc);
            });
        };
        ViewerPageLogic.prototype.addEventListenerToCreateFolderForm = function () {
            var _this = this;
            var createFolderForm = document.getElementById("create-folder-form");
            createFolderForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var folderName = createFolderForm.folder_name;
                var data = {
                    path: ViewerPageLogic.makePath(_this.currentDir, folderName.value)
                };
                ViewerPageLogic.ajaxHelper(_this.apiRoute + "create", data);
            });
        };
        ViewerPageLogic.prototype.addEventListenerToDeleteFileFolderForm = function () {
            var _this = this;
            var deleteFileFolderForm = document.getElementById("delete-file-folder-form");
            deleteFileFolderForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var fileName = deleteFileFolderForm.file_name;
                var data = {
                    path: ViewerPageLogic.makePath(_this.currentDir, fileName.value)
                };
                ViewerPageLogic.ajaxHelper(_this.apiRoute + "delete", data);
            });
        };
        ViewerPageLogic.prototype.addEventListenerToDeleteAllForm = function () {
            var _this = this;
            var deleteAllForm = document.getElementById("delete-all-form");
            deleteAllForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var path;
                if (_this.currentDir === "") {
                    path = "/";
                }
                else {
                    path = _this.currentDir;
                }
                var data = {
                    path: path
                };
                ViewerPageLogic.ajaxHelper(_this.apiRoute + "delete-all", data);
            });
        };
        ViewerPageLogic.ajaxHelper = function (url, data) {
            var errFunc = function (resp) {
                NotificationHandler_3.NotificationHandler.displayError(resp.error.message);
            };
            var okFunc = function () {
                location.reload(true);
            };
            AjaxHandler_3.AjaxHandler.submitJSON(url, data, errFunc, okFunc);
        };
        ViewerPageLogic.makePath = function (currentDir, fileName) {
            var index = (currentDir === "");
            if (index) {
                return fileName;
            }
            return currentDir + "/" + fileName;
        };
        return ViewerPageLogic;
    }());
    exports.ViewerPageLogic = ViewerPageLogic;
});
define("Logic/UserPageLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_4, AjaxHandler_4) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var UserPageLogic = (function () {
        function UserPageLogic() {
            this.userApiRoute = "/api/user/";
            this.addEventListenerToChangeNameForm();
            this.addEventListenerToChangePasswordForm();
            this.addEventListenerToDeleteAccountForm();
        }
        UserPageLogic.prototype.addEventListenerToChangeNameForm = function () {
            var _this = this;
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
                    NotificationHandler_4.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    NotificationHandler_4.NotificationHandler.displaySuccess(resp.data.content);
                    location.reload(true);
                };
                AjaxHandler_4.AjaxHandler.submitJSON(_this.userApiRoute + "change-name", data, errFunc, okFunc);
            });
        };
        UserPageLogic.prototype.addEventListenerToChangePasswordForm = function () {
            var _this = this;
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
                    NotificationHandler_4.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    NotificationHandler_4.NotificationHandler.displaySuccess(resp.data.content);
                    changePasswordForm.reset();
                };
                AjaxHandler_4.AjaxHandler.submitJSON(_this.userApiRoute + "change-password", data, errFunc, okFunc);
            });
        };
        UserPageLogic.prototype.addEventListenerToDeleteAccountForm = function () {
            var _this = this;
            var deleteAccountForm = document.getElementById("delete-account-form");
            deleteAccountForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var password = deleteAccountForm.password;
                var data = {
                    password: password.value
                };
                var errFunc = function (resp) {
                    NotificationHandler_4.NotificationHandler.displayError(resp.error.message);
                    deleteAccountForm.reset();
                };
                var okFunc = function () {
                    window.location.href = "/login";
                };
                AjaxHandler_4.AjaxHandler.submitJSON(_this.userApiRoute + "delete", data, errFunc, okFunc);
            });
        };
        return UserPageLogic;
    }());
    exports.UserPageLogic = UserPageLogic;
});
define("Logic/AdminPageLogic", ["require", "exports", "Handler/NotificationHandler", "Handler/AjaxHandler"], function (require, exports, NotificationHandler_5, AjaxHandler_5) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var AdminPageLogic = (function () {
        function AdminPageLogic() {
            this.adminApiRoute = "/api/admin/";
            this.addEventListenerToChangeUsernameForm();
            this.addEventListenerToChangeDirectoryRootForm();
            this.addEventListenerToChangeAdminStatusForm();
            this.addEventListenerToCreateUserForm();
            this.addEventListenerToDeleteUserForm();
        }
        AdminPageLogic.prototype.addEventListenerToChangeUsernameForm = function () {
            var _this = this;
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
                    NotificationHandler_5.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    var username = document.getElementById("username");
                    if (data.current_username === username.innerText) {
                        location.reload(true);
                        return;
                    }
                    NotificationHandler_5.NotificationHandler.displaySuccess(resp.data.content);
                };
                AjaxHandler_5.AjaxHandler.submitJSON(_this.adminApiRoute + "change-username", data, errFunc, okFunc);
            });
        };
        AdminPageLogic.prototype.addEventListenerToChangeDirectoryRootForm = function () {
            var _this = this;
            var changeDirForm = document.getElementById("change-dir-root-form");
            changeDirForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var dirRoot = changeDirForm.dir_root;
                var data = {
                    dir_root: dirRoot.value
                };
                var errFunc = function (resp) {
                    NotificationHandler_5.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    NotificationHandler_5.NotificationHandler.displaySuccess(resp.data.content);
                    changeDirForm.reset();
                };
                AjaxHandler_5.AjaxHandler.submitJSON(_this.adminApiRoute + "change-dir-root", data, errFunc, okFunc);
            });
        };
        AdminPageLogic.prototype.addEventListenerToChangeAdminStatusForm = function () {
            var _this = this;
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
                    NotificationHandler_5.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    NotificationHandler_5.NotificationHandler.displaySuccess(resp.data.content);
                    changeAdminStatusForm.reset();
                };
                AjaxHandler_5.AjaxHandler.submitJSON(_this.adminApiRoute + "change-admin-status", data, errFunc, okFunc);
            });
        };
        AdminPageLogic.prototype.addEventListenerToCreateUserForm = function () {
            var _this = this;
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
                    NotificationHandler_5.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    NotificationHandler_5.NotificationHandler.displaySuccess(resp.data.content);
                    createUserForm.reset();
                };
                AjaxHandler_5.AjaxHandler.submitJSON(_this.adminApiRoute + "create-user", data, errFunc, okFunc);
            });
        };
        AdminPageLogic.prototype.addEventListenerToDeleteUserForm = function () {
            var _this = this;
            var deleteUserForm = document.getElementById("delete-user-form");
            deleteUserForm.addEventListener("submit", function (event) {
                event.preventDefault();
                var userID = deleteUserForm.user_id;
                var data = {
                    user_id: parseInt(userID.value)
                };
                var errFunc = function (resp) {
                    NotificationHandler_5.NotificationHandler.displayError(resp.error.message);
                };
                var okFunc = function (resp) {
                    NotificationHandler_5.NotificationHandler.displaySuccess(resp.data.content);
                    deleteUserForm.reset();
                };
                AjaxHandler_5.AjaxHandler.submitJSON(_this.adminApiRoute + "delete-user", data, errFunc, okFunc);
            });
        };
        return AdminPageLogic;
    }());
    exports.AdminPageLogic = AdminPageLogic;
});
define("LogicController", ["require", "exports", "Logic/LoginPageLogic", "Logic/NavbarLogic", "Logic/ViewerPageLogic", "Logic/UserPageLogic", "Logic/AdminPageLogic"], function (require, exports, LoginPageLogic_1, NavbarLogic_1, ViewerPageLogic_1, UserPageLogic_1, AdminPageLogic_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var LogicController = (function () {
        function LogicController() {
        }
        LogicController.run = function () {
            var page = window.location.pathname;
            if (page === "/login") {
                new LoginPageLogic_1.LoginPageLogic();
                return;
            }
            new NavbarLogic_1.NavbarLogic();
            var isViewerPage = page.search("/viewer/") !== -1;
            if (isViewerPage) {
                new ViewerPageLogic_1.ViewerPageLogic();
                return;
            }
            switch (page) {
                case "/user":
                    new UserPageLogic_1.UserPageLogic();
                    break;
                case "/admin":
                    new AdminPageLogic_1.AdminPageLogic();
                    break;
            }
        };
        return LogicController;
    }());
    LogicController.run();
});
