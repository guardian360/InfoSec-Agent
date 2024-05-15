import {openPersonalizePage} from './personalize.js';
import {openAboutPage} from './about.js';
import {openHomePage} from './home.js';
import {openSecurityDashboardPage} from './security-dashboard.js';
import {openPrivacyDashboardPage} from './privacy-dashboard.js';
import {openIssuesPage} from './issues.js';
import {openIntegrationPage} from './integration.js';
import {ChangeLanguage as changeLanguage, LogError as logError} from '../../wailsjs/go/main/Tray.js';
import * as runTime from '../../wailsjs/runtime/runtime.js';
/**
 * Initiates a language update operation.
 * Calls the ChangeLanguage function and handles the result or error.
 */
async function updateLanguage() {
  await changeLanguage()
    .then(async () => {
      sessionStorage.setItem('languageChanged', JSON.stringify(true));
      runTime.WindowReload();
    })
    .catch((err) => {
      logError('Error changing language:' + err);
      console.error(err);
    });
}
/** Opens the settings page after window is reloaded after updateLanguage() is called */
if (sessionStorage.getItem('languageChanged') != null) {
  let page = sessionStorage.getItem('savedPage');
  page = parseInt(page);
  switch (page) {
  case 1:
    openHomePage();
    break;
  case 2:
    openSecurityDashboardPage();
    break;
  case 3:
    openPrivacyDashboardPage();
    break;
  case 4:
    openIssuesPage();
    break;
  case 5:
    openIntegrationPage();
    break;
  case 6:
    openAboutPage();
    break;
  case 7:
    openPersonalizePage();
    break;
  default:
    openHomePage();
    console.log('Invalid option selected');
  }
  sessionStorage.removeItem('languageChanged');
}

document.getElementById('personalize-button').addEventListener('click', () => openPersonalizePage());
document.getElementById('language-button').addEventListener('click', () => updateLanguage());

