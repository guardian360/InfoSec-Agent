import {Graph} from './graph.js';
import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {adjustWithRiskCounters, setMaxInterval, addGraphFunctions} from './security-dashboard.js';
import {scanTest} from './database.js';
import {suggestedIssue} from './home.js';
import {openAllChecksPage} from './all-checks.js';

/** Load the content of the Privacy Dashboard page */
export function openPrivacyDashboardPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('privacy-dashboard-button');
  sessionStorage.setItem('savedPage', '3');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard">
    <div class="container-dashboard"> <!-- title top container -->
      <div class="dashboard-segment dashboard-title"> <!-- title top segment -->
        <p class="lang-privacy-dashboard"><p> 
      </div>
    </div>
    <div class="container-dashboard"> <!-- top container -->
      <div class="dashboard-segment"> <!-- Privacy status segment -->
        <div class="data-segment-header">
          <p class="lang-privacy-stat"></p>
        </div>
        <div class="security-status">
          <p class="status-descriptor"></p>
        </div>
      </div> 
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="lang-choose-issue-description"></p>
        </div>
        <a id="suggested-issue" class="privacy-button lang-suggested-issue"><p></p></a>
        <a id="scan-now" class="privacy-button lang-scan-now"></a>
      </div>
      <div class="dashboard-segment risk-areas">
        <div class="data-segment-header">
          <p class="lang-privacy-risk-areas"></p>
        </div>
        <div class="security-area-buttons">
          <div class="security-area privacy-risk-button" id="privacy-button-permissions">
            <a>
              <p>
                <span class="lang-permissions"></span>
                <span class="material-symbols-outlined">person_check</span>
              </p>
            </a>
          </div>
          <div class="security-area privacy-risk-button" id="privacy-button-browser">
            <a>
              <p><span class="lang-browser"></span><span class="material-symbols-outlined">travel_explore</span></p>
            </a>
          </div>
          <div class="security-area privacy-risk-button" id="privacy-button-other">
            <a>
              <p><span class="lang-other"></span><span class="material-symbols-outlined">view_cozy</span></p>
            </a>
          </div>
        </div>
      </div>    
    </div>
    <div class="container-dashboard"> <!-- bottom container -->
      <div class="dashboard-segment"> <!-- Privacy risk counters segment -->
        <div class="data-segment-header">
          <p class="lang-risk-level-counters"></p>
        </div>
        <div class="risk-counter high-risk">
          <div><p class="lang-high-risk-issues"></p></div>
          <div><p id="high-risk-counter">0</p></div>
        </div>
        <div class="risk-counter medium-risk">
          <div><p class="lang-medium-risk-issues"></p></div>
          <div><p id="medium-risk-counter">0</p></div>
        </div>
        <div class="risk-counter low-risk">
          <div><p class="lang-low-risk-issues"></p></div>
          <div><p id="low-risk-counter">0</p></div>
        </div>
        <div class="risk-counter info-risk">
          <div><p class="lang-info-risk-issues"></p></div>
          <div><p id="info-risk-counter">0</p></div>
        </div>
        <div class="risk-counter no-risk">
          <div><p class="lang-acceptable-issues"></p></div>
          <div><p id="no-risk-counter">0</p></div>
        </div>
      </div> 
      <div class="dashboard-segment"> <!-- pie chart segment -->
        <div class="data-segment-header">
            <p class="lang-risk-level-distribution piechart-header"></p>
        </div>
        <div class="pie-chart-container">
          <canvas class="pie-chart" id="pie-chart-privacy"></canvas>
        </div>
      </div>
      <div class="dashboard-segment"> <!-- graph segement -->
        <div class="data-segment-header">
          <p class="lang-risk-level-distribution"></p>
        </div>
        <div class="graph-segment-content">
          <div class="graph-buttons">
            <p class="lang-bar-graph-description">
            </p>
            <button id="dropbtn" class="dropbtn privacy-button"><span class="lang-select-risks"></span></button>
            <div class="dropdown-selector" id="myDropdown">
              <p><input type="checkbox" checked="true" value="true" id="select-high-risk">
                <label for="select-high-risk" class="lang-high-risk-issues"></label><br>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-medium-risk">
                <label for="select-medium-risk" class="lang-medium-risk-issues"></label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-low-risk">
                <label for="select-low-risk" class="lang-low-risk-issues"></label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-info-risk">
                <label for="select-info-risk" class="lang-info-risk-issues"></label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-no-risk">
                <label for="select-no-risk" class="lang-acceptable-issues"></label>
              </p>
            </div>
            <a class="interval-button">
              <p class="lang-change-interval"></p>
              <input type="number" value="1" id="graph-interval" min="1">
            </a>
          </div>
          <div class="interval-graph-container">
            <canvas id="interval-graph"></canvas>
          </div>
        </div>
      </div>
    </div>
  </div>
  `;
  // Set counters on the page to the right values
  let rc = JSON.parse(sessionStorage.getItem('PrivacyRiskCounters'));
  adjustWithRiskCounters(rc, document, true);
  setMaxInterval(rc, document);

  // Localize the static content of the dashboard
  const staticDashboardContent = [
    'lang-security-dashboard',
    'lang-issues',
    'lang-high-risk-issues',
    'lang-medium-risk-issues',
    'lang-low-risk-issues',
    'lang-info-risk-issues',
    'lang-acceptable-issues',
    'lang-privacy-stat',
    'lang-risk-level-counters',
    'lang-risk-level-distribution',
    'lang-suggested-issue',
    'lang-scan-now',
    'lang-privacy-risk-areas',
    'lang-permissions',
    'lang-browser',
    'lang-other',
    'lang-select-risks',
    'lang-change-interval',
    'lang-choose-issue-description',
    'lang-bar-graph-description',
  ];
  const localizationIds = [
    'Navigation.PrivacyDashboard',
    'Dashboard.Issues',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.InfoRisk',
    'Dashboard.Acceptable',
    'Dashboard.PrivacyStatus',
    'Dashboard.RiskLevelCounters',
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.ScanNow',
    'Dashboard.PrivacyRiskAreas',
    'Dashboard.Permissions',
    'Dashboard.Browser',
    'Dashboard.Other',
    'Dashboard.SelectRisks',
    'Dashboard.ChangeInterval',
    'Dashboard.ChooseIssueDescription',
    'Dashboard.BarGraphDescription',
  ];
  for (let i = 0; i < staticDashboardContent.length; i++) {
    getLocalization(localizationIds[i], staticDashboardContent[i]);
  }

  // Create charts
  new PieChart('pie-chart-privacy', rc, 'Privacy');
  const g = new Graph('interval-graph', rc);
  addGraphFunctions(g);
  document.getElementById('scan-now').addEventListener('click', async () => {
    await scanTest(true);
    rc = JSON.parse(sessionStorage.getItem('PrivacyRiskCounters'));
    adjustWithRiskCounters(rc, document, true);
    setMaxInterval(rc, document);
    g.rc = rc;
    await g.changeGraph();
  });
  document.getElementById('suggested-issue').addEventListener('click', () => suggestedIssue('Privacy'));

  // Add links to checks page
  document.getElementById('privacy-button-permissions').addEventListener('click',
    () => openAllChecksPage('permissions'));
  document.getElementById('privacy-button-browser').addEventListener('click',
    () => openAllChecksPage('browser'));
  document.getElementById('privacy-button-other').addEventListener('click',
    () => openAllChecksPage('privacy-other'));
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('privacy-dashboard-button').addEventListener('click', () => openPrivacyDashboardPage());
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}
