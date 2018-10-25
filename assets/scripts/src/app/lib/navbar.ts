import { displayError } from "../bin/display";
import { submitAjaxJSON, JSONErrorResponse } from "../bin/ajax";

const addEventListenerToMobileMenuButton = () => {
  // extend and collapse navigation menu for mobile
  const mobileMenuButton = document.getElementById(
    "mobile-menu-button"
  ) as HTMLElement;

  mobileMenuButton.addEventListener("click", () => {
    const mobileMenu = document.getElementById("mobile-menu") as HTMLElement;

    if (
      mobileMenu.classList.contains("is-active") ||
      mobileMenuButton.classList.contains("is-active")
    ) {
      mobileMenu.classList.remove("is-active");
      mobileMenuButton.classList.remove("is-active");
    } else {
      mobileMenuButton.classList.add("is-active");
      mobileMenu.classList.add("is-active");
    }
  });
};

const addEventListenerToLogoutButton = () => {
  const logoutButton = document.getElementById("logout-button") as HTMLElement;

  logoutButton.addEventListener("click", () => {
    const errFunc = (resp: JSONErrorResponse) => {
      displayError(resp.error.message);
    };
    const okFunc = () => {
      window.location.href = "/login";
    };
    submitAjaxJSON("/api/user/logout", undefined, errFunc, okFunc);
  });
};

export const addEventListenersNavbar = () => {
  addEventListenerToMobileMenuButton();
  addEventListenerToLogoutButton();
};
