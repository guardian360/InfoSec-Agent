import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals'
// import {adjustWithRiskCounters, setMaxInterval} from '../src/js/security-dashboard.js';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
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
          <p class="bar-graph-description"></p>
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
            <input type="number" value="5" id="graph-interval" min=1>
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
`, {
  url: 'http://localhost',
});
global.document = dom.window.document
global.window = dom.window

// Mock scanTest
jest.unstable_mockModule('../src/js/database.js', () => ({
  scanTest: jest.fn()
}))

// Mock LogError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn()
}))

// test cases
describe('Security dashboard', function() {
  it('adjustWithRiskCounters should show data from risk counters', async function() {
    // arrange
    const mockRiskCounters = {
      lastHighRisk: 2,
      lastMediumRisk: 3,
      lastLowRisk: 4,
      lastnoRisk: 5,
      count: 5,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    // act
    dashboard.adjustWithRiskCounters(mockRiskCounters, global.document);

    // assert
    test.value(document.getElementById('high-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastHighRisk);
    test.value(document.getElementById('medium-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastMediumRisk);
    test.value(document.getElementById('low-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastLowRisk);
    test.value(document.getElementById('no-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastnoRisk);
  });
  it('Should display the right security status', async function() {
    // arrange
    const expectedColors = ['rgb(255, 255, 255)', 'rgb(255, 255, 255)', 'rgb(0, 0, 0)', 'rgb(0, 0, 0)'];
    // const expectedBackgroundColors = ['rgb(0, 255, 255)', 'rgb(0, 0, 255)', 'rgb(255, 0, 0)', 'rgb(255, 255, 0)'];
    // const expectedText = ['Critical', 'Medium concern', 'Light concern', 'Safe'];

    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',

      lastHighRisk: 10,
      lastMediumRisk: 10,
      lastLowRisk: 10,
      lastnoRisk: 10,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    expectedColors.forEach((element, index) => {
      // act
      dashboard.adjustWithRiskCounters(mockRiskCounters, dom.window.document);

      // assert
      test.value(dom.window.document.getElementById('high-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastHighRisk);
      test.value(dom.window.document.getElementById('medium-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastMediumRisk);
      test.value(dom.window.document.getElementById('low-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastLowRisk);
      test.value(dom.window.document.getElementById('no-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastnoRisk);
    });
  });
  it('adjustWithRiskCounters should display the right security status', async function() {
    // arrange
    const expectedColors = ['rgb(255, 255, 255)', 'rgb(255, 255, 255)', 'rgb(0, 0, 0)', 'rgb(0, 0, 0)'];
    const expectedBackgroundColors = ['rgb(0, 255, 255)', 'rgb(0, 0, 255)', 'rgb(255, 0, 0)', 'rgb(255, 255, 0)'];
    const expectedText = ['Critical', 'Medium concern', 'Light concern', 'Safe'];
    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',

      lastHighRisk: 10,
      lastMediumRisk: 10,
      lastLowRisk: 10,
      lastnoRisk: 10,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    expectedColors.forEach((element, index) => {
      // act
      if (index == 1) mockRiskCounters.lastHighRisk = 0;
      if (index == 2) mockRiskCounters.lastMediumRisk = 0;
      if (index == 3) mockRiskCounters.lastLowRisk = 0;
      dashboard.adjustWithRiskCounters(mockRiskCounters, dom.window.document);
      // assert
      test.value(dom.window.document.getElementsByClassName('status-descriptor')[0].innerHTML)
        .isEqualTo(expectedText[index]);
      test.value(dom.window.document.getElementsByClassName('status-descriptor')[0].style.backgroundColor)
        .isEqualTo(expectedBackgroundColors[index]);
      test.value(dom.window.document.getElementsByClassName('status-descriptor')[0].style.color)
        .isEqualTo(expectedColors[index]);
    });
  });
  it('setMaxInterval should set the max value of the graph interval to the maximum amount of data', async function() {
    // arrange
    const mockRiskCounters = {
      count: 5,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    // act
    dashboard.setMaxInterval(mockRiskCounters, dom.window.document);

    // assert
    test.value(dom.window.document.getElementById('graph-interval').max).isEqualTo(mockRiskCounters.count);
  });
});
