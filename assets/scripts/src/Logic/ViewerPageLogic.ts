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
    let createFolderForm = document.getElementById("create-folder-form") as HTMLFormElement;

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
    let deleteFileFolderForm = document.getElementById("delete-file-folder-form") as HTMLFormElement;

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
    let deleteAllForm = document.getElementById("delete-all-form") as HTMLFormElement;

    deleteAllForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        let path: string;
        if (currentDir === "") {
            path = "/";
        } else {
            path = currentDir;
        }

        const data: pathInput = {
            path: path
        };
        ajaxHelper(apiRoute + "delete-all", data);
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
