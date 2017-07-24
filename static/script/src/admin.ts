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

    // handle create user form logic
    let adminCreateUserForm = document.getElementById("create-user-form") as HTMLFormElement;
    adminCreateUserForm.addEventListener('submit', (event: Event) => {
        event.preventDefault();

        const url = adminApiRoute + "create-user";
        const data: createUserInput = {
            username: adminCreateUserForm.username.value as string,
            password: adminCreateUserForm.password.value as string,
            first_name: adminCreateUserForm.first_name.value as string,
            last_name: adminCreateUserForm.last_name.value as string,
            directory_root: adminCreateUserForm.directory_root.value as string,
            is_admin: adminCreateUserForm.is_admin.checked as boolean
        };
        const errFunc = (resp: JsonErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };
        const okFunc = (resp: JsonDataResponse) => {
            displaySuccessNotification(resp.data.content);
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });

    // handle delete user form logic
    let adminDeleteUserForm = document.getElementById("delete-user-form") as HTMLFormElement;
    adminDeleteUserForm.addEventListener('submit', (event: Event) => {
        event.preventDefault();

        const url = adminApiRoute + "delete-user";
        const data: deleteUserInput = {
            user_id: parseInt(adminDeleteUserForm.user_id.value)
        };
        const errFunc = (resp: JsonErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };
        const okFunc = (resp: JsonDataResponse) => {
            displaySuccessNotification(resp.data.content);
            adminDeleteUserForm.user_id.value = "";
        };
        submitAjaxJson(url, data, errFunc, okFunc);
    });
}
