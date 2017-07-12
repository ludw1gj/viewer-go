"use strict";
/**
 * This function displays error notification.
 * @param {String} msg
 * A message to display in the notification
 */
function displayErrorNotification(msg) {
    notification.classList.remove("is-success", "hidden");
    notification.classList.add("is-danger");
    notification.innerText = msg;
}

/**
 * This function displays success notification.
 * @param {String} msg
 * A message to display in the notification
 */
function displaySuccessNotification(msg) {
    notification.classList.remove("is-danger", "hidden");
    notification.classList.add("is-success");
    notification.innerText = msg;
}

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
