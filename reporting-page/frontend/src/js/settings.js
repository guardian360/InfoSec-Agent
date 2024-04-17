import {openPersonalizePage} from './personalize.js';
import {ChangeLanguage as changeLanguage, LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import * as runTime from '../../wailsjs/runtime/runtime.js';
/**
 * Initiates a language update operation.
 * Calls the ChangeLanguage function and handles the result or error.
 */
async function updateLanguage() {
  await changeLanguage()
    .then(async (result) => {
      sessionStorage.setItem('languageChanged',JSON.stringify(true))
      runTime.WindowReload();
      openSettingsPage();
    })
    .catch((err) => {
      logError('Error changing language:' + err);
      console.error(err);
    });
}
/** Opens the settings page after window is reloaded after updateLanguage() is called */
if (sessionStorage.getItem('languageChanged') != null) {
  openSettingsPage();
  sessionStorage.removeItem('languageChanged');
}

/** Load the content of the Settings page */
function openSettingsPage() {
  closeNavigation();
  markSelectedNavigationItem('settings-button');

  document.getElementById('page-contents').innerHTML = `
  <div class="setting personalize">
    <span class="setting-description personalize-title">Personalization</span>
    <button class="setting-button personalize-button" type="button">Personalize</button>    
  </div> 
  <hr class="solid">
  <div class="setting language">
    <span class="setting-description language-title">Language</span>
    <button class="setting-button language-button" type="button">Change Language</button>
  </div> 
  `;

  // Localize the static content of the settings page
  const staticSettingsContent = ['personalize-title', 'personalize-button', 'language-title', 'language-button'];
  const localizationIds = [
    'Settings.PersonalizeTitle',
    'Settings.PersonalizeButton',
    'Settings.ChangeLanguageTitle',
    'Settings.ChangeLanguageButton',
  ];
  for (let i = 0; i < staticSettingsContent.length; i++) {
    getLocalization(localizationIds[i], staticSettingsContent[i]);
  }

  document.getElementsByClassName('language-button')[0].addEventListener('click', () => updateLanguage());
  document.getElementsByClassName('personalize-button')[0].addEventListener('click', () => openPersonalizePage());
  document.onload = retrieveTheme();
}

document.getElementById('settings-button').addEventListener('click', () => openSettingsPage());
