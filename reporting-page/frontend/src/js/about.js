import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the About page */
export function openAboutPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('about-button');
  sessionStorage.setItem('savedPage', '6');

  document.getElementById('page-contents').innerHTML = `
  <div class="container-about">
    <div class="about-header">
        <h1 class="lang-about-title"></h1>
        p class='lang-about-info'></p>
    </div>
    <div class="project">
        <h2 class="lang-summary-title">Infosec Agent</h2>
        <p class='lang-summary-info'></p>
    </div>
    <div class="project">
        <h2 class="lang-affiliations-title">Little Brother</h2>
        <p class='lang-affiliations-info'>
        </p>
    </div>
    <div class="project">
        <h2 class="lang-contributing-title">Contributing</h2>
        <p class='lang-contributing-info'></p>
    </div>
  </div>`;

  // Localize the static content of the about page
  const staticAboutPageConent = [
    'lang-about-title',
    'lang-about-info',
    'lang-summary-title',
    'lang-summary-info',
    'lang-affiliations-title',
    'lang-affiliations-info',
    'lang-contributing-title',
    'lang-contributing-info',
  ];
  const localizationIds = [
    'About.AboutTitle',
    'About.AboutInfo',
    'About.SummaryTitle',
    'About.SummaryInfo',
    'About.AffiliationsTitle',
    'About.AffiliationsInfo',
    'About.ContributingTitle',
    'About.ContributingInfo',
  ];
  for (let i = 0; i < staticAboutPageConent.length; i++) {
    getLocalization(localizationIds[i], staticAboutPageConent[i]);
  }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('about-button').addEventListener('click', () => openAboutPage());
  } catch (error) {
    logError('Error in about.js: ' + error);
  }
}
