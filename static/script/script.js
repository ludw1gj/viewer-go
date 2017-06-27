var burgerIcon = document.getElementById("burger-icon");
var burgerMenu = document.getElementById("burger-menu");

burgerIcon.addEventListener("click", function () {
    if (burgerIcon.classList.contains("is-active") || burgerMenu.classList.contains("is-active")) {
        burgerIcon.classList.remove("is-active");
        burgerMenu.classList.remove("is-active");
    } else {
        burgerIcon.classList.add("is-active");
        burgerMenu.classList.add("is-active");
    }
}, false);
