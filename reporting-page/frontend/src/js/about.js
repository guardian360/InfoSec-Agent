import {closeNavigation, markSelectedNavigationItem} from './navigation-menu';
import {retrieveTheme} from './personalize';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the About page */
export function openAboutPage() {
  document.onload = retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('about-button');
  sessionStorage.setItem('savedPage', '6');

  document.getElementById('page-contents').innerHTML = `
  <div class="container-about">
    <div class="about-header">
        <h1 class="about-title">About Infosec Agent</h1>
    </div>
    <div class="project-info">
        <div class="project">
            <h2>Infosec Agent</h2>
            <p class='infosec-info'>The InfoSec Agent project aims to improve 
                the security and privacy of Windows computer users. 
                Currently, 
                there are applications available that do this, but they are mainly targeted at large companies. 
                The goal of this project is to make this accessible to everyone.
                An application is being developed that collects information 
                about the user's system to discover any security or privacy-related vulnerabilities. 
                The results will be presented to the user in a special dashboard, 
                showing the current status of the system, including recommended actions to improve it.</p>
        </div>
        <div class="project">
            <h2>Little Brother</h2>
            <p class='little-brother-info'>
              This project is a collaborative effort involving nine students from Utrecht University in The Netherlands
              , in partnership with the Dutch IT company Guardian360. 
              It serves as the Software Project for the Bachelor's Programme in Computing Sciences at the UU.
            </p>
        </div>
    </div>
    <div class="contribute">
        <h2 class="about-contribute">Contributing</h2>
        <p class='contribute-info'>InfoSec-Agent is an Open-Source project licensed under the GPL-3.0 License. 
            However, due to its origins as a Utrecht University assignment, 
            public contributions to this repository will only be merged after the completion of this assignment, 
            which is scheduled for June 24, 2024.
            Feel free to report any bugs or issues you encounter. 
            Your feedback is valuable and helps improve the InfoSec-Agent project.
            You can email us at <a href="mailto:infosecagentuu@gmail.com">infosecagentuu@gmail.com</a>.</p>
    </div>
</div>`;

  // Localize the static content of the home page
  const staticAboutPageConent = [
    'about-title',
    'infosec-info',
    'little-brother-info',
    'about-contribute',
    'contribute-info',
  ];
  const localizationIds = [
    'About.about-title',
    'About.infosec-info',
    'About.little-brother-info',
    'About.about-contribute',
    'About.contribute-info',
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
