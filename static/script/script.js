
document.getElementById("delete-all").value = window.location.pathname;

/**
 * This function submits an ajax request of content-type json to the 'api' route
 * @param {Object} data
 * A Form Element which has been serialized by serializeFormObj function
 */
function submitAjax(data) {
    var url = '/api' + window.location.pathname;
    var request = new XMLHttpRequest();
    request.open('POST', url, true);

    request.onload = function () {
        if (this.status >= 200) {
            var resp = JSON.parse(this.response);
            if (resp.error) {
                displayCalcCard(resp.error);
            } else {
                displayCalcCard(resp.content);
            }
        } else {
            displayCalcCard('The server has encountered a problem.');
        }
    };

    request.onerror = function () {
        console.log('There was a connection issue. Check your internet connection or the sever might be down.')
    };

    request.setRequestHeader('Content-Type', 'application/json');
    request.send(JSON.stringify(data));
}
