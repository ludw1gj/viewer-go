import { addEventListenersLoginPage } from "./app/lib/loginPage";
import { addEventListenersNavbar } from "./app/lib/navbar";
import { addEventListenersViewerPage } from "./app/lib/viewerPage";
import { addEventListenersUserPage } from "./app/lib/userPage";
import { addEventListenersAdminPage } from "./app/lib/adminPage";

(() => {
  const page = window.location.pathname;

  if (page === "/login") {
    addEventListenersLoginPage();
    return;
  }
  addEventListenersNavbar();

  const isViewerPage = page.search("/viewer/") !== -1;
  if (isViewerPage) {
    addEventListenersViewerPage();
    return;
  }

  switch (page) {
    case "/user":
      addEventListenersUserPage();
      break;
    case "/admin":
      addEventListenersAdminPage();
      break;
  }
})();
