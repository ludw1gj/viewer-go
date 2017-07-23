// loginInput contains the required data structure.
interface loginInput {
    username: string;
    password: string;
}

// initLoginPage function should be run at initialisation of login page.
function initLoginPage(): void {
    // handle login user form logic
    let loginForm = document.getElementById("login-form") as HTMLFormElement;
    loginForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const url = "/api/user/login";
        const data: loginInput = {
            username: loginForm.username.value,
            password: loginForm.password.value
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
        submitAjaxJson(url, data, errFunc, okFunc);
    });
}

// loadLoginPage loads login page script if at login page.
function loadLoginPage(): void {
    if (window.location.pathname === "/login") {
        initLoginPage();
    }
}

// run init.
loadLoginPage();