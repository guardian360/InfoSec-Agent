//import logo from "../assets/images/logoTeamA-transformed.png";
import { GetLocalization } from "./localize.js";

//document.getElementById("logo").src = logo;

/** Give the selected navigation item a different color 
 * @param {string} item - The navigation item that is selected
*/
export function MarkSelectedNavigationItem(item) {
    let navItems = document.getElementsByClassName("nav-link");
    let stylesheet = getComputedStyle(document.documentElement)
    for (let i = 1; i < navItems.length; i++) {
        navItems[i].style.backgroundColor = stylesheet.getPropertyValue('--background-color-left-nav');
    }
    if (item === "settings-button" ) {
        return;
    }
    document.getElementById(item).style.backgroundColor = stylesheet.getPropertyValue('--background-nav-hover');
}

export function loadPersonalizeNavigation(){
    let navItems = document.getElementsByClassName("nav-link");
    let stylesheet = getComputedStyle(document.documentElement)
    for (let i = 1; i < navItems.length; i++) {
        navItems[i].style.backgroundColor = stylesheet.getPropertyValue('--background-color-left-nav');
    }
}
/** Close the navigation menu when a navigation item is clicked, only when screen size is less than 800px */
export function CloseNavigation() {
    if (document.body.offsetWidth < 800) {
        document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
    }
}

/** Open or close the navigation menu when user clicks on hamburger menu */
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

/** Open or close the navigation menu when user resizes the screen */
export function ToggleNavigationResize() {
    if (document.body.offsetWidth > 799) {
        document.getElementsByClassName("left-nav")[0].style.visibility = "visible";
    }
    else {
        document.getElementsByClassName("left-nav")[0].style.visibility = "hidden";
    }
}

if (typeof document !== 'undefined') {
    try {
        document.getElementById("header-hambuger").addEventListener("click", () => ToggleNavigationHamburger());
        document.body.onresize = () => ToggleNavigationResize();

        let navbarItems = ["settings", "home", "security-dashboard","privacy-dashboard", "issues", "integration", "about"]
        let localizationIds = ["Navigation.Settings", "Navigation.Home", "Navigation.SecurityDashboard", "Navigation.PrivacyDashboard", "Navigation.Issues", "Navigation.Integration", "Navigation.About"]
        for (let i = 0; i < navbarItems.length; i++) {
            GetLocalization(localizationIds[i], navbarItems[i])
        }
    } catch (error) {
        console.log("Error in navigation-menu.js: " + error)
    }
}
