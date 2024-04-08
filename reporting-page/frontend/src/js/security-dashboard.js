import {Graph} from './graph.js';
import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {ScanNow as scanNowGo} from '../../wailsjs/go/main/Tray.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';

/** Load the content of the Security Dashboard page */
function openSecurityDashboardPage() {
  closeNavigation();
  markSelectedNavigationItem('security-dashboard-button');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard-data">
    <div class="data-column risk-analysis">
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="security-stat">Security status</p>
        </div>
        <div class="security-status">
          <p class="status-descriptor"></p>
        </div>
      </div>
      <div class="data-segment">
        <div class="data-segment-header">
          <p class="risk-counters-header">Risk level counters</p>
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
        <div class="risk-counter no-risk">
          <div><p class="safe-issues">Safe issues</p></div>
          <div><p id="no-risk-counter">0</p></div>
        </div>
      </div>
    </div>
    <div class="data-column">
      <div class="data-segment piechart">
        <div class="data-segment-header">
            <p class="piechart-header">Risk level distribution</p>
        </div>
        <div class="piechart-container">
          <canvas id="pieChart"></canvas>
        </div>
      </div>
      <div class="data-segment graph-row">
        <div class="data-segment-header">
          <p class="bar-graph-header">Risk level distribution</p>
        </div>
        <div class="graph-segment-content">
          <div class="graph-buttons dropdown">
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
              <p><input type="checkbox" checked="true" value="true" id="select-no-risk">
                <label for="select-no-risk" class="safe-issues"> Safe</label>
              </p>
            </div>
            <a class="interval-button">
              <p class="change-interval">Change interval</p>
              <input type="number" value="1" id="graph-interval" min="1">
            </a>
          </div>
          <div class="graph-column issues-graph">
            <canvas id="interval-graph"></canvas>
          </div>
        </div>
      </div>
    </div>
    <div class="data-column actions">
      <div class="data-segment issue-buttons">
        <div class="data-segment-header">
          <p class="choose-issue-description"></p>
        </div>
        <a class="issue-button suggested-issue"><p>Suggested Issue</p></a>
        <a class="issue-button quick-fix"><p>Quick Fix</p></a>
        <a class="issue-button scan-now">Scan Now</a>
      </div>
      <div class="data-segment risk-areas">
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
  const rc = JSON.parse(sessionStorage.getItem('RiskCounters'));
  console.log(rc);
  adjustWithRiskCounters(rc, document);
  setMaxInterval(rc, document);
  // Create charts

  // Localize the static content of the dashboard
  const staticDashboardContent = [
    'issues',
    'high-risk-issues',
    'medium-risk-issues',
    'low-risk-issues',
    'safe-issues',
    'security-stat',
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
    'Dashboard.Safe',
    'Dashboard.SecurityStatus',
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
  new PieChart('pieChart', rc);
  const g = new Graph('interval-graph', rc);
  addGraphFunctions(g);
  document.getElementsByClassName('scan-now')[0].addEventListener('click', () => {
    scanNowGo();
    g.changeGraph();
  });
  document.onload = retrieveTheme();
}

if (typeof document !== 'undefined') {
  try {
    document.getElementById('security-dashboard-button').addEventListener('click', () => openSecurityDashboardPage());
  } catch (error) {
    console.log('Error in security-dashboard.js: ' + error);
  }
}

/** Changes the risk counters to show the correct values
 *
 * @param {RiskCounters} rc Risk counters from which the data is taken
 * @param {Document} doc Document in which the counters are located
 */
export function adjustWithRiskCounters(rc, doc) {
  // change counters according to collected data
  doc.getElementById('high-risk-counter').innerHTML = rc.lastHighRisk;
  doc.getElementById('medium-risk-counter').innerHTML = rc.lastMediumRisk;
  doc.getElementById('low-risk-counter').innerHTML = rc.lastLowRisk;
  doc.getElementById('no-risk-counter').innerHTML = rc.lastnoRisk;

  const securityStatus = doc.getElementsByClassName('status-descriptor')[0];
  if (rc.lastHighRisk > 1) {
    try {
      getLocalization('Dashboard.Critical', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Critical';
    }
    securityStatus.style.backgroundColor = rc.highRiskColor;
    securityStatus.style.color = 'rgb(255, 255, 255)';
  } else if (rc.lastMediumRisk > 1) {
    try {
      getLocalization('Dashboard.MediumConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Medium concern';
    }
    securityStatus.style.backgroundColor = rc.mediumRiskColor;
    securityStatus.style.color = 'rgb(255, 255, 255)';
  } else if (rc.lastLowRisk > 1) {
    try {
      getLocalization('Dashboard.LightConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Light concern';
    }
    securityStatus.style.backgroundColor = rc.lowRiskColor;
    securityStatus.style.color = 'rgb(0, 0, 0)';
  } else {
    try {
      getLocalization('Dashboard.NoConcern', 'status-descriptor');
    } catch (error) {
      securityStatus.innerHTML = 'Safe';
    }
    securityStatus.style.backgroundColor = rc.noRiskColor;
    securityStatus.style.color = 'rgb(0, 0, 0)';
  }
}

/** Set the max number input of the 'graph-interval' element
 *
 * @param {RiskCounters} rc Risk counters from which the max count is taken
 * @param {Document} doc Document in which the counters are located
 */
export function setMaxInterval(rc, doc) {
  doc.getElementById('graph-interval').max = rc.count;
}

/** Adds eventlisteners to elements in graph-row section of the dashboard page
 *
 * @param {Graph} g Graph class containing the functions to be called
 */
export function addGraphFunctions(g) {
  document.getElementById('dropbtn').addEventListener('click', () => g.graphDropdown());
  document.getElementById('graph-interval').addEventListener('change', () => g.changeGraph());
  document.getElementById('select-high-risk').addEventListener('change', () => g.toggleRisks('high'));
  document.getElementById('select-medium-risk').addEventListener('change', () => g.toggleRisks('medium'));
  document.getElementById('select-low-risk').addEventListener('change', () => g.toggleRisks('low'));
  document.getElementById('select-no-risk').addEventListener('change', () => g.toggleRisks('no'));
}

