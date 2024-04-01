import logo from "../assets/images/logoTeamA-transformed.png";
import { GetLocalization } from "./localize.js";

document.getElementById("logo").src = logo;

export function CloseNavigationHamburger() {
    if (document.body.offsetWidth < 800) {
        console.log("i closed navigation: " + document.body.offsetWidth);
        document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
    }
}

export function ToggleNavigationHamburger() {
    if (document.body.offsetWidth < 800) {
        console.log("you clicked hamburger: " + document.body.offsetWidth);
        if (document.getElementsByClassName("left-nav")[0].style.visibility === "visible") {
            console.log("great success you toggled to hidden!");
            document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
            return;
        }
        else {
            console.log("great success you toggled to visible!");
            document.getElementsByClassName("left-nav")[0].style.visibility = "visible";
        }
    }
}

export function ToggleNavigationResize() {
    if (document.body.offsetWidth > 799) {
        document.getElementsByClassName("left-nav")[0].style.visibility = "visible";
        console.log("you resized to visible: " + document.body.offsetWidth);
    }
    else {
        document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
        console.log("you resized to invisible: " + document.body.offsetWidth);
    }
}

document.getElementById("header-hambuger").addEventListener("click", () => ToggleNavigationHamburger());
document.body.onresize = () => ToggleNavigationResize();

let navbarItems = ["settings", "home", "security-dashboard","privacy-dashboard", "issues", "integration", "about"]
let localizationIds = ["Navigation.Settings", "Navigation.Home", "Navigation.SecurityDashboard", "Navigation.PrivacyDashboard", "Navigation.Issues", "Navigation.Integration", "Navigation.About"]
for (let i = 0; i < navbarItems.length; i++) {
    GetLocalization(localizationIds[i], navbarItems[i])
}
