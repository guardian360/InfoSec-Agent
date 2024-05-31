import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {openIssuePage} from './issue.js';
import {saveProgress, shareProgress, selectSocialMedia} from './share.js';
import data from '../databases/database.en-GB.json' assert { type: 'json' };
import { showModal } from './settings.js';

/** Load the content of the Home page */
export function openHomePage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('home-button');
  sessionStorage.setItem('savedPage', 1);

  document.getElementById('page-contents').innerHTML = `
  <div class="home-page">
    <div class="container-home"> 
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-piechart-header"></p>
        </div>
        <div class="pie-chart-container">
          <canvas class="pie-chart" id="pie-chart-home"></canvas>
        </div>
      </div>
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-choose-issue-description"></p>
        </div>
        <a id="suggested-issue" class="issue-button lang-suggested-issue"></a>
        <a id="scan-now" class="issue-button lang-scan-now"></a>
        <a id="share-progress" class="issue-button">Share progress</a>
      </div>
    </div>
    <div class="container-home"> 
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-title-medals"></p>
        </div>
        <div class="medals">
          <div class="medal-layout">
            <img id="medal" alt="Photo of medal"></img>
            <p class="medal-name"> Medal 1</p>
          </div>
          <div class="medal-layout">
            <img id="medal2" alt="Photo of medal"></img>
            <p class="medal-name"> Medal 2</p>
          </div>
          <div class="medal-layout">
            <img id="medal3" alt="Photo of medal"></img>
            <p class="medal-name"> Medal 3</p>
          </div><div class="medal-layout">
            <img id="medal4" alt="Photo of medal"></img>
            <p class="medal-name"> Medal 4</p>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div id="share-modal" class="modal">
    <div class="modal-content">
      <div class="modal-header">
        <span id="close-share-modal" class="close">&times;</span>
        <p>Select where to share your progress, Save and download it, then share it with others!</p>
      </div>
      <div id="share-node" class="modal-body"><img class="api-key-image" src="https://placehold.co/600x315" alt="Step 1 Image"></div>
      <div id="share-buttons" class="modal-body">
        <a id="share-save-button" class="modal-button share-button">Save</a>
        <a class="share-button-break">|</a>
        <a id="select-facebook" class="select-button selected">Facebook</a>
        <a id="select-x" class="select-button">X</a>
        <a id="select-linkedin" class="select-button">LinkedIn</a>
        <a id="select-instagram" class="select-button">Instagram</a>
        <a class="share-button-break">|</a>
        <a id="share-button" class="modal-button share-button">Share</a>
      </div>
    </div>
  </div>
  `;

  const medal = 'frontend/src/assets/images/img_medal1.png';
  document.getElementById('medal').src = medal;
  document.getElementById('medal2').src = medal;
  document.getElementById('medal3').src = medal;
  document.getElementById('medal4').src = medal;

  const rc = JSON.parse(sessionStorage.getItem('RiskCounters'));
  new PieChart('pie-chart-home', rc, 'Total');

  // Localize the static content of the home page
  const staticHomePageContent = [
    'lang-piechart-header',
    'lang-suggested-issue',
    'lang-scan-now',
    'lang-title-medals',
    'lang-choose-issue-description',
  ];
  const localizationIds = [
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.ScanNow',
    'Dashboard.Medals',
    'Dashboard.ChooseIssueDescription',
  ];
  for (let i = 0; i < staticHomePageContent.length; i++) {
    getLocalization(localizationIds[i], staticHomePageContent[i]);
  }

  document.getElementById('scan-now').addEventListener('click', () => scanTest(true));
  document.getElementById('suggested-issue').addEventListener('click', () => suggestedIssue(''));
  document.getElementById('share-progress').addEventListener('click', () => showModal('share-modal'));
  document.getElementById('share-save-button').addEventListener('click', 
    () => saveProgress(document.getElementById('share-node')));
  document.getElementById('share-button').addEventListener('click', () => shareProgress());

  document.getElementById('select-facebook').addEventListener('click', () => selectSocialMedia('facebook'));
  document.getElementById('select-x').addEventListener('click', () => selectSocialMedia('x'));
  document.getElementById('select-linkedin').addEventListener('click', () => selectSocialMedia('linkedin'));
  document.getElementById('select-instagram').addEventListener('click', () => selectSocialMedia('instagram'));
}

/** Opens the issue page of the issue with highest risk level
 *
 * @param {string} type Type of issue to open the issue page of (e.g. 'Security', 'Privacy', and '' for all types)
*/
export function suggestedIssue(type) {
  // Get the issues from the database
  const issues = JSON.parse(sessionStorage.getItem('DataBaseData'));

  // Skip informative issues
  let issueIndex = 0;
  let maxSeverityIssue = issues[issueIndex];
  while (maxSeverityIssue.severity === 4) {
    issueIndex++;
    maxSeverityIssue = issues[issueIndex];
  }

  // Find the issue with the highest severity
  for (let i = 0; i < issues.length; i++) {
    if (maxSeverityIssue.severity < issues[i].severity && issues[i].severity !== 4) {
      if (type == '' || data[issues[i].jsonkey].Type === type) {
        maxSeverityIssue = issues[i];
      }
    }
  }

  // Open the issue page of the issue with the highest severity
  openIssuePage(maxSeverityIssue.jsonkey, maxSeverityIssue.severity);
  document.getElementById('scan-now').addEventListener('click', () => scanTest(true));
}


if (typeof document !== 'undefined') {
  try {
    document.getElementById('logo-button').addEventListener('click', () => openHomePage());
    document.getElementById('home-button').addEventListener('click', () => openHomePage());
  } catch (error) {
    /* istanbul ignore next */
    logError('Error in security-dashboard.js: ' + error);
  }
}


window.onload = function() {
  const savedImage = localStorage.getItem('picture');
  const savedText = localStorage.getItem('title');
  const savedIcon = localStorage.getItem('favicon');
  if (savedImage) {
    const logo = document.getElementById('logo');
    logo.src = savedImage;
  }
  if (savedText) {
    const title = document.getElementById('title');
    title.textContent = savedText;
  }
  if (savedIcon) {
    const favicon = document.getElementById('favicon');
    favicon.href = savedIcon;
  }
};


