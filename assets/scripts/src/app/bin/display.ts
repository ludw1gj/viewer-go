export const displayError = (msg: string) => {
  const notification = document.getElementById("notification") as HTMLElement;

  notification.classList.remove("is-success", "hidden");
  notification.classList.add("is-danger");
  notification.innerText = msg;
};

export const displaySuccess = (msg: string) => {
  const notification = document.getElementById("notification") as HTMLElement;

  notification.classList.remove("is-danger", "hidden");
  notification.classList.add("is-success");
  notification.innerText = msg;
};
