// JsonErrorResponse is the expected json data response structure.
interface JsonErrorResponse {
    error: {
        message: string;
    };
}

// JsonDataResponse is the expected json error response structure.
interface JsonDataResponse {
    data: {
        content: string;
    };
}

// submitAjaxJson submits an AJAX POST request.
function submitAjaxJson(url: string, data: object | undefined, errFunc: (resp: JsonErrorResponse) => void,
                        okFunc: (resp: JsonDataResponse) => void) {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = () => {
        const DONE = 4;
        if (xhr.readyState === DONE) {
            const resp = JSON.parse(xhr.responseText) as JsonDataResponse | JsonErrorResponse;
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp as JsonErrorResponse);
                return;
            } else if ("data" in resp) {
                okFunc(resp as JsonDataResponse);
                return;
            } else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(JSON.stringify(data));
}

// submitAjaxFormData uploads files via AJAX.
function submitAjaxFormData(url: string, uploadForm: HTMLFormElement, errFunc: (resp: JsonErrorResponse) => void,
                            okFunc: (resp: JsonDataResponse) => void): void {
    const formData = new FormData(uploadForm);
    let xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.onreadystatechange = () => {
        const DONE = 4;
        if (xhr.readyState === DONE) {
            const resp = JSON.parse(xhr.responseText) as JsonDataResponse | JsonErrorResponse;
            if ("error" in resp || xhr.status === 401 || xhr.status === 500) {
                errFunc(resp as JsonErrorResponse);
                return;
            } else if ("data" in resp) {
                okFunc(resp as JsonDataResponse);
                return;
            } else {
                displayErrorNotification("There has been an error.");
            }
        }
    };
    xhr.send(formData);
}
