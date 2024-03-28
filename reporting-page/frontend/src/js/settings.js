import { openPersonalizePage } from "./personalize";
import { retrieveTheme } from "./personalize";

/** Load the content of the Settings page */
function openSettingsPage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="setting dark-mode">
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
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Personalize page</span>
    <div id="personalize-button-box">
      <div id="personalize-button">Click me</div>
    </div>
  </div> 
  `;

  document.getElementById("personalize-button").addEventListener("click", () => openPersonalizePage());
  document.onload = retrieveTheme();
}

document.getElementById("settings-button").addEventListener("click", () => openSettingsPage());