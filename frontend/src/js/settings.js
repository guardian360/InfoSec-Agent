import "../css/home.css";
import "../css/settings.css";
import "../css/color-palette.css";

function openSettingsPage() {
    document.getElementById("page-contents").innerHTML = `
    <div class="setting" id="dark-mode">
      <span class="setting-description">Dark mode</span>
      <label class="switch">
        <input type="checkbox" id="dark-mode-switch">
        <span class="slider round"></span>
      </label>
    </div>
    <hr class="solid">
    <div class="setting">
      <span class="setting-description">Other setting</span>
      <label class="switch">
        <input type="checkbox">
        <span class="slider round"></span>
      </label>
    </div>
    <hr class="solid">
    <div class="setting">
      <span class="setting-description">Other setting</span>
      <label class="switch">
        <input type="checkbox">
        <span class="slider round"></span>
      </label>
    </div>
    `;

    // function toggleDarkMode() {
    //     console.log("toggle dark mode");
    //     var darkModeSwitch = document.getElementById("dark-mode-switch");
    
    //     if (darkModeSwitch.checked == true) {
    //         document.documentElement.style.setProperty('--background-color', '#333');
    //         document.documentElement.style.setProperty('--dark-text', '#d4d4d4');
    //     }
    //     else {
    //         document.documentElement.style.setProperty('--background-color', '#ffffff');
    //         document.documentElement.style.setProperty('--dark-text', '#333');
    //     }    
    // }

    // document.getElementById("dark-mode-switch").addEventListener("click", () => toggleDarkMode());
}

document.getElementById("settings-button").addEventListener("click", () => openSettingsPage());