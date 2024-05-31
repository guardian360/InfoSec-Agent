import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {openIssuePage} from './issue.js';
import data from '../databases/database.en-GB.json' assert { type: 'json' };

/** Load the content of the Home page */
export function openHomePage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('home-button');
  sessionStorage.setItem('savedPage', 1);

  document.getElementById('page-contents').innerHTML = `
  <video autoplay muted loop class="video-background">
        <source id="lighthouse-background" type="video/mp4">
        Your browser does not support HTML5 video.
    </video>
  <div class="home-page">
    <div class="container-home"> 
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-piechart-header">Risk level distribution</p>
        </div>
        <div class="pie-chart-container">
          <canvas class="pie-chart" id="pie-chart-home"></canvas>
        </div>
      </div>
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="lang-choose-issue-description">Actions</p>
        </div>
        <a id="suggested-issue" class="issue-button lang-suggested-issue">Suggested Issue</a>
        <a id="scan-now" class="issue-button lang-scan-now">Scan Now</a>
      </div>
    </div>
  </div>
  `;

  const lighthouseState = 'src/assets/images/regular1.mp4';
  document.getElementById('lighthouse-background').src = lighthouseState;

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
