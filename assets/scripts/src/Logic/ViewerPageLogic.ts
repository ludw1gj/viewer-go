import {NotificationHandler} from "../Handler/NotificationHandler";
import {AjaxHandler, JSONErrorResponse} from "../Handler/AjaxHandler";

interface pathInput {
    path: string;
}

class ViewerPageLogic {

    private readonly apiRoute: string = "/api/viewer/";
    private readonly currentDir: string = (document.getElementById("current-dir")  as HTMLElement).innerText.slice(1);

    constructor() {
        this.addEventListenerToUploadFileForm();
        this.addEventListenerToCreateFolderForm();
        this.addEventListenerToDeleteFileFolderForm();
        this.addEventListenerToDeleteAllForm();
    }

    private addEventListenerToUploadFileForm(): void {
        let uploadForm = document.getElementById("upload-form") as HTMLFormElement;

        uploadForm.addEventListener("submit", (event: Event) => {
            event.preventDefault();

            const errFunc = (resp: JSONErrorResponse) => {
                NotificationHandler.displayError(resp.error.message);
            };

            const okFunc = () => {
                location.reload(true);
            };

            AjaxHandler.submitFormData(this.apiRoute + "upload/" + this.currentDir, uploadForm, errFunc, okFunc);
        });
    }

    private addEventListenerToCreateFolderForm(): void {
        let createFolderForm = document.getElementById("create-folder-form") as HTMLFormElement;

        createFolderForm.addEventListener("submit", (event: Event) => {
            event.preventDefault();

            const folderName: HTMLInputElement = createFolderForm.folder_name;

            const data: pathInput = {
                path: ViewerPageLogic.makePath(this.currentDir, folderName.value)
            };
            ViewerPageLogic.ajaxHelper(this.apiRoute + "create", data);
        });
    }

    private addEventListenerToDeleteFileFolderForm(): void {
        let deleteFileFolderForm = document.getElementById("delete-file-folder-form") as HTMLFormElement;

        deleteFileFolderForm.addEventListener("submit", (event: Event) => {
            event.preventDefault();

            const fileName: HTMLInputElement = deleteFileFolderForm.file_name;

            const data: pathInput = {
                path: ViewerPageLogic.makePath(this.currentDir, fileName.value)
            };
            ViewerPageLogic.ajaxHelper(this.apiRoute + "delete", data);
        });
    }

    private addEventListenerToDeleteAllForm(): void {
        let deleteAllForm = document.getElementById("delete-all-form") as HTMLFormElement;

        deleteAllForm.addEventListener("submit", (event: Event) => {
            event.preventDefault();

            let path: string;
            if (this.currentDir === "") {
                path = "/";
            } else {
                path = this.currentDir;
            }

            const data: pathInput = {
                path: path
            };
            ViewerPageLogic.ajaxHelper(this.apiRoute + "delete-all", data);
        });
    }

    private static ajaxHelper(url: string, data: object): void {
        const errFunc = (resp: JSONErrorResponse) => {
            NotificationHandler.displayError(resp.error.message);
        };
        const okFunc = () => {
            location.reload(true);
        };
        AjaxHandler.submitJSON(url, data, errFunc, okFunc)
    }

    private static makePath(currentDir: string, fileName: string): string {
        const index: boolean = (currentDir === "");
        if (index) {
            return fileName;
        }
        return currentDir + "/" + fileName;
    }

}

export {ViewerPageLogic}
