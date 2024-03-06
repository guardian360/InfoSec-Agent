import "../css/home.css";
import "../css/color-palette.css";

function openHomePage() {
    document.getElementById("page-contents").innerHTML = ``;
}

document.getElementById("logo-button").addEventListener("click", () => openHomePage());
document.getElementById("home-button").addEventListener("click", () => openHomePage());