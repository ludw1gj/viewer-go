import { displayError } from "../bin/display";
import {
  submitAjaxFormData,
  submitAjaxJSON,
  JSONErrorResponse
} from "../bin/ajax";

interface pathInput {
  path: string;
}

const apiRoute = "/api/viewer/";

const getCurrentDir = (): string => {
  return (document.getElementById(
    "current-dir"
  ) as HTMLElement).innerText.slice(1);
};

const addEventListenerToUploadFileForm = () => {
  const uploadForm = document.getElementById("upload-form") as HTMLFormElement;

  uploadForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };

    const okFunc = () => {
      location.reload(true);
    };

    submitAjaxFormData(
      apiRoute + "upload/" + getCurrentDir(),
      uploadForm,
      errFunc,
      okFunc
    );
  });
};

const addEventListenerToCreateFolderForm = () => {
  const createFolderForm = document.getElementById(
    "create-folder-form"
  ) as HTMLFormElement;

  createFolderForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const folderName: HTMLInputElement = createFolderForm.folder_name;

    const data: pathInput = {
      path: makePath(getCurrentDir(), folderName.value)
    };
    ajaxHelper(apiRoute + "create", data);
  });
};

const addEventListenerToDeleteFileFolderForm = () => {
  const deleteFileFolderForm = document.getElementById(
    "delete-file-folder-form"
  ) as HTMLFormElement;

  deleteFileFolderForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const fileName: HTMLInputElement = deleteFileFolderForm.file_name;

    const data: pathInput = {
      path: makePath(getCurrentDir(), fileName.value)
    };
    ajaxHelper(apiRoute + "delete", data);
  });
};

const addEventListenerToDeleteAllForm = () => {
  const deleteAllForm = document.getElementById(
    "delete-all-form"
  ) as HTMLFormElement;

  deleteAllForm.addEventListener("submit", (event: Event) => {
    event.preventDefault();

    const sendAjax = function(deletePath: string) {
      const data: pathInput = {
        path: deletePath
      };
      ajaxHelper(apiRoute + "delete-all", data);
    };

    if (getCurrentDir() === "") {
      sendAjax("/");
    } else {
      sendAjax(getCurrentDir());
    }
  });
};

const ajaxHelper = (url: string, data: object) => {
  const errFunc = (resp: JSONErrorResponse) => {
    displayError(resp.error.message);
  };
  const okFunc = () => {
    location.reload(true);
  };
  submitAjaxJSON(url, data, errFunc, okFunc);
};

const makePath = (currentDir: string, fileName: string): string => {
  const index: boolean = currentDir === "";
  if (index) {
    return fileName;
  }
  return currentDir + "/" + fileName;
};

export const addEventListenersViewerPage = () => {
  addEventListenerToUploadFileForm();
  addEventListenerToCreateFolderForm();
  addEventListenerToDeleteFileFolderForm();
  addEventListenerToDeleteAllForm();
};
