import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {getLocalization} from './localize.js';

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
    </div>
    <div class="project-info">
        <div class="project">
            <h2>Infosec Agent</h2>
            <p class='lang-infosec-info'></p>
        </div>
        <div class="project">
            <h2>Little Brother</h2>
            <p class='lang-little-brother-info'>
            </p>
        </div>
    </div>
    <div class="contribute">
        <h2 class="lang-contribute-title">Contributing</h2>
        <p class='lang-contribute-info'></p>
    </div>
  </div>`;

  // Localize the static content of the about page
  const staticAboutPageConent = [
    'lang-about-title',
    'lang-infosec-info',
    'lang-little-brother-info',
    'lang-contribute-title',
    'lang-contribute-info',
  ];
  const localizationIds = [
    'About.about-title',
    'About.infosec-info',
    'About.little-brother-info',
    'About.contribute-title',
    'About.contribute-info',
  ];
  for (let i = 0; i < staticAboutPageConent.length; i++) {
    getLocalization(localizationIds[i], staticAboutPageConent[i]);
  }
}

document.getElementById('about-button').addEventListener('click', () => openAboutPage());
