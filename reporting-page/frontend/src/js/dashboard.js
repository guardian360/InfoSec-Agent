import {RiskCounters} from "./risk-counters.js";
import {Graph} from "./graph.js";
import {PieChart} from "./piechart.js";

/** Load the content of the Dashboard page */
export function openDashboardPage() {
  document.getElementById("page-contents").innerHTML = `
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
  `;  
  // Set counters on the page to the right values
  let rc = new RiskCounters();
  AdjustWithRiskCounters(rc);
  SetMaxInterval(rc);    
  // Create charts
  new PieChart("pieChart",rc);
  let g = new Graph("interval-graph",rc);
  AddGraphFunctions(g);
}

if (typeof document !== 'undefined') {
  document.getElementById("dashboard-button").addEventListener("click", () => openDashboardPage());
}

/** Changes the risk counters to show the correct values 
 * 
 * @param {RiskCounters} rc Risk counters from which the data is taken 
 */
export function AdjustWithRiskCounters(rc) {
  // change counters according to collected data
  document.getElementById("high-risk-counter").innerHTML = rc.lastHighRisk;
  document.getElementById("medium-risk-counter").innerHTML = rc.lastMediumRisk;
  document.getElementById("low-risk-counter").innerHTML = rc.lastLowRisk;
  document.getElementById("no-risk-counter").innerHTML = rc.lastnoRisk; 

  let securityStatus = document.getElementById("security-status");  
  if (rc.lastHighRisk > 1) {
    securityStatus.innerHTML = "Critical";
    securityStatus.style.backgroundColor = rc.highRiskColor;
    securityStatus.style.color = "rgb(255, 255, 255)";
  } else if (rc.lastMediumRisk > 1) {
    securityStatus.innerHTML = "Medium concern";
    securityStatus.style.backgroundColor = rc.mediumRiskColor; 
    securityStatus.style.color = "rgb(255, 255, 255)";
  } else if (rc.lastLowRisk > 1) {
    securityStatus.innerHTML = "Light concern";
    securityStatus.style.backgroundColor = rc.lowRiskColor;
    securityStatus.style.color = "rgb(0, 0, 0)";  
  } else {
    securityStatus.innerHTML = "Safe";
    securityStatus.style.backgroundColor = rc.noRiskColor;
    securityStatus.style.color = "rgb(0, 0, 0)";  
  }  
}

/** Set the max number input of the 'graph-interval' element
 * 
 * @param {RiskCounters} rc Risk counters from which the max count is taken
 */
export function SetMaxInterval(rc) {
  document.getElementById("graph-interval").max = rc.count;
}

/** Adds eventlisteners to elements in graph-row section of the dashboard page 
 * 
 * @param {Graph} g Graph class containing the functions to be called
 */
export function AddGraphFunctions(g) {
  document.getElementById("dropbtn").addEventListener("click", () => g.GraphDropdown());
  document.getElementById("graph-interval").addEventListener("change", () => g.ChangeGraph());
  document.getElementById("select-high-risk").addEventListener("change", () => g.ToggleRisks("high"));
  document.getElementById("select-medium-risk").addEventListener("change", () => g.ToggleRisks("medium"));
  document.getElementById("select-low-risk").addEventListener("change", () => g.ToggleRisks("low"));
  document.getElementById("select-no-risk").addEventListener("change", () => g.ToggleRisks("no"));
}

