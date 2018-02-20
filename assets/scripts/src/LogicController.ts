import {initiateLoginPage} from "./Logic/LoginPageLogic";
import {initiateNavbar} from "./Logic/NavbarLogic";
import {initiateViewerPage} from "./Logic/ViewerPageLogic";
import {initiateUserPage} from "./Logic/UserPageLogic";
import {initiateAdminPage} from "./Logic/AdminPageLogic";

function run(): void {
    const page = window.location.pathname;

    if (page === "/login") {
        initiateLoginPage();
        return;
    }
    initiateNavbar();

    const isViewerPage = page.search("/viewer/") !== -1;
    if (isViewerPage) {
        initiateViewerPage();
        return;
    }

    switch (page) {
        case "/user":
            initiateUserPage();
            break;
        case "/admin":
            initiateAdminPage();
            break;
    }
}

run();
