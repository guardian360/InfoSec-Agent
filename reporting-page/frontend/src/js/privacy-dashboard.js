import {Graph} from './graph.js';
import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {adjustWithRiskCounters, setMaxInterval, addGraphFunctions} from './security-dashboard.js';
import {scanTest} from './database.js';

/** Load the content of the Privacy Dashboard page */
export function openPrivacyDashboardPage() {
  document.onload = retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('privacy-dashboard-button');
  sessionStorage.setItem('savedPage', '3');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard">
    <div class="container-dashboard">
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="privacy-stat">Privacy status</p>
        </div>
        <div class="security-status">
          <p class="status-descriptor"></p>
        </div>
      </div>
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p>Risk level counters</p>
        </div>
        <div class="risk-counter high-risk">
          <div><p class="high-risk-issues">High risk issues</p></div>
          <div><p id="high-risk-counter">0</p></div>
        </div>
        <div class="risk-counter medium-risk">
          <div><p class="medium-risk-issues">Medium risk issues</p></div>
          <div><p id="medium-risk-counter">0</p></div>
        </div>
        <div class="risk-counter low-risk">
          <div><p class="low-risk-issues">Low risk issues</p></div>
          <div><p id="low-risk-counter">0</p></div>
        </div>
        <div class="risk-counter info-risk">
          <div><p class="info-risk-issues">Informative</p></div>
          <div><p id="info-risk-counter">0</p></div>
        </div>
        <div class="risk-counter no-risk">
          <div><p class="safe-issues">Safe issues</p></div>
          <div><p id="no-risk-counter">0</p></div>
        </div>
      </div>      
    </div>
    <div class="container-dashboard">
      <div class="dashboard-segment">
        <div class="data-segment-header">
            <p class="piechart-header">Risk level distribution</p>
        </div>
        <div class="pie-chart-container">
          <canvas id="pie-chart"></canvas>
        </div>
      </div>
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p>Risk level distribution</p>
        </div>
        <div class="graph-segment-content">
          <div class="graph-buttons">
            <p class="bar-graph-description">
              In this graph you are able to see the distribution of different issues 
              we have found over the past times we ran a check.
            </p>
            <button id="dropbtn" class="dropbtn"><span class="select-risks">Select Risks</span></button>
            <div class="dropdown-selector" id="myDropdown">
              <p><input type="checkbox" checked="true" value="true" id="select-high-risk">
                <label for="select-high-risk" class="high-risk-issues"> High risks</label><br>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-medium-risk">
                <label for="select-medium-risk" class="medium-risk-issues"> Medium risks</label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-low-risk">
                <label for="select-low-risk" class="low-risk-issues"> Low risks</label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-info-risk">
                <label for="select-info-risk" class="info-risk-issues"> Informative</label>
              </p>
              <p><input type="checkbox" checked="true" value="true" id="select-no-risk">
                <label for="select-no-risk" class="safe-issues"> Safe</label>
              </p>
            </div>
            <a class="interval-button">
              <p class="change-interval">Change interval</p>
              <input type="number" value="1" id="graph-interval" min="1">
            </a>
          </div>
          <div class="interval-graph-container">
            <canvas id="interval-graph"></canvas>
          </div>
        </div>
      </div>
    </div>
    <div class="container-dashboard">
      <div class="dashboard-segment">
        <div class="data-segment-header">
          <p class="choose-issue-description"></p>
        </div>
        <a class="issue-button suggested-issue"><p>Suggested Issue</p></a>
        <a class="issue-button quick-fix"><p>Quick Fix</p></a>
        <a class="issue-button scan-now">Scan Now</a>
      </div>
      <div class="dashboard-segment risk-areas">
        <div class="data-segment-header">
          <p id="risk-areas">Areas of security risks</p>
        </div>
        <div class="security-area">
          <a>
            <p>
              <span class="applications">Applications</span>
              <span class="material-symbols-outlined">apps_outage</span>
            </p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="browser">Browser</span><span class="material-symbols-outlined">travel_explore</span></p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="devices">Devices</span><span class="material-symbols-outlined">devices</span></p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p>
              <span class="operating-system">Operating system</span>
              <span class="material-symbols-outlined">desktop_windows</span>
            </p>        
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="passwords">Passwords</span><span class="material-symbols-outlined">key</span></p>
          </a>
        </div>
        <div class="security-area">
          <a>
            <p><span class="other">Other</span><span class="material-symbols-outlined">view_cozy</span></p>
          </a>
        </div>
      </div>
    </div>
  </div>
  `;
  // Set counters on the page to the right values
  let rc = JSON.parse(sessionStorage.getItem('PrivacyRiskCounters'));
  adjustWithRiskCounters(rc, document);
  setMaxInterval(rc, document);

  // Localize the static content of the dashboard
  const staticDashboardContent = [
    'issues',
    'high-risk-issues',
    'medium-risk-issues',
    'low-risk-issues',
    'info-risk-issues',
    'safe-issues',
    'privacy-stat',
    'suggested-issue',
    'quick-fix',
    'scan-now',
    'applications',
    'browser',
    'devices',
    'operating-system',
    'passwords',
    'other',
    'select-risks',
    'change-interval',
    'choose-issue-description',
    'bar-graph-description',
  ];
  const localizationIds = [
    'Dashboard.Issues',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.InfoRisk',
    'Dashboard.Safe',
    'Dashboard.PrivacyStatus',
    'Dashboard.SuggestedIssue',
    'Dashboard.QuickFix',
    'Dashboard.ScanNow',
    'Dashboard.Applications',
    'Dashboard.Browser',
    'Dashboard.Devices',
    'Dashboard.OperatingSystem',
    'Dashboard.Passwords',
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
  new PieChart('pie-chart', rc, 'Privacy');
  const g = new Graph('interval-graph', rc);
  addGraphFunctions(g);
  document.getElementsByClassName('scan-now')[0].addEventListener('click', async () => {
    await scanTest();
    rc = JSON.parse(sessionStorage.getItem('PrivacyRiskCounters'));
    adjustWithRiskCounters(rc, document);
    setMaxInterval(rc, document);
    g.rc = rc;
    await g.changeGraph();
  });
}

document.getElementById('privacy-dashboard-button').addEventListener('click', () => openPrivacyDashboardPage());
