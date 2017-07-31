// loginInput contains the required data structure.
interface loginInput {
    username: string;
    password: string;
}

// addEventListenersLoginForm function should be run at initialisation of login page.
function addEventListenerLoginForm(): void {
    // handle login user form logic
    let loginForm = document.getElementById("login-form") as HTMLFormElement;
    loginForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const username: HTMLInputElement = loginForm.username;
        const password: HTMLInputElement = loginForm.password;

        const data: loginInput = {
            username: username.value,
            password: password.value
        };

        const errFunc = (resp: JsonErrorResponse) => {
            let notification = document.getElementById("login-error-notification") as HTMLFormElement;
            notification.classList.remove("hidden");
            notification.classList.add("is-danger");
            notification.innerText = resp.error.message;
        };

        const okFunc = () => {
            window.location.href = "/viewer/";
        };
        submitAjaxJson("/api/user/login", data, errFunc, okFunc);
    });
}

// loadLoginPage loads login page script if at login page.
function loadLoginPageScript(): void {
    if (window.location.pathname === "/login") {
        addEventListenerLoginForm();
    }
}
