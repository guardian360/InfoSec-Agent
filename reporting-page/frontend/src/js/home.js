import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

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
        <a class="issue-button lang-suggested-issue"></a>
        <a class="issue-button lang-quick-fix"></a>
        <a id="scan-now" class="issue-button lang-scan-now"></a>
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
    'lang-quick-fix',
    'lang-scan-now',
    'lang-title-medals',
    'lang-choose-issue-description',
  ];
  const localizationIds = [
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.QuickFix',
    'Dashboard.ScanNow',
    'Dashboard.Medals',
    'Dashboard.ChooseIssueDescription',
  ];
  for (let i = 0; i < staticHomePageContent.length; i++) {
    getLocalization(localizationIds[i], staticHomePageContent[i]);
  }

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
