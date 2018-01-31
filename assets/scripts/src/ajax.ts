// JSONErrorResponse is the expected json data response structure.
interface JSONErrorResponse {
    error: {
        message: string;
    };
}

// JSONDataResponse is the expected json error response structure.
interface JSONDataResponse {
    data: {
        content: string;
    };
}

// submitAjaxJSON submits an AJAX POST request.
function submitAjaxJSON(url: string,
                        data: object | undefined,
                        errFunc: (resp: JSONErrorResponse) => void,
                        okFunc: (resp: JSONDataResponse) => void) {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = () => {
        const DONE = 4;
        if (xhr.readyState === DONE) {
            const resp = JSON.parse(xhr.responseText) as JSONDataResponse | JSONErrorResponse;
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp as JSONErrorResponse);
            } else if ("data" in resp) {
                okFunc(resp as JSONDataResponse);
            } else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(JSON.stringify(data));
}

// submitAjaxFormData uploads files via AJAX.
function submitAjaxFormData(url: string,
                            uploadForm: HTMLFormElement,
                            errFunc: (resp: JSONErrorResponse) => void,
                            okFunc: (resp: JSONDataResponse) => void): void {
    const formData = new FormData(uploadForm);
    let xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.onreadystatechange = () => {
        const DONE = 4;
        if (xhr.readyState === DONE) {
            const resp = JSON.parse(xhr.responseText) as JSONDataResponse | JSONErrorResponse;
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp as JSONErrorResponse);
            } else if ("data" in resp) {
                okFunc(resp as JSONDataResponse);
            } else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(formData);
}
