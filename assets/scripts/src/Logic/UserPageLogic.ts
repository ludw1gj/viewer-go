import {displayError, displaySuccess} from "../Handler/NotificationHandler";
import {ajaxSubmitJSON, JSONDataResponse, JSONErrorResponse} from "../Handler/AjaxHandler";

interface changeNameInput {
    first_name: string;
    last_name: string
}

interface changePasswordInput {
    old_password: string;
    new_password: string
}

interface deleteAccountInput {
    password: string
}

const userApiRoute = "/api/user/";

function addEventListenerToChangeNameForm(): void {
    const changeNameForm = document.getElementById("change-name-form") as HTMLFormElement;

    changeNameForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const firstName: HTMLInputElement = changeNameForm.first_name;
        const lastName: HTMLInputElement = changeNameForm.last_name;

        const data: changeNameInput = {
            first_name: firstName.value,
            last_name: lastName.value
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayError(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccess(resp.data.content);
            location.reload(true);
        };
        ajaxSubmitJSON(userApiRoute + "change-name", data, errFunc, okFunc);
    });
}

function addEventListenerToChangePasswordForm(): void {
    const changePasswordForm = document.getElementById("change-password-form") as HTMLFormElement;

    changePasswordForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const oldPassword: HTMLInputElement = changePasswordForm.old_password;
        const newPassword: HTMLInputElement = changePasswordForm.new_password;

        const data: changePasswordInput = {
            old_password: oldPassword.value,
            new_password: newPassword.value
        };

        const errFunc = (resp: JSONErrorResponse) => {
            displayError(resp.error.message);
        };

        const okFunc = (resp: JSONDataResponse) => {
            displaySuccess(resp.data.content);
            changePasswordForm.reset();
        };
        ajaxSubmitJSON(userApiRoute + "change-password", data, errFunc, okFunc);
    });
}

function addEventListenerToDeleteAccountForm(): void {
    const deleteAccountForm = document.getElementById("delete-account-form") as HTMLFormElement;

    deleteAccountForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const password: HTMLInputElement = deleteAccountForm.password;

        const data: deleteAccountInput = {
            password: password.value
        };

        const errFunc = function (resp: JSONErrorResponse) {
            displayError(resp.error.message);
            deleteAccountForm.reset();
        };

        const okFunc = function () {
            window.location.href = "/login";
        };
        ajaxSubmitJSON(userApiRoute + "delete", data, errFunc, okFunc);
    });
}

export function initiateUserPage() {
    addEventListenerToChangeNameForm();
    addEventListenerToChangePasswordForm();
    addEventListenerToDeleteAccountForm()
}