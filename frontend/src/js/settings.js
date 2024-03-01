import "../css/home.css";
import "../css/settings.css";
import "../css/color-palette.css";

function toggleDarkMode() {
    console.log("toggle dark mode");
    var darkModeSwitch = document.getElementById("dark-mode-switch");

    if (darkModeSwitch.checked == true) {
        document.documentElement.style.setProperty('--background-color', '#333');
        document.documentElement.style.setProperty('--dark-text', '#d4d4d4');
    }
    else {
        document.documentElement.style.setProperty('--background-color', '#ffffff');
        document.documentElement.style.setProperty('--dark-text', '#333');
    }

    
}