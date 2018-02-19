import {LoginPageLogic} from "./Logic/LoginPageLogic";
import {NavbarLogic} from "./Logic/NavbarLogic";
import {ViewerPageLogic} from "./Logic/ViewerPageLogic";
import {UserPageLogic} from "./Logic/UserPageLogic";
import {AdminPageLogic} from "./Logic/AdminPageLogic";

class LogicController {

    public static run(): void {
        const page = window.location.pathname;

        if (page === "/login") {
            new LoginPageLogic();
            return;
        }

        new NavbarLogic();

        const isViewerPage = page.search("/viewer/") !== -1;
        if (isViewerPage) {
            new ViewerPageLogic();
            return;
        }

        switch (page) {
            case "/user":
                new UserPageLogic();
                break;
            case "/admin":
                new AdminPageLogic();
                break;
        }
    }

}

LogicController.run();
