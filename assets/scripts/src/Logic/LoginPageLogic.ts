import {AjaxHandler, JSONErrorResponse} from "../Handler/AjaxHandler";

interface loginInput {
    username: string;
    password: string;
}

class LoginPageLogic {

    constructor() {
        this.addEventListenerLoginForm();
    }

    private addEventListenerLoginForm(): void {
        let loginForm = document.getElementById("login-form") as HTMLFormElement;

        loginForm.addEventListener("submit", (event: Event) => {
            event.preventDefault();

            const username: HTMLInputElement = loginForm.username;
            const password: HTMLInputElement = loginForm.password;

            const data: loginInput = {
                username: username.value,
                password: password.value
            };

            const errFunc = (resp: JSONErrorResponse) => {
                let notification = document.getElementById("login-error-notification") as HTMLFormElement;
                notification.classList.remove("hidden");
                notification.classList.add("is-danger");
                notification.innerText = resp.error.message;
            };

            const okFunc = () => {
                window.location.href = "/viewer/";
            };
            AjaxHandler.submitJSON("/api/user/login", data, errFunc, okFunc);
        });
    }

}

export {LoginPageLogic}