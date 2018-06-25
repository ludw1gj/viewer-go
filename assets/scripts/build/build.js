System.register("app/bin/display", [], function (exports_1, context_1) {
    "use strict";
    var displayError, displaySuccess;
    var __moduleName = context_1 && context_1.id;
    return {
        setters: [],
        execute: function () {
            exports_1("displayError", displayError = function (msg) {
                var notification = document.getElementById("notification");
                notification.classList.remove("is-success", "hidden");
                notification.classList.add("is-danger");
                notification.innerText = msg;
            });
            exports_1("displaySuccess", displaySuccess = function (msg) {
                var notification = document.getElementById("notification");
                notification.classList.remove("is-danger", "hidden");
                notification.classList.add("is-success");
                notification.innerText = msg;
            });
        }
    };
});
System.register("app/bin/ajax", ["app/bin/display"], function (exports_2, context_2) {
    "use strict";
    var display_1, submitAjaxJSON, submitAjaxFormData;
    var __moduleName = context_2 && context_2.id;
    return {
        setters: [
            function (display_1_1) {
                display_1 = display_1_1;
            }
        ],
        execute: function () {
            exports_2("submitAjaxJSON", submitAjaxJSON = function (url, data, errFunc, okFunc) {
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
                            display_1.displayError("There has been an error.");
                        }
                    }
                };
                xhr.send(JSON.stringify(data));
            });
            exports_2("submitAjaxFormData", submitAjaxFormData = function (url, uploadForm, errFunc, okFunc) {
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
                            display_1.displayError("There has been an error.");
                        }
                    }
                };
                xhr.send(formData);
            });
        }
    };
});
System.register("app/lib/loginPage", ["app/bin/ajax"], function (exports_3, context_3) {
    "use strict";
    var ajax_1, addEventListenerLoginForm, addEventListenersLoginPage;
    var __moduleName = context_3 && context_3.id;
    return {
        setters: [
            function (ajax_1_1) {
                ajax_1 = ajax_1_1;
            }
        ],
        execute: function () {
            addEventListenerLoginForm = function () {
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
                    ajax_1.submitAjaxJSON("/api/user/login", data, errFunc, okFunc);
                });
            };
            exports_3("addEventListenersLoginPage", addEventListenersLoginPage = function () {
                addEventListenerLoginForm();
            });
        }
    };
});
System.register("app/lib/navbar", ["app/bin/display", "app/bin/ajax"], function (exports_4, context_4) {
    "use strict";
    var display_2, ajax_2, addEventListenerToMobileMenuButton, addEventListenerToLogoutButton, addEventListenersNavbar;
    var __moduleName = context_4 && context_4.id;
    return {
        setters: [
            function (display_2_1) {
                display_2 = display_2_1;
            },
            function (ajax_2_1) {
                ajax_2 = ajax_2_1;
            }
        ],
        execute: function () {
            addEventListenerToMobileMenuButton = function () {
                var mobileMenuButton = document.getElementById("mobile-menu-button");
                mobileMenuButton.addEventListener("click", function () {
                    var mobileMenu = document.getElementById("mobile-menu");
                    if (mobileMenuButton.classList.contains("is-active") ||
                        mobileMenuButton.classList.contains("is-active")) {
                        mobileMenu.classList.remove("is-active");
                        mobileMenuButton.classList.remove("is-active");
                    }
                    else {
                        mobileMenuButton.classList.add("is-active");
                        mobileMenu.classList.add("is-active");
                    }
                });
            };
            addEventListenerToLogoutButton = function () {
                var logoutButton = document.getElementById("logout-button");
                logoutButton.addEventListener("click", function () {
                    var errFunc = function (resp) {
                        display_2.displayError(resp.error.message);
                    };
                    var okFunc = function () {
                        window.location.href = "/login";
                    };
                    ajax_2.submitAjaxJSON("/api/user/logout", undefined, errFunc, okFunc);
                });
            };
            exports_4("addEventListenersNavbar", addEventListenersNavbar = function () {
                addEventListenerToMobileMenuButton();
                addEventListenerToLogoutButton();
            });
        }
    };
});
System.register("app/lib/viewerPage", ["app/bin/display", "app/bin/ajax"], function (exports_5, context_5) {
    "use strict";
    var display_3, ajax_3, apiRoute, getCurrentDir, addEventListenerToUploadFileForm, addEventListenerToCreateFolderForm, addEventListenerToDeleteFileFolderForm, addEventListenerToDeleteAllForm, ajaxHelper, makePath, addEventListenersViewerPage;
    var __moduleName = context_5 && context_5.id;
    return {
        setters: [
            function (display_3_1) {
                display_3 = display_3_1;
            },
            function (ajax_3_1) {
                ajax_3 = ajax_3_1;
            }
        ],
        execute: function () {
            apiRoute = "/api/viewer/";
            getCurrentDir = function () {
                return document.getElementById("current-dir").innerText.slice(1);
            };
            addEventListenerToUploadFileForm = function () {
                var uploadForm = document.getElementById("upload-form");
                uploadForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var errFunc = function (resp) {
                        display_3.displayError(resp.error.message);
                    };
                    var okFunc = function () {
                        location.reload(true);
                    };
                    ajax_3.submitAjaxFormData(apiRoute + "upload/" + getCurrentDir(), uploadForm, errFunc, okFunc);
                });
            };
            addEventListenerToCreateFolderForm = function () {
                var createFolderForm = document.getElementById("create-folder-form");
                createFolderForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var folderName = createFolderForm.folder_name;
                    var data = {
                        path: makePath(getCurrentDir(), folderName.value)
                    };
                    ajaxHelper(apiRoute + "create", data);
                });
            };
            addEventListenerToDeleteFileFolderForm = function () {
                var deleteFileFolderForm = document.getElementById("delete-file-folder-form");
                deleteFileFolderForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var fileName = deleteFileFolderForm.file_name;
                    var data = {
                        path: makePath(getCurrentDir(), fileName.value)
                    };
                    ajaxHelper(apiRoute + "delete", data);
                });
            };
            addEventListenerToDeleteAllForm = function () {
                var deleteAllForm = document.getElementById("delete-all-form");
                deleteAllForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var sendAjax = function (deletePath) {
                        var data = {
                            path: deletePath
                        };
                        ajaxHelper(apiRoute + "delete-all", data);
                    };
                    if (getCurrentDir() === "") {
                        sendAjax("/");
                    }
                    else {
                        sendAjax(getCurrentDir());
                    }
                });
            };
            ajaxHelper = function (url, data) {
                var errFunc = function (resp) {
                    display_3.displayError(resp.error.message);
                };
                var okFunc = function () {
                    location.reload(true);
                };
                ajax_3.submitAjaxJSON(url, data, errFunc, okFunc);
            };
            makePath = function (currentDir, fileName) {
                var index = currentDir === "";
                if (index) {
                    return fileName;
                }
                return currentDir + "/" + fileName;
            };
            exports_5("addEventListenersViewerPage", addEventListenersViewerPage = function () {
                addEventListenerToUploadFileForm();
                addEventListenerToCreateFolderForm();
                addEventListenerToDeleteFileFolderForm();
                addEventListenerToDeleteAllForm();
            });
        }
    };
});
System.register("app/lib/userPage", ["app/bin/display", "app/bin/ajax"], function (exports_6, context_6) {
    "use strict";
    var display_4, ajax_4, userApiRoute, addEventListenerToChangeNameForm, addEventListenerToChangePasswordForm, addEventListenerToDeleteAccountForm, addEventListenersUserPage;
    var __moduleName = context_6 && context_6.id;
    return {
        setters: [
            function (display_4_1) {
                display_4 = display_4_1;
            },
            function (ajax_4_1) {
                ajax_4 = ajax_4_1;
            }
        ],
        execute: function () {
            userApiRoute = "/api/user/";
            addEventListenerToChangeNameForm = function () {
                var changeNameForm = document.getElementById("change-name-form");
                changeNameForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var firstname = changeNameForm.first_name;
                    var lastname = changeNameForm.last_name;
                    var data = {
                        first_name: firstname.value,
                        last_name: lastname.value
                    };
                    var errFunc = function (resp) {
                        display_4.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        display_4.displaySuccess(resp.data.content);
                        var firstname = document.getElementById("firstname");
                        var lastname = document.getElementById("lastname");
                        firstname.innerText = data.first_name;
                        lastname.innerText = data.last_name;
                    };
                    ajax_4.submitAjaxJSON(userApiRoute + "change-name", data, errFunc, okFunc);
                });
            };
            addEventListenerToChangePasswordForm = function () {
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
                        display_4.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        display_4.displaySuccess(resp.data.content);
                        changePasswordForm.reset();
                    };
                    ajax_4.submitAjaxJSON(userApiRoute + "change-password", data, errFunc, okFunc);
                });
            };
            addEventListenerToDeleteAccountForm = function () {
                var deleteAccountForm = document.getElementById("delete-account-form");
                deleteAccountForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var password = deleteAccountForm.password;
                    var data = {
                        password: password.value
                    };
                    var errFunc = function (resp) {
                        display_4.displayError(resp.error.message);
                        deleteAccountForm.reset();
                    };
                    var okFunc = function () {
                        window.location.href = "/login";
                    };
                    ajax_4.submitAjaxJSON(userApiRoute + "delete", data, errFunc, okFunc);
                });
            };
            exports_6("addEventListenersUserPage", addEventListenersUserPage = function () {
                addEventListenerToChangeNameForm();
                addEventListenerToChangePasswordForm();
                addEventListenerToDeleteAccountForm();
            });
        }
    };
});
System.register("app/lib/adminPage", ["app/bin/display", "app/bin/ajax"], function (exports_7, context_7) {
    "use strict";
    var display_5, ajax_5, adminApiRoute, addEventListenerToChangeUsernameForm, addEventListenerToChangeDirectoryRootForm, addEventListenerToChangeAdminStatusForm, addEventListenerToCreateUserForm, addEventListenerToDeleteUserForm, addEventListenersAdminPage;
    var __moduleName = context_7 && context_7.id;
    return {
        setters: [
            function (display_5_1) {
                display_5 = display_5_1;
            },
            function (ajax_5_1) {
                ajax_5 = ajax_5_1;
            }
        ],
        execute: function () {
            adminApiRoute = "/api/admin/";
            addEventListenerToChangeUsernameForm = function () {
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
                        display_5.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        var username = document.getElementById("username");
                        if (data.current_username === username.innerText) {
                            username.innerText = data.new_username;
                        }
                        display_5.displaySuccess(resp.data.content);
                        changeUsernameForm.reset();
                    };
                    ajax_5.submitAjaxJSON(adminApiRoute + "change-username", data, errFunc, okFunc);
                });
            };
            addEventListenerToChangeDirectoryRootForm = function () {
                var changeDirForm = document.getElementById("change-dir-root-form");
                changeDirForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var dirRoot = changeDirForm.dir_root;
                    var data = {
                        dir_root: dirRoot.value
                    };
                    var errFunc = function (resp) {
                        display_5.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        display_5.displaySuccess(resp.data.content);
                        changeDirForm.reset();
                    };
                    ajax_5.submitAjaxJSON(adminApiRoute + "change-dir-root", data, errFunc, okFunc);
                });
            };
            addEventListenerToChangeAdminStatusForm = function () {
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
                        display_5.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        display_5.displaySuccess(resp.data.content);
                        changeAdminStatusForm.reset();
                    };
                    ajax_5.submitAjaxJSON(adminApiRoute + "change-admin-status", data, errFunc, okFunc);
                });
            };
            addEventListenerToCreateUserForm = function () {
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
                        display_5.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        display_5.displaySuccess(resp.data.content);
                        createUserForm.reset();
                    };
                    ajax_5.submitAjaxJSON(adminApiRoute + "create-user", data, errFunc, okFunc);
                });
            };
            addEventListenerToDeleteUserForm = function () {
                var deleteUserForm = document.getElementById("delete-user-form");
                deleteUserForm.addEventListener("submit", function (event) {
                    event.preventDefault();
                    var userID = deleteUserForm.user_id;
                    var data = {
                        user_id: parseInt(userID.value)
                    };
                    var errFunc = function (resp) {
                        display_5.displayError(resp.error.message);
                    };
                    var okFunc = function (resp) {
                        display_5.displaySuccess(resp.data.content);
                        deleteUserForm.reset();
                    };
                    ajax_5.submitAjaxJSON(adminApiRoute + "delete-user", data, errFunc, okFunc);
                });
            };
            exports_7("addEventListenersAdminPage", addEventListenersAdminPage = function () {
                addEventListenerToChangeUsernameForm();
                addEventListenerToChangeDirectoryRootForm();
                addEventListenerToChangeAdminStatusForm();
                addEventListenerToCreateUserForm();
                addEventListenerToDeleteUserForm();
            });
        }
    };
});
System.register("index", ["app/lib/loginPage", "app/lib/navbar", "app/lib/viewerPage", "app/lib/userPage", "app/lib/adminPage"], function (exports_8, context_8) {
    "use strict";
    var loginPage_1, navbar_1, viewerPage_1, userPage_1, adminPage_1;
    var __moduleName = context_8 && context_8.id;
    return {
        setters: [
            function (loginPage_1_1) {
                loginPage_1 = loginPage_1_1;
            },
            function (navbar_1_1) {
                navbar_1 = navbar_1_1;
            },
            function (viewerPage_1_1) {
                viewerPage_1 = viewerPage_1_1;
            },
            function (userPage_1_1) {
                userPage_1 = userPage_1_1;
            },
            function (adminPage_1_1) {
                adminPage_1 = adminPage_1_1;
            }
        ],
        execute: function () {
            (function () {
                var page = window.location.pathname;
                if (page === "/login") {
                    loginPage_1.addEventListenersLoginPage();
                    return;
                }
                navbar_1.addEventListenersNavbar();
                var isViewerPage = page.search("/viewer/") !== -1;
                if (isViewerPage) {
                    viewerPage_1.addEventListenersViewerPage();
                    return;
                }
                switch (page) {
                    case "/user":
                        userPage_1.addEventListenersUserPage();
                        break;
                    case "/admin":
                        adminPage_1.addEventListenersAdminPage();
                        break;
                }
            })();
        }
    };
});
