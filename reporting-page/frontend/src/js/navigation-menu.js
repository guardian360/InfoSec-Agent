import logo from "../assets/images/logoTeamA-transformed.png";
import { GetLocalization } from "./localize.js";

document.getElementById("logo").src = logo;

export function MarkSelectedNavigationItem() {
    
}

export function CloseNavigationHamburger() {
    if (document.body.offsetWidth < 800) {
        document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
    }
}

export function ToggleNavigationHamburger() {
    if (document.body.offsetWidth < 800) {
        if (document.getElementsByClassName("left-nav")[0].style.visibility === "visible") {
            document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
            return;
        }
        else {
            document.getElementsByClassName("left-nav")[0].style.visibility = "visible";
        }
    }
}

export function ToggleNavigationResize() {
    if (document.body.offsetWidth > 799) {
        document.getElementsByClassName("left-nav")[0].style.visibility = "visible";
    }
    else {
        document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
    }
}

document.getElementById("header-hambuger").addEventListener("click", () => ToggleNavigationHamburger());
document.body.onresize = () => ToggleNavigationResize();

let navbarItems = ["settings", "home", "security-dashboard","privacy-dashboard", "issues", "integration", "about"]
let localizationIds = ["Navigation.Settings", "Navigation.Home", "Navigation.SecurityDashboard", "Navigation.PrivacyDashboard", "Navigation.Issues", "Navigation.Integration", "Navigation.About"]
for (let i = 0; i < navbarItems.length; i++) {
    GetLocalization(localizationIds[i], navbarItems[i])
}
