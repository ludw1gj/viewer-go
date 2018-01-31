// changeDirRootInput contains the required data structure.
interface changeUsernameInput {
    current_username: string;
    new_username: string
}

// changeDirRootInput contains the required data structure.
interface changeDirRootInput {
    dir_root: string;
}

// changeAdminStatusInput contains the required data structure.
interface changeAdminStatusInput {
    user_id: number;
    is_admin: boolean;
}

// createUserInput contains the required data structure.
interface createUserInput {
    username: string;
    password: string;
    first_name: string;
    last_name: string;
    directory_root: string;
    is_admin: boolean;
}

// deleteUserInput contains the required data structure.
interface deleteUserInput {
    user_id: number;
}

// addEventListenersAdminForms function should be run at initialisation of admin page.
function addEventListenersAdminForms(): void {
    const adminApiRoute = "/api/admin/";

    // handle change directory root form logic
    let changeUsernameForm = document.getElementById("change-username-form") as HTMLFormElement;
    changeUsernameForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const currentUsername: HTMLInputElement = changeUsernameForm.current_username;
        const newUsername: HTMLInputElement = changeUsernameForm.new_username;

        const data: changeUsernameInput = {
            current_username: currentUsername.value,
            new_username: newUsername.value
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            let username = document.getElementById("username") as HTMLSpanElement;
            if (data.current_username === username.innerText) {
                location.reload(true);
                return;
            }
            displaySuccessNotification(resp.data.content);
        };
        submitAjaxJSON(adminApiRoute + "change-username", data, errFunc, okFunc)
    });

    // handle change directory root form logic
    let changeDirForm = document.getElementById("change-dir-root-form") as HTMLFormElement;
    changeDirForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const dirRoot: HTMLInputElement = changeDirForm.dir_root;

        const data: changeDirRootInput = {
            dir_root: dirRoot.value
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccessNotification(resp.data.content);
            changeDirForm.reset();
        };
        submitAjaxJSON(adminApiRoute + "change-dir-root", data, errFunc, okFunc)
    });

    // handle change admin status form logic
    let changeAdminStatusForm = document.getElementById("change-admin-status-form") as HTMLFormElement;
    changeAdminStatusForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const userID: HTMLInputElement = changeAdminStatusForm.user_id;
        const isAdmin: HTMLInputElement = changeAdminStatusForm.is_admin;

        const data: changeAdminStatusInput = {
            user_id: parseInt(userID.value),
            is_admin: isAdmin.checked
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccessNotification(resp.data.content);
            changeDirForm.reset();
        };
        submitAjaxJSON(adminApiRoute + "change-admin-status", data, errFunc, okFunc)
    });

    // handle create user form logic
    let createUserForm = document.getElementById("create-user-form") as HTMLFormElement;
    createUserForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const username: HTMLInputElement = createUserForm.username;
        const password: HTMLInputElement = createUserForm.password;
        const firstName: HTMLInputElement = createUserForm.first_name;
        const lastName: HTMLInputElement = createUserForm.last_name;
        const DirRoot: HTMLInputElement = createUserForm.directory_root;
        const isAdmin: HTMLInputElement = createUserForm.is_admin;

        const data: createUserInput = {
            username: username.value,
            password: password.value,
            first_name: firstName.value,
            last_name: lastName.value,
            directory_root: DirRoot.value,
            is_admin: isAdmin.checked
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccessNotification(resp.data.content);
            createUserForm.reset();
        };
        submitAjaxJSON(adminApiRoute + "create-user", data, errFunc, okFunc);
    });

    // handle delete user form logic
    let deleteUserForm = document.getElementById("delete-user-form") as HTMLFormElement;
    deleteUserForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const userID: HTMLInputElement = deleteUserForm.user_id;

        const data: deleteUserInput = {
            user_id: parseInt(userID.value)
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccessNotification(resp.data.content);
            createUserForm.reset();
        };
        submitAjaxJSON(adminApiRoute + "delete-user", data, errFunc, okFunc);
    });
}
