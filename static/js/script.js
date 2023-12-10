const menuBtn = document.getElementById("menu-btn");
const mobileNav = document.getElementById("mobile-nav");

menuBtn.addEventListener("click", navToggle)

function navToggle() {
    menuBtn.classList.toggle("open");
    mobileNav.classList.toggle("flex");
    mobileNav.classList.toggle("hidden");
}