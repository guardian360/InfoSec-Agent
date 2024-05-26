import {openPersonalizePage} from './personalize.js';
import {openAboutPage} from './about.js';
import {openHomePage} from './home.js';
import {openSecurityDashboardPage} from './security-dashboard.js';
import {openPrivacyDashboardPage} from './privacy-dashboard.js';
import {openIssuePage} from './issue.js';
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
  if (sessionStorage.getItem('languageChanged') != null || sessionStorage.getItem('WindowsVersionChanged') != null) {
    const page = sessionStorage.getItem('savedPage');
    switch (page) {
    case '1':
      openHomePage();
      break;
    case '2':
      openSecurityDashboardPage();
      break;
    case '3':
      openPrivacyDashboardPage();
      break;
    case '4':
      openIssuesPage();
      break;
    case '5':
      openIntegrationPage();
      break;
    case '6':
      openAboutPage();
      break;
    case '7':
      openPersonalizePage();
      break;
    case 8:
      const issueId = sessionStorage.getItem('issueId');
      const severity = sessionStorage.getItem('severity');
      openIssuePage(issueId, severity);
      break;
    default:
      try {
        const issuepage = JSON.parse(page);
        openIssuePage(issuepage[0], issuepage[1]);
      } catch {
        logError('Invalid option selected in reloadPage()');
      }
    }
    sessionStorage.removeItem('languageChanged');
    sessionStorage.removeItem('WindowsVersionChanged');
  }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('personalize-button').addEventListener('click', () => openPersonalizePage());
    document.getElementById('language-button').addEventListener('click', () => updateLanguage());
    document.getElementById('windows-version-button').addEventListener('click', () => showWindowsVersion());
    document.getElementById('windows-10').addEventListener('click', () => selectWindowsVersion(10));
    document.getElementById('windows-11').addEventListener('click', () => selectWindowsVersion(11));
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}

/** displays the popup to select the currently used windows version */
export function showWindowsVersion() {
  // Get the modal
  const modal = document.getElementById('window-version-modal');

  modal.style.display = 'block';

  // Get the <span> element that closes the modal
  const close = document.getElementById('close-windows-select');

  // When the user clicks on <span> (x), close the modal
  close.onclick = function() {
    reloadPage();
    modal.style.display = 'none';
  };

  // When the user clicks anywhere outside of the modal, close it
  window.onclick = function(event) {
    if (event.target == modal) {
      reloadPage();
      modal.style.display = 'none';
    }
  };

  const version = sessionStorage.getItem('WindowsVersion');
  document.getElementById('windows-' + version + '-button').classList.add('selected');
}

/**
 * Select the windows version and set it in sessionstorage
 * @param {string} version windows version to select
 */
export function selectWindowsVersion(version) {
  if (version == '10') {
    document.getElementById('windows-10-button').classList.add('selected');
    document.getElementById('windows-11-button').classList.remove('selected');
  } else if (version == '11') {
    document.getElementById('windows-11-button').classList.add('selected');
    document.getElementById('windows-10-button').classList.remove('selected');
  }
  sessionStorage.setItem('WindowsVersion', version);
  sessionStorage.setItem('WindowsVersionChanged', JSON.stringify(true));
}
