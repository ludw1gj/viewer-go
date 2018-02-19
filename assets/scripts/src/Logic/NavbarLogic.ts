import {NotificationHandler} from "../Handler/NotificationHandler";
import {AjaxHandler, JSONErrorResponse} from "../Handler/AjaxHandler";

class NavbarLogic {

    constructor() {
        this.addEventListenerToMobileMenuButton();
        this.addEventListenerToLogoutButton();
    }

    private addEventListenerToMobileMenuButton(): void {
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
    }

    private addEventListenerToLogoutButton(): void {
        let logoutButton = document.getElementById("logout-button") as HTMLElement;

        logoutButton.addEventListener('click', () => {
            const errFunc = (resp: JSONErrorResponse) => {
                NotificationHandler.displayError(resp.error.message);
            };
            const okFunc = () => {
                window.location.href = "/login";
            };
            AjaxHandler.submitJSON("/api/user/logout", undefined, errFunc, okFunc);
        });
    }

}

export {NavbarLogic}