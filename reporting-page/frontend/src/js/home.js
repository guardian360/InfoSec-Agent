import {PieChart} from './piechart.js';
import {LogMessage, ScanNow as scanNowGo} from '../../wailsjs/go/main/Tray.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import medal from '../assets/images/img_medal1.jpg';
import {retrieveTheme} from './personalize.js';

/** Load the content of the Home page */
export function openHomePage() {
  LogMessage('Opening Home Page');
  closeNavigation();
  markSelectedNavigationItem('home-button');

  document.getElementById('page-contents').innerHTML = `
  <div class="home-data">
    <div class="container-data home-column-one"> 
      <div class="data-column risk-counters">     
        <div class="data-column data-segment piechart">
          <div class="data-segment-header">
            <p class="piechart-header">Risk level distribution</p>
          </div>
          <div class="piechart-container">
            <canvas id="pieChart"></canvas>
          </div>
        </div>
      </div>
      <div class="data-column data-segment issue-buttons">
        <div class="data-segment-header">
          <p class="choose-issue-description"></p>
        </div>
        <a class="issue-button suggested-issue">Suggested Issue</a>
        <a class="issue-button quick-fix">Quick Fix</a>
        <a class="issue-button scan-now">Scan Now</a>
      </div>
    </div>
    <div class="data-segment progress">  
      <div class="data-segment-header">
        <p class="title-medals"></p>
      </div>
      <div class="medals">
        <div class="medal-layout">
          <img id="medal" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 1</p>
          <p> 01-04-2024</p>
        </div>
        <div class="medal-layout">
          <img id="medal2" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 2</p>
          <p> 01-04-2024</p>
        </div>
        <div class="medal-layout">
          <img id="medal3" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 3</p>
          <p> 01-04-2024</p>
        </div><div class="medal-layout">
          <img id="medal4" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 4</p>
          <p> 01-04-2024</p>
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
  new PieChart('pieChart', rc);

  // Localize the static content of the home page
  const staticHomePageConent = [
    'suggested-issue',
    'quick-fix',
    'scan-now',
    'title-medals',
    'security-status',
    'high-risk-issues',
    'medium-risk-issues',
    'low-risk-issues',
    'safe-issues',
    'choose-issue-description',
  ];
  const localizationIds = [
    'Dashboard.SuggestedIssue',
    'Dashboard.QuickFix',
    'Dashboard.ScanNow',
    'Dashboard.Medals',
    'Dashboard.SecurityStatus',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.Safe',
    'Dashboard.ChooseIssueDescription',
  ];
  for (let i = 0; i < staticHomePageConent.length; i++) {
    getLocalization(localizationIds[i], staticHomePageConent[i]);
  }

  document.getElementsByClassName('scan-now')[0].addEventListener('click', () => scanNow());
  document.getElementById('home-button').addEventListener('click', () => openHomePage());
  document.getElementById('logo').innerHTML = localStorage.getItem('picture');

  document.onload = retrieveTheme();
}

document.getElementById('logo-button').addEventListener('click', () => openHomePage());
document.getElementById('home-button').addEventListener('click', () => openHomePage());

/**
 * Initiates a scan operation immediately.
 * Calls the ScanNow function and handles the result or error.
 */
function scanNow() {
  scanNowGo()
    .then((result) => {
    })
    .catch((err) => {
      console.error(err);
    });
}

// document.onload = openHomePage();

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
