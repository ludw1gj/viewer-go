class NotificationHandler {

    public static displayError(msg: string): void {
        let notification = document.getElementById("notification") as HTMLElement;
        notification.classList.remove("is-success", "hidden");
        notification.classList.add("is-danger");
        notification.innerText = msg;
    }

    public static displaySuccess(msg: string): void {
        let notification = document.getElementById("notification") as HTMLElement;
        notification.classList.remove("is-danger", "hidden");
        notification.classList.add("is-success");
        notification.innerText = msg;
    }
}

export {NotificationHandler}
