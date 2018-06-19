import { displayError } from "./display";

export interface JSONErrorResponse {
  error: {
    message: string;
  };
}

export interface JSONDataResponse {
  data: {
    content: string;
  };
}

export const submitAjaxJSON = (
  url: string,
  data: object | undefined,
  errFunc: (resp: JSONErrorResponse) => void,
  okFunc: (resp: JSONDataResponse) => void
) => {
  const xhr = new XMLHttpRequest();

  xhr.open("POST", url, true);
  xhr.setRequestHeader("Content-Type", "application/json");

  xhr.onreadystatechange = () => {
    const DONE = 4;

    if (xhr.readyState === DONE) {
      const resp = JSON.parse(xhr.responseText) as
        | JSONDataResponse
        | JSONErrorResponse;
      if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
        errFunc(resp as JSONErrorResponse);
      } else if ("data" in resp) {
        okFunc(resp as JSONDataResponse);
      } else {
        displayError("There has been an error.");
      }
    }
  };
  xhr.send(JSON.stringify(data));
};

export const submitAjaxFormData = (
  url: string,
  uploadForm: HTMLFormElement,
  errFunc: (resp: JSONErrorResponse) => void,
  okFunc: (resp: JSONDataResponse) => void
) => {
  const formData = new FormData(uploadForm);

  const xhr = new XMLHttpRequest();
  xhr.open("POST", url, true);
  xhr.onreadystatechange = () => {
    const DONE = 4;

    if (xhr.readyState === DONE) {
      const resp = JSON.parse(xhr.responseText) as
        | JSONDataResponse
        | JSONErrorResponse;
      if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
        errFunc(resp as JSONErrorResponse);
      } else if ("data" in resp) {
        okFunc(resp as JSONDataResponse);
      } else {
        displayError("There has been an error.");
      }
    }
  };
  xhr.send(formData);
};
