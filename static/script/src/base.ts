// addEventListenersBaseNav function should be run at initialisation of base page.
function addEventListenersBaseNav(): void {
    // extend and collapse navigation menu for mobile
    let mobileMenuButton = document.getElementById("mobile-menu-button") as HTMLElement;
    mobileMenuButton.addEventListener("click", () => {
        let mobileMenu = document.getElementById("mobile-menu") as HTMLElement;
        if (mobileMenuButton.classList.contains("is-active") || mobileMenuButton.classList.contains("is-active")) {
            mobileMenu.classList.remove("is-active");
            mobileMenuButton.classList.remove("is-active");
        } else {
            mobileMenuButton.classList.add("is-active");
            mobileMenu.classList.add("is-active");
        }
    });

    // handle logout user
    let logoutButton = document.getElementById("logout-button") as HTMLElement;
    logoutButton.addEventListener('click', () => {
        const errFunc = (resp: JsonErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };
        const okFunc = () => {
            window.location.href = "/login";
        };
        submitAjaxJson("/api/user/logout", undefined, errFunc, okFunc);
    });
}

// displayErrorNotification displays error notification.
function displayErrorNotification(msg: string): void {
    let notification = document.getElementById("notification") as HTMLElement;
    notification.classList.remove("is-success", "hidden");
    notification.classList.add("is-danger");
    notification.innerText = msg;
}

// displaySuccessNotification displays success notification.
function displaySuccessNotification(msg: string): void {
    let notification = document.getElementById("notification") as HTMLElement;
    notification.classList.remove("is-danger", "hidden");
    notification.classList.add("is-success");
    notification.innerText = msg;
}

// load authorized page's script.
function loadAuthorizedPages(): void {
    const page = window.location.pathname;
    if (page !== "/login") {
        addEventListenersBaseNav();
    }
    if (page.search("/viewer/") !== -1) {
        // user is on the viewer page.
        addEventListenersViewerForms();
        return
    }
    switch (page) {
        case "/user":
            addEventListenersUserForms();
            break;
        case "/admin":
            addEventListenersAdminForms();
            break;
    }
}
