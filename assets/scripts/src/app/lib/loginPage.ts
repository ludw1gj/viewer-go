import { submitAjaxJSON, JSONErrorResponse } from "../bin/ajax";

interface loginInput {
  username: string;
  password: string;
}

const addEventListenerLoginForm = () => {
  const loginForm = document.getElementById("login-form") as HTMLFormElement;

  loginForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const username: HTMLInputElement = loginForm.username;
    const password: HTMLInputElement = loginForm.password;

    const data: loginInput = {
      username: username.value,
      password: password.value
    };

    const errFunc = (resp: JSONErrorResponse) => {
      const notification = document.getElementById(
        "login-error-notification"
      ) as HTMLFormElement;
      notification.classList.remove("hidden");
      notification.classList.add("is-danger");
      notification.innerText = resp.error.message;
    };

    const okFunc = () => {
      window.location.href = "/viewer/";
    };
    submitAjaxJSON("/api/user/login", data, errFunc, okFunc);
  });
};

export const addEventListenersLoginPage = () => {
  addEventListenerLoginForm();
};
