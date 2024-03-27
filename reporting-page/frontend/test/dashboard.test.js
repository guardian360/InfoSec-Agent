import test from 'unit.js';
import { JSDOM } from "jsdom";
import { AdjustWithRiskCounters, SetMaxInterval } from '../src/js/security-dashboard.js';

// Mock page
const dom = new JSDOM(`
  <div class="dashboard-data">
  <div class="data-column risk-counters">
    <div class="security-status">
      <div><p>Security status</p></div>
      <div><p id="security-status">Critical</p></div>
    </div>
    <div class="risk-counter high-risk">
      <div><p>High risk issues</p></div>
      <div><p id="high-risk-counter">0</p></div>
    </div>
    <div class="risk-counter medium-risk">
      <div><p>Medium risk issues</p></div>
      <div><p id="medium-risk-counter">0</p></div>
    </div>
    <div class="risk-counter low-risk">
      <div><p>Low risk issues</p></div>
      <div><p id="low-risk-counter">0</p></div>
    </div>
    <div class="risk-counter no-risk">
      <div><p>Safe issues</p></div>
      <div><p id="no-risk-counter">0</p></div>
    </div>
  </div>
  <div class="data-column piechart">
    <canvas id="pieChart"></canvas>
  </div>
  <div class="data-column issue-buttons">
    <H2>You have some issues you can fix. 
      To start resolving an issue either navigate to the issues page, or pick a suggested issue below
    </H2>
    <a class="issue-button"><p>Suggested Issue</p></a>
    <a class="issue-button"><p>Quick Fix</p></a>
  </div>
  </div>
  <div class="second-row">
  <h2>Areas of security/privacy risks</h2>
  <div class="security-areas">
    <div class="security-area">
      <a>
        <p><span class="material-symbols-outlined">apps_outage</span><span>Applications</span></p>
      </a>
      <a class="areas-issues-button">
        <p>Issues</p>
      </a>
    </div>
    <div class="security-area">
      <a>
        <p><span class="material-symbols-outlined">travel_explore</span><span>Browser</span></p>
      </a>
      <a class="areas-issues-button">
        <p>Issues</p>
      </a>
    </div>
    <div class="security-area">
      <a>
        <p><span class="material-symbols-outlined">devices</span><span>Devices</span></p>
      </a>
      <a class="areas-issues-button">
        <p>Issues</p>
      </a>
    </div>
    <div class="security-area">
      <a>
        <p><span class="material-symbols-outlined">desktop_windows</span><span>Operating system</span></p>
      </a>
      <a class="areas-issues-button">
        <p>Issues</p>
      </a>
    </div>
    <div class="security-area">
      <a>
        <p><span class="material-symbols-outlined">key</span><span>Passwords</span></p>
      </a>
      <a class="areas-issues-button">
        <p>Issues</p>
      </a>
    </div>
    <div class="security-area">
      <a>
        <p><span class="material-symbols-outlined">view_cozy</span><span>Other</span></p>
      </a>
      <a class="areas-issues-button">
        <p>Issues</p>
      </a>
    </div>
  </div>
  </div>
  <div class="graph-row">
  <div class="graph-column issues-graph-buttons">
    <H2>In this graph you are able to see the distribution of different issues we have found over the past 5 times we ran a check.</H2>
    <div class="dropdown">
      <button class="dropbtn" id="dropbtn">Select Risks</button>
      <div class="dropdown-selector" id="myDropdown">
        <p><input type="checkbox" checked="true" value="true" id="select-high-risk">
          <label for="select-high-risk"> High risks</label><br>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-medium-risk">
          <label for="select-medium-risk"> Medium risks</label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-low-risk">
          <label for="select-low-risk"> Low risks</label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-no-risk">
          <label for="select-no-risk"> Safe</label>
        </p>
      </div>
    </div>
    <a class="interval-button"><p>Change interval</p><input type="number" value="5" id="graph-interval" min=1></a>
  </div>
  <div class="graph-column issues-graph">
    <canvas id="interval-graph"></canvas>
  </div>
  </div>
`, {
url: 'http://localhost'
});
global.document = dom.window.document
global.window = dom.window

// test cases
describe("dashboard", function() {
    it("Should show data from risk counters", function() {
      const mockRiskCounters = {  
        lastHighRisk : 2,
        lastMediumRisk : 3,
        lastLowRisk : 4,
        lastnoRisk : 5,

        count : 5,
      };

      // act
      AdjustWithRiskCounters(mockRiskCounters);

      // assert
      test.value(document.getElementById("high-risk-counter").innerHTML).isEqualTo(mockRiskCounters.lastHighRisk);
      test.value(document.getElementById("medium-risk-counter").innerHTML).isEqualTo(mockRiskCounters.lastMediumRisk);
      test.value(document.getElementById("low-risk-counter").innerHTML).isEqualTo(mockRiskCounters.lastLowRisk);
      test.value(document.getElementById("no-risk-counter").innerHTML).isEqualTo(mockRiskCounters.lastnoRisk);

    })
    it("Should display the right security status", function() {
      // arrange
      const expectedColors = ["rgb(255, 255, 255)","rgb(255, 255, 255)","rgb(0, 0, 0)","rgb(0, 0, 0)"]
      const expectedBackgroundColors = ["rgb(0, 255, 255)","rgb(0, 0, 255)","rgb(255, 0, 0)","rgb(255, 255, 0)"]
      const expectedText = ["Critical","Medium concern","Light concern","Safe"]

      const mockRiskCounters = {
        highRiskColor : "rgb(0, 255, 255)",
        mediumRiskColor : "rgb(0, 0, 255)",
        lowRiskColor : "rgb(255, 0, 0)",
        noRiskColor : "rgb(255, 255, 0)",
  
        lastHighRisk : 10,
        lastMediumRisk : 10,
        lastLowRisk : 10,
        lastnoRisk : 10,
      };
      expectedColors.forEach((element,index) => {
        // act
        if (index == 1) mockRiskCounters.lastHighRisk = 0;
        if (index == 2) mockRiskCounters.lastMediumRisk = 0;
        if (index == 3) mockRiskCounters.lastLowRisk = 0;
        AdjustWithRiskCounters(mockRiskCounters);

        // assert
        test.value(document.getElementById("security-status").innerHTML).isEqualTo(expectedText[index]);
        test.value(document.getElementById("security-status").style.backgroundColor).isEqualTo(expectedBackgroundColors[index]);
        test.value(document.getElementById("security-status").style.color).isEqualTo(expectedColors[index]);
      });
    })
    it("Should set the max value of the graph interval to the maximum amount of data", function() {
      // arrange
      const mockRiskCounters = {  
        count : 5,
      };

      // act
      SetMaxInterval(mockRiskCounters);

      // assert
      test.value(document.getElementById("graph-interval").max).isEqualTo(mockRiskCounters.count);
    })
})