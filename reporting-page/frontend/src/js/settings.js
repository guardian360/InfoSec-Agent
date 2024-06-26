import {openPersonalizePage} from './personalize.js';
import {openAboutPage} from './about.js';
import {openHomePage} from './home.js';
import {openSecurityDashboardPage} from './security-dashboard.js';
import {openPrivacyDashboardPage} from './privacy-dashboard.js';
import {openIssuePage} from './issue.js';
import {openIssuesPage} from './issues.js';
import {openIntegrationPage} from './integration.js';
import {openAllChecksPage} from './all-checks.js';
import {openProgramsPage} from './programs.js';

import {ChangeLanguage as changeLanguage,
  ChangeScanInterval as changeScanInterval,
  LogError as logError} from '../../wailsjs/go/main/Tray.js';

// On reload makes sure modal is openable
sessionStorage.removeItem('ModalOpen');

/**
 * Initiates a language update operation.
 * Calls the ChangeLanguage function and handles the result or error.
 */
export async function updateLanguage() {
  // remove modal is open flag if it exists
  sessionStorage.removeItem('ModalOpen');
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
      openProgramsPage();
      break;
    case '6':
      openAllChecksPage();
      break;
    case '7':
      openIntegrationPage();
      break;
    case '8':
      openAboutPage();
      break;
    case '9':
      openPersonalizePage();
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
    document.getElementById('windows-version-button')
      .addEventListener('click', () => showModal('window-version-modal'));
    document.getElementById('scan-interval-button').addEventListener('click', () => changeScanInterval());

    document.getElementById('windows-10').addEventListener('click', () => selectWindowsVersion(10));
    document.getElementById('windows-11').addEventListener('click', () => selectWindowsVersion(11));
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}

/**
 * displays the popup to select the currently used windows version
 * @param {string} id id of modal element
 */
export function showModal(id) {
  const open = sessionStorage.getItem('ModalOpen');
  if (open == undefined || open == null) {
    sessionStorage.setItem('ModalOpen', true);
    // Get the modal
    const modal = document.getElementById(id);

    modal.style.display = 'block';

    // Get the <span> element that closes the modal
    const close = document.getElementById('close-' + id);

    // When the user clicks on <span> (x), close the modal
    close.onclick = function() {
      reloadPage();
      modal.style.display = 'none';
      sessionStorage.removeItem('ModalOpen');
    };

    // When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
      if (event.target == modal) {
        reloadPage();
        modal.style.display = 'none';
        sessionStorage.removeItem('ModalOpen');
      }
    };

    if (id === 'window-version-modal') {
      const version = sessionStorage.getItem('WindowsVersion');
      document.getElementById('windows-' + version + '-button').classList.add('selected');
    }
  }
}

/**
 * Select the windows version and set it in sessionstorage
 * @param {string} version windows version to select
 */
export function selectWindowsVersion(version) {
  if (version == '10') {
    document.getElementById('windows-10-button').classList.add('selected');
    document.getElementById('windows-11-button').classList.remove('selected');
  } else {
    document.getElementById('windows-11-button').classList.add('selected');
    document.getElementById('windows-10-button').classList.remove('selected');
  }
  sessionStorage.setItem('WindowsVersion', version);
  sessionStorage.setItem('WindowsVersionChanged', JSON.stringify(true));
}
