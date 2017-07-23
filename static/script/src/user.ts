// changePasswordInput contains the required data structure.
interface changePasswordInput {
    old_password: string;
    new_password: string
}

// deleteAccountInput contains the required data structure.
interface deleteAccountInput {
    password: string
}

// initUserPage function should be run at initialisation of user page.
function initUserPage(): void {
    const userApiRoute = "/api/user/";

    // handle change password form logic
    let userChangePasswordForm = document.getElementById("change-password-form") as HTMLFormElement;
    userChangePasswordForm.addEventListener('submit', (event: Event) => {
        event.preventDefault();

        const url: string = userApiRoute + "change-password";
        const oldPw = userChangePasswordForm.elements.item(0) as HTMLInputElement;
        const newPw = userChangePasswordForm.elements.item(1) as HTMLInputElement;
        const data: changePasswordInput = {
            old_password: oldPw.value,
            new_password: newPw.value
        };
        const errFunc = (resp: JsonErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };
        const okFunc = (resp: JsonDataResponse) => {
            displaySuccessNotification(resp.data.content);
            oldPw.value = "";
            newPw.value = "";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });

    // handle delete user form logic
    let userDeleteAccountForm = document.getElementById("delete-account-form") as HTMLFormElement;
    userDeleteAccountForm.addEventListener('submit', (event: Event) => {
        event.preventDefault();

        const url: string = userApiRoute + "delete";
        const pw = userDeleteAccountForm.elements.item(0) as HTMLInputElement;
        const data: deleteAccountInput = {
            password: pw.value
        };
        const errFunc = function (resp: JsonErrorResponse) {
            displayErrorNotification(resp.error.message);
        };
        const okFunc = function () {
            window.location.href = "/login";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });
}
