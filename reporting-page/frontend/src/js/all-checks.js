import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the About page */
export function openAllChecksPage() {
  console.log('test');
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('all-checks-button');
  sessionStorage.setItem('savedPage', '8');

  document.getElementById('page-contents').innerHTML = `
  <div class="container-check-types">
  </div>`;

  // Localize the static content of the about page
  const staticAboutPageConent = [
  ];
  const localizationIds = [
  ];
  for (let i = 0; i < staticAboutPageConent.length; i++) {
    getLocalization(localizationIds[i], staticAboutPageConent[i]);
  }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('all-checks-button').addEventListener('click', () => openAllChecksPage());
  } catch (error) {
    logError('Error in all-checks.js: ' + error);
  }
}