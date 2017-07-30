// pathInput contains the required data structure.
interface pathInput {
    path: string;
}

// addEventListenersViewerForms function should be run at initialisation of viewer page.
function addEventListenersViewerForms(): void {
    const apiRoute: string = "/api/viewer/";
    const currentDir: string = (document.getElementById("current-dir")  as HTMLElement).innerText.slice(1);

    // handle upload form logic
    let uploadForm = document.getElementById("upload-form") as HTMLFormElement;
    uploadForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const errFunc = (resp: JsonErrorResponse) => {
            displayErrorNotification(resp.error.message);
        };

        const okFunc = () => {
            location.reload(true);
        };
        submitAjaxFormData(apiRoute + "upload", uploadForm, errFunc, okFunc);
    });

    // handle create folder form logic
    let createFolderForm = document.getElementById("create-folder-form") as HTMLFormElement;
    createFolderForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const folderName: HTMLInputElement = createFolderForm.folder_name;

        const data: pathInput = {
            path: makePath(currentDir, folderName.value)
        };
        viewerAjaxHelper(apiRoute + "create", data);
    });

    // handle delete file/folder form logic
    let deleteFileFolderForm = document.getElementById("delete-file-folder-form") as HTMLFormElement;
    deleteFileFolderForm.addEventListener("submit", (event: Event) => {
        event.preventDefault();

        const fileName: HTMLInputElement = deleteFileFolderForm.file_name;

        const data: pathInput = {
            path: makePath(currentDir, fileName.value)
        };
        viewerAjaxHelper(apiRoute + "delete", data);
    });

    // handle delete all form logic
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
        viewerAjaxHelper(apiRoute + "delete-all", data);
    });
}

// viewerAjaxHelper is a wrapper for submitAjaxJson function.
function viewerAjaxHelper(url: string, data: object): void {
    const errFunc = (resp: JsonErrorResponse) => {
        displayErrorNotification(resp.error.message);
    };
    const okFunc = () => {
        location.reload(true);
    };
    submitAjaxJson(url, data, errFunc, okFunc)
}

// makePath generates a path.
function makePath(currentDir: string, fileName: string): string {
    const index: boolean = (currentDir === "");
    if (index) {
        return fileName;
    } else {
        return currentDir + "/" + fileName;
    }
}
