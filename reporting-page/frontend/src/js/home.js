import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import medal from '../assets/images/img_medal1.png';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';

/** Load the content of the Home page */
export function openHomePage() {
  document.onload = retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('home-button');
  sessionStorage.setItem('savedPage', 1);
  document.getElementById('page-contents').innerHTML = `
  <div class="home">
    <div class="container-home"> 
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="piechart-header">Risk level distribution</p>
        </div>
        <div class="pie-chart-container">
          <canvas id="pie-chart"></canvas>
        </div>
      </div>
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="choose-issue-description">Select an issue</p>
        </div>
        <a class="issue-button suggested-issue"></a>
        <a class="issue-button quick-fix">Quick Fix</a>
        <a class="issue-button scan-now">Scan Now</a>
      </div>
    </div>
    <div class="container-home"> 
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="title-medals"></p>
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

  document.getElementById('medal').src = medal;
  document.getElementById('medal2').src = medal;
  document.getElementById('medal3').src = medal;
  document.getElementById('medal4').src = medal;

  const rc = JSON.parse(sessionStorage.getItem('RiskCounters'));
  new PieChart('pie-chart', rc, 'Total');

  // Localize the static content of the home page
  const staticHomePageContent = [
    'piechart-header',
    'suggested-issue',
    'quick-fix',
    'scan-now',
    'title-medals',
    'security-status',
    'high-risk-issues',
    'medium-risk-issues',
    'low-risk-issues',
    'info-risk-issues',
    'safe-issues',
    'choose-issue-description',
  ];
  const localizationIds = [
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.QuickFix',
    'Dashboard.ScanNow',
    'Dashboard.Medals',
    'Dashboard.SecurityStatus',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.InfoRisk',
    'Dashboard.Safe',
    'Dashboard.ChooseIssueDescription',
  ];
  for (let i = 0; i < staticHomePageContent.length; i++) {
    getLocalization(localizationIds[i], staticHomePageContent[i]);
  }

  document.getElementsByClassName('scan-now')[0].addEventListener('click', () => scanTest());
  document.getElementById('logo').innerHTML = localStorage.getItem('picture');
}

document.getElementById('logo-button').addEventListener('click', () => openHomePage());
document.getElementById('home-button').addEventListener('click', () => openHomePage());

window.onload = function() {
  markSelectedNavigationItem('home-button');
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
