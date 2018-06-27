import { displayError, displaySuccess } from "../bin/display";
import {
  submitAjaxJSON,
  JSONDataResponse,
  JSONErrorResponse
} from "../bin/ajax";

interface changeNameInput {
  first_name: string;
  last_name: string;
}

interface changePasswordInput {
  old_password: string;
  new_password: string;
}

interface deleteAccountInput {
  password: string;
}

const addEventListenerToChangeNameForm = (apiRoute: string) => {
  const changeNameForm = document.getElementById(
    "change-name-form"
  ) as HTMLFormElement;

  changeNameForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const firstname: HTMLInputElement = changeNameForm.first_name;
    const lastname: HTMLInputElement = changeNameForm.last_name;

    const data: changeNameInput = {
      first_name: firstname.value,
      last_name: lastname.value
    };

    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };

    const okFunc = (resp: JSONDataResponse) => {
      displaySuccess(resp.data.content);

      const firstname = document.getElementById("firstname") as HTMLSpanElement;
      const lastname = document.getElementById("lastname") as HTMLSpanElement;

      firstname.innerText = data.first_name;
      lastname.innerText = data.last_name;
    };
    submitAjaxJSON(apiRoute + "change-name", data, errFunc, okFunc);
  });
};

const addEventListenerToChangePasswordForm = (apiRoute: string) => {
  const changePasswordForm = document.getElementById(
    "change-password-form"
  ) as HTMLFormElement;

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
    submitAjaxJSON(apiRoute + "change-password", data, errFunc, okFunc);
  });
};

const addEventListenerToDeleteAccountForm = (apiRoute: string) => {
  const deleteAccountForm = document.getElementById(
    "delete-account-form"
  ) as HTMLFormElement;

  deleteAccountForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const password: HTMLInputElement = deleteAccountForm.password;

    const data: deleteAccountInput = {
      password: password.value
    };

    const errFunc = function(resp: JSONErrorResponse) {
      displayError(resp.error.message);
      deleteAccountForm.reset();
    };

    const okFunc = function() {
      window.location.href = "/login";
    };
    submitAjaxJSON(apiRoute + "delete", data, errFunc, okFunc);
  });
};

export const addEventListenersUserPage = () => {
  const userApiRoute = "/api/user/";

  addEventListenerToChangeNameForm(userApiRoute);
  addEventListenerToChangePasswordForm(userApiRoute);
  addEventListenerToDeleteAccountForm(userApiRoute);
};
