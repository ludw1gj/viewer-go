import {displayError} from "../Handler/NotificationHandler";
import {ajaxSubmitFormData, ajaxSubmitJSON, JSONErrorResponse} from "../Handler/AjaxHandler";

interface pathInput {
    path: string;
}

const apiRoute = "/api/viewer/";
const currentDir: string = (document.getElementById("current-dir")  as HTMLElement).innerText.slice(1);

function addEventListenerToUploadFileForm(): void {
    const uploadForm = document.getElementById("upload-form") as HTMLFormElement;

    uploadForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const errFunc = (resp: JSONErrorResponse) => {
            displayError(resp.error.message);
        };

        const okFunc = () => {
            location.reload(true);
        };

        ajaxSubmitFormData(apiRoute + "upload/" + currentDir, uploadForm, errFunc, okFunc);
    });
}

function addEventListenerToCreateFolderForm(): void {
    const createFolderForm = document.getElementById("create-folder-form") as HTMLFormElement;

    createFolderForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const folderName: HTMLInputElement = createFolderForm.folder_name;

        const data: pathInput = {
            path: makePath(currentDir, folderName.value)
        };
        ajaxHelper(apiRoute + "create", data);
    });
}

function addEventListenerToDeleteFileFolderForm(): void {
    const deleteFileFolderForm = document.getElementById("delete-file-folder-form") as HTMLFormElement;

    deleteFileFolderForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const fileName: HTMLInputElement = deleteFileFolderForm.file_name;

        const data: pathInput = {
            path: makePath(currentDir, fileName.value)
        };
        ajaxHelper(apiRoute + "delete", data);
    });
}

function addEventListenerToDeleteAllForm(): void {
    const deleteAllForm = document.getElementById("delete-all-form") as HTMLFormElement;

    deleteAllForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const sendAjax = function (deletePath: string) {
            const data: pathInput = {
                path: deletePath
            };
            ajaxHelper(apiRoute + "delete-all", data);
        };

        if (currentDir === "") {
            sendAjax("/");
        } else {
            sendAjax(currentDir);
        }
    });
}

function ajaxHelper(url: string, data: object): void {
    const errFunc = (resp: JSONErrorResponse) => {
        displayError(resp.error.message);
    };
    const okFunc = () => {
        location.reload(true);
    };
    ajaxSubmitJSON(url, data, errFunc, okFunc)
}

function makePath(currentDir: string, fileName: string): string {
    const index: boolean = (currentDir === "");
    if (index) {
        return fileName;
    }
    return currentDir + "/" + fileName;
}

export function initiateViewerPage() {
    addEventListenerToUploadFileForm();
    addEventListenerToCreateFolderForm();
    addEventListenerToDeleteFileFolderForm();
    addEventListenerToDeleteAllForm();
}
