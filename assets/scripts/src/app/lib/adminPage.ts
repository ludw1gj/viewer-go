import { displayError, displaySuccess } from "../bin/display";
import {
  submitAjaxJSON,
  JSONDataResponse,
  JSONErrorResponse
} from "../bin/ajax";

interface changeUsernameInput {
  current_username: string;
  new_username: string;
}

interface changeDirRootInput {
  dir_root: string;
}

interface changeAdminStatusInput {
  user_id: number;
  is_admin: boolean;
}

interface createUserInput {
  username: string;
  password: string;
  first_name: string;
  last_name: string;
  directory_root: string;
  is_admin: boolean;
}

interface deleteUserInput {
  user_id: number;
}

const adminApiRoute = "/api/admin/";

const addEventListenerToChangeUsernameForm = () => {
  const changeUsernameForm = document.getElementById(
    "change-username-form"
  ) as HTMLFormElement;

  changeUsernameForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const currentUsername: HTMLInputElement =
      changeUsernameForm.current_username;
    const newUsername: HTMLInputElement = changeUsernameForm.new_username;

    const data: changeUsernameInput = {
      current_username: currentUsername.value,
      new_username: newUsername.value
    };

    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };

    const okFunc = (resp: JSONDataResponse) => {
      const username = document.getElementById("username") as HTMLSpanElement;
      if (data.current_username === username.innerText) {
        location.reload(true);
        return;
      }
      displaySuccess(resp.data.content);
    };
    submitAjaxJSON(adminApiRoute + "change-username", data, errFunc, okFunc);
  });
};

const addEventListenerToChangeDirectoryRootForm = () => {
  const changeDirForm = document.getElementById(
    "change-dir-root-form"
  ) as HTMLFormElement;

  changeDirForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const dirRoot: HTMLInputElement = changeDirForm.dir_root;

    const data: changeDirRootInput = {
      dir_root: dirRoot.value
    };

    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };

    const okFunc = (resp: JSONDataResponse) => {
      displaySuccess(resp.data.content);
      changeDirForm.reset();
    };
    submitAjaxJSON(adminApiRoute + "change-dir-root", data, errFunc, okFunc);
  });
};

const addEventListenerToChangeAdminStatusForm = () => {
  const changeAdminStatusForm = document.getElementById(
    "change-admin-status-form"
  ) as HTMLFormElement;

  changeAdminStatusForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const userID: HTMLInputElement = changeAdminStatusForm.user_id;
    const isAdmin: HTMLInputElement = changeAdminStatusForm.is_admin;

    const data: changeAdminStatusInput = {
      user_id: parseInt(userID.value),
      is_admin: isAdmin.checked
    };

    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };

    const okFunc = (resp: JSONDataResponse) => {
      displaySuccess(resp.data.content);
      changeAdminStatusForm.reset();
    };
    submitAjaxJSON(
      adminApiRoute + "change-admin-status",
      data,
      errFunc,
      okFunc
    );
  });
};

const addEventListenerToCreateUserForm = () => {
  const createUserForm = document.getElementById(
    "create-user-form"
  ) as HTMLFormElement;

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
      displayError(resp.error.message);
    };

    const okFunc = (resp: JSONDataResponse) => {
      displaySuccess(resp.data.content);
      createUserForm.reset();
    };
    submitAjaxJSON(adminApiRoute + "create-user", data, errFunc, okFunc);
  });
};

const addEventListenerToDeleteUserForm = () => {
  const deleteUserForm = document.getElementById(
    "delete-user-form"
  ) as HTMLFormElement;

  deleteUserForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const userID: HTMLInputElement = deleteUserForm.user_id;

    const data: deleteUserInput = {
      user_id: parseInt(userID.value)
    };

    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };

    const okFunc = (resp: JSONDataResponse) => {
      displaySuccess(resp.data.content);
      deleteUserForm.reset();
    };
    submitAjaxJSON(adminApiRoute + "delete-user", data, errFunc, okFunc);
  });
};

export const addEventListenersAdminPage = () => {
  addEventListenerToChangeUsernameForm();
  addEventListenerToChangeDirectoryRootForm();
  addEventListenerToChangeAdminStatusForm();
  addEventListenerToCreateUserForm();
  addEventListenerToDeleteUserForm();
};
