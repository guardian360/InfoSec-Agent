import {openPersonalizePage} from './personalize.js';
import {openAboutPage} from './about.js';
import {openHomePage} from './home.js';
import {openSecurityDashboardPage} from './security-dashboard.js';
import {openPrivacyDashboardPage} from './privacy-dashboard.js';
import {openIssuesPage} from './issues.js';
import {openIssuePage} from './issue.js';
import {openIntegrationPage} from './integration.js';
import {ChangeLanguage as changeLanguage, LogError as logError} from '../../wailsjs/go/main/Tray.js';

/**
 * Initiates a language update operation.
 * Calls the ChangeLanguage function and handles the result or error.
 */
export async function updateLanguage() {
  await changeLanguage()
    .then(async () => {
      sessionStorage.setItem('languageChanged', JSON.stringify(true));
      reloadPage();
    })
    .catch((err) => {
      logError('Error changing language:' + err);
    });
}

/** opens the page corresponding to the page set in sessionstorage */
export function reloadPage() {
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
    case 8:
      const issueId = sessionStorage.getItem('issueId');
      const severity = sessionStorage.getItem('severity');
      openIssuePage(issueId, severity);
      break;
    default:
      logError('Invalid option selected in reloadPage()');
    }
    sessionStorage.removeItem('languageChanged');
  }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('personalize-button').addEventListener('click', () => openPersonalizePage());
    document.getElementById('language-button').addEventListener('click', () => updateLanguage());
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}
