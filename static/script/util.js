"use strict";
/**
 * This function serialises a Form Element Object into a general Javascript Object
 * @param {Object} form
 * A DOM Form Element Object
 * @returns {Object}
 * A general Javascript Object
 */
function serializeFormObj(form) {
    var elems = form.elements;
    var obj = {};

    for (var i = 0; i < elems.length; i += 1) {
        var element = elems[i];
        var type = element.type;
        var name = element.name;
        var value = element.value;

        switch (type) {
            case "text":
                obj[name] = value;
                break;
            case "password":
                obj[name] = value;
                break;
            case "checkbox":
                obj[name] = element.checked;
                break;
            case "number":
                obj[name] = parseInt(value);
                break;
            default:
                break;
        }
    }
    return obj;
}

/**
 * This function submits an AJAX POST request
 * @param url {String}
 * The URL to send the request to
 * @param data {Object}
 * The data to send
 * @param errFunc {Function}
 * A function to execute if server returned an error
 * @param okFunc {Function}
 * A function to execute if request was successful
 */
function submitAjaxJson(url, data, errFunc, okFunc) {
    var httpStatusUnauthorised = 401,
        httpStatusInternalServerError = 500,
        httpStatusOK = 200;

    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === httpStatusUnauthorised || this.status === httpStatusInternalServerError) {
            errFunc(resp);
        } else if (this.status === httpStatusOK) {
            okFunc(resp);
        }
    };
    xhr.send(JSON.stringify(data));
}

/**
 * The function uploads files via AJAX
 * @param url {String}
 * The URL to send the request to
 * @param uploadForm {Object}
 * A form object to be sent as AJAX, enctype="multipart/form-data"
 * @param errFunc {Function}
 * A function to execute if server returned an error
 * @param okFunc {Function}
 * A function to execute if request was successful
 */
function submitAjaxFormData(url, uploadForm, errFunc, okFunc) {
    var formData = new FormData(uploadForm);
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.onload = function () {
        var resp = JSON.parse(xhr.responseText);
        if (this.status === 401 || this.status === 500) {
            errFunc(resp);
        } else if (this.status === 200) {
            okFunc(resp);
        }
    };
    xhr.send(formData);
}
