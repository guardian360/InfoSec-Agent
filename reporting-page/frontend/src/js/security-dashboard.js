import {Graph} from './graph.js';
import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';
import {suggestedIssue} from './home.js';
import {openAllChecksPage} from './all-checks.js';

/** Load the content of the Security Dashboard page */
export function openSecurityDashboardPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('security-dashboard-button');
  sessionStorage.setItem('savedPage', '2');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard">
    <div class="container-dashboard"> <!-- title top container -->
      <div class="dashboard-segment dashboard-title"> <!-- title top segment -->
        <p class="lang-security-dashboard"><p> 
      </div>
    </div>
    <div class="container-dashboard"> <!-- top container -->
      <div class="dashboard-segment"> <!-- Security status segment -->
        <div class="data-segment-header">
          <p class="lang-security-stat"></p>
        </div>
        <div class="security-status">
          <p class="status-descriptor"></p>
        </div>
      </div>
      <div class="dashboard-segment"> <!-- functional buttons segment -->
        <div class="data-segment-header">
          <p class="lang-choose-issue-description"></p>
        </div>
        <a id="suggested-issue" class="security-button lang-suggested-issue"><p></p></a>
        <a id="scan-now" class="security-button lang-scan-now"></a>
      </div>
      <div class="dashboard-segment risk-areas"> <!-- informative buttons segment -->
        <div class="data-segment-header">
          <p class="lang-security-risk-areas"></p>
        </div>
        <div class="security-area-buttons">
          <div class="security-area security-risk-button" id="security-button-applications">
            <a>
              <p>
                <span class="lang-applications"></span>
                <span class="material-symbols-outlined">apps_outage</span>
              </p>
            </a>
          </div>
          <div class="security-area security-risk-button" id="security-button-devices">
            <a>
              <p><span class="lang-devices"></span><span class="material-symbols-outlined">devices</span></p>
            </a>
          </div>
          <div class="security-area security-risk-button" id="security-button-network">
          <a>
            <p><span class="lang-network"></span><span class="material-symbols-outlined">lan</span></p>
          </a>
        </div>
          <div class="security-area security-risk-button" id="security-button-os">
            <a>
              <p>
                <span class="lang-operating-system"></span>
                <span class="material-symbols-outlined">desktop_windows</span>
              </p>        
            </a>
          </div>
          <div class="security-area security-risk-button" id="security-button-passwords">
            <a>
              <p><span class="lang-passwords"></span><span class="material-symbols-outlined">key</span></p>
            </a>
          </div>
          <div class="security-area security-risk-button" id="security-button-other">
            <a>
              <p><span class="lang-other"></span><span class="material-symbols-outlined">view_cozy</span></p>
            </a>
          </div>
        </div>
      </div>     
    </div>
    <div class="container-dashboard"> <!-- bottom container -->
      <div class="dashboard-segment"> <!-- Security risk counters segment -->
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
          <canvas class="pie-chart" id="pie-chart-security"></canvas>
        </div>
      </div>
      <div class="dashboard-segment"> <!-- graph segment -->
        <div class="data-segment-header">
          <p class="lang-risk-level-distribution"></p>
        </div>
        <div class="graph-segment-content">
          <div class="graph-buttons">
            <p class="lang-bar-graph-description">
            </p>
            <button id="dropbtn" class="dropbtn security-button"><span class="lang-select-risks"></span></button>
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
  let rc = JSON.parse(sessionStorage.getItem('SecurityRiskCounters'));
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
    'lang-security-stat',
    'lang-risk-level-counters',
    'lang-risk-level-distribution',
    'lang-suggested-issue',
    'lang-scan-now',
    'lang-security-risk-areas',
    'lang-applications',
    'lang-devices',
    'lang-network',
    'lang-operating-system',
    'lang-passwords',
    'lang-other',
    'lang-select-risks',
    'lang-change-interval',
    'lang-choose-issue-description',
    'lang-bar-graph-description',
  ];
  const localizationIds = [
    'Navigation.SecurityDashboard',
    'Dashboard.Issues',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.InfoRisk',
    'Dashboard.Acceptable',
    'Dashboard.SecurityStatus',
    'Dashboard.RiskLevelCounters',
    'Dashboard.RiskLevelDistribution',
    'Dashboard.SuggestedIssue',
    'Dashboard.ScanNow',
    'Dashboard.SecurityRiskAreas',
    'Dashboard.Applications',
    'Dashboard.Devices',
    'Dashboard.Network',
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
  new PieChart('pie-chart-security', rc, 'Security');
  const g = new Graph('interval-graph', rc);
  addGraphFunctions(g);
  document.getElementById('scan-now').addEventListener('click', async () => {
    await scanTest(true);
    rc = JSON.parse(sessionStorage.getItem('SecurityRiskCounters'));
    adjustWithRiskCounters(rc, document, true);
    setMaxInterval(rc, document);
    g.rc = rc;
    await g.changeGraph();
  });
  document.getElementById('suggested-issue').addEventListener('click', () => suggestedIssue('Security'));

  // Add links to checks page
  document.getElementById('security-button-applications').addEventListener('click',
    () => openAllChecksPage('applications'));
  document.getElementById('security-button-devices').addEventListener('click',
    () => openAllChecksPage('devices'));
  document.getElementById('security-button-network').addEventListener('click',
    () => openAllChecksPage('network'));
  document.getElementById('security-button-os').addEventListener('click',
    () => openAllChecksPage('os'));
  document.getElementById('security-button-passwords').addEventListener('click',
    () => openAllChecksPage('passwords'));
  document.getElementById('security-button-other').addEventListener('click',
    () => openAllChecksPage('security-other'));
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('security-dashboard-button').addEventListener('click', () => openSecurityDashboardPage());
  } catch (error) {
    logError('Error in security-dashboard.js: ' + error);
  }
}

/** Changes the risk counters to show the correct values
 *
 * @param {RiskCounters} rc Risk counters from which the data is taken
 * @param {Document} doc Document in which the counters are located
 * @param {boolean} retrieveStyling Boolean to determine if the colors of the risk levels should be retrieved
 */
export function adjustWithRiskCounters(rc, doc, retrieveStyling) {
  // change counters according to collected data
  doc.getElementById('high-risk-counter').innerHTML = rc.lastHighRisk;
  doc.getElementById('medium-risk-counter').innerHTML = rc.lastMediumRisk;
  doc.getElementById('low-risk-counter').innerHTML = rc.lastLowRisk;
  doc.getElementById('info-risk-counter').innerHTML = rc.lastInfoRisk;
  doc.getElementById('no-risk-counter').innerHTML = rc.lastNoRisk;

  if (retrieveStyling) {
    rc.highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--high-risk-color');
    rc.mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--medium-risk-color');
    rc.lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--low-risk-color');
    rc.infoColor = getComputedStyle(document.documentElement).getPropertyValue('--info-color');
    rc.noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--no-risk-color');
  }

  const securityStatus = doc.getElementsByClassName('status-descriptor')[0];
  if (rc.lastHighRisk > 1) {
    try {
      getLocalization('Dashboard.Critical', 'status-descriptor');
    } catch (error) {
      /* istanbul ignore next */
      securityStatus.innerHTML = 'Critical';
    }
    securityStatus.style.backgroundColor = rc.highRiskColor;
  } else if (rc.lastMediumRisk > 1) {
    try {
      getLocalization('Dashboard.MediumConcern', 'status-descriptor');
    } catch (error) {
      /* istanbul ignore next */
      securityStatus.innerHTML = 'Medium concern';
    }
    securityStatus.style.backgroundColor = rc.mediumRiskColor;
  } else if (rc.lastLowRisk > 1) {
    try {
      getLocalization('Dashboard.LowConcern', 'status-descriptor');
    } catch (error) {
      /* istanbul ignore next */
      securityStatus.innerHTML = 'Low concern';
    }
    securityStatus.style.backgroundColor = rc.lowRiskColor;
  } else if (rc.lastInfoRisk > 1) {
    try {
      getLocalization('Dashboard.InfoConcern', 'status-descriptor');
    } catch (error) {
      /* istanbul ignore next */
      securityStatus.innerHTML = 'Informative';
    }
    securityStatus.style.backgroundColor = rc.infoColor;
  } else {
    try {
      getLocalization('Dashboard.NoConcern', 'status-descriptor');
    } catch (error) {
      /* istanbul ignore next */
      securityStatus.innerHTML = 'Acceptable';
    }
    securityStatus.style.backgroundColor = rc.noRiskColor;
  }
  securityStatus.style.color = 'rgb(255, 255, 255)';
}

/** Set the max number input of the 'graph-interval' element
 *
 * @param {RiskCounters} rc Risk counters from which the max count is taken
 * @param {Document} doc Document in which the counters are located
 */
export function setMaxInterval(rc, doc) {
  doc.getElementById('graph-interval').max = rc.count;
}

/** Adds event listeners to elements in graph-row section of the dashboard page
 *
 * @param {Graph} g Graph class containing the functions to be called
 */
export function addGraphFunctions(g) {
  document.getElementById('dropbtn').addEventListener('click', () => g.graphDropdown());
  document.getElementById('graph-interval').addEventListener('change', () => g.changeGraph());
  document.getElementById('select-high-risk').addEventListener('change', () => g.toggleRisks('high'));
  document.getElementById('select-medium-risk').addEventListener('change', () => g.toggleRisks('medium'));
  document.getElementById('select-low-risk').addEventListener('change', () => g.toggleRisks('low'));
  document.getElementById('select-info-risk').addEventListener('change', () => g.toggleRisks('info'));
  document.getElementById('select-no-risk').addEventListener('change', () => g.toggleRisks('no'));
}
