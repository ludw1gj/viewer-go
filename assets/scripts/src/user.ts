// changeNameInput contains the required data structure.
interface changeNameInput {
    first_name: string;
    last_name: string
}

// changePasswordInput contains the required data structure.
interface changePasswordInput {
    old_password: string;
    new_password: string
}

// deleteAccountInput contains the required data structure.
interface deleteAccountInput {
    password: string
}

// addEventListenersUserForms function should be run at initialisation of user page.
function addEventListenersUserForms(): void {
    const userApiRoute = "/api/user/";

    // handle change name form logic
    let changeNameForm = document.getElementById("change-name-form") as HTMLFormElement;
    changeNameForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const firstName: HTMLInputElement = changeNameForm.first_name;
        const lastName: HTMLInputElement = changeNameForm.last_name;

        const data: changeNameInput = {
            first_name: firstName.value,
            last_name: lastName.value
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccessNotification(resp.data.content);
            location.reload(true);
        };
        submitAjaxJSON(userApiRoute + "change-name", data, errFunc, okFunc);
    });

    // handle change password form logic
    let changePasswordForm = document.getElementById("change-password-form") as HTMLFormElement;
    changePasswordForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const oldPassword: HTMLInputElement = changePasswordForm.old_password;
        const newPassword: HTMLInputElement = changePasswordForm.new_password;

        const data: changePasswordInput = {
            old_password: oldPassword.value,
            new_password: newPassword.value
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccessNotification(resp.data.content);
            changePasswordForm.reset();
        };
        submitAjaxJSON(userApiRoute + "change-password", data, errFunc, okFunc);
    });

    // handle delete user form logic
    let deleteAccountForm = document.getElementById("delete-account-form") as HTMLFormElement;
    deleteAccountForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const password: HTMLInputElement = deleteAccountForm.password;

        const data: deleteAccountInput = {
            password: password.value
        };

        const errFunc = function (resp: JSONErrorResponse) {
            displayErrorNotification(resp.error.message);
            deleteAccountForm.reset();
        };

        const okFunc = function () {
            window.location.href = "/login";
        };
        submitAjaxJSON(userApiRoute + "delete", data, errFunc, okFunc);
    });
}
