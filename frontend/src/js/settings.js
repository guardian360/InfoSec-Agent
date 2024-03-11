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
}

document.getElementById("settings-button").addEventListener("click", () => openSettingsPage());