import * as rc from "./risk-counters";
import * as graph from "./graph";
import * as piechart from "./piechart";
import { GetLocalization } from './localize.js';

/** Load the content of the Security Dashboard page */
function openSecurityDashboardPage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="dashboard-data">
    <div class="data-column risk-counters">
      <div class="security-status">
        <div><p class="security-stat">Security status</p></div>
        <div><p class="status-descriptor"></p></div>
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
    <div class="data-column piechart">
      <canvas id="pieChart"></canvas>
    </div>
    <div class="data-column issue-buttons">
      <H2 class="choose-issue-description">You have some issues you can fix. 
        To start resolving an issue either navigate to the issues page, or pick a suggested issue below.
      </H2>
      <a class="issue-button suggested-issue"><p>Suggested Issue</p></a>
      <a class="issue-button quick-fix"><p>Quick Fix</p></a>
    </div>
  </div>
  <div class="second-row">
    <h2 id="risk-areas">Areas of security/privacy risks</h2>
    <div class="security-areas">
      <div class="security-area">
        <a>
          <p><span class="material-symbols-outlined">apps_outage</span><span class="applications">Applications</span></p>
        </a>
        <a class="areas-issues-button">
          <p class="issues">Issues</p>
        </a>
      </div>
      <div class="security-area">
        <a>
          <p><span class="material-symbols-outlined">travel_explore</span><span class="browser">Browser</span></p>
        </a>
        <a class="areas-issues-button">
          <p class="issues">Issues</p>
        </a>
      </div>
      <div class="security-area">
        <a>
          <p><span class="material-symbols-outlined">devices</span><span class="devices">Devices</span></p>
        </a>
        <a class="areas-issues-button">
          <p class="issues">Issues</p>
        </a>
      </div>
      <div class="security-area">
        <a>
          <p><span class="material-symbols-outlined">desktop_windows</span><span class="operating-system">Operating system</span></p>
        </a>
        <a class="areas-issues-button">
          <p class="issues">Issues</p>
        </a>
      </div>
      <div class="security-area">
        <a>
          <p><span class="material-symbols-outlined">key</span><span class="passwords">Passwords</span></p>
        </a>
        <a class="areas-issues-button">
          <p class="issues">Issues</p>
        </a>
      </div>
      <div class="security-area">
        <a>
          <p><span class="material-symbols-outlined">view_cozy</span><span class="other">Other</span></p>
        </a>
        <a class="areas-issues-button">
          <p class="issues">Issues</p>
        </a>
      </div>
    </div>
  </div>
  <div class="graph-row">
    <div class="graph-column issues-graph-buttons">
      <H2 class="bar-graph-description">In this graph you are able to see the distribution of different issues we have found over the past 5 times we ran a check.</H2>
      <div class="dropdown">
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
      </div>
      <a class="interval-button"><p class="change-interval">Change interval</p><input type="number" value="5" id="graph-interval" min=1></a>
    </div>
    <div class="graph-column issues-graph">
      <canvas id="interval-graph"></canvas>
    </div>
  </div>
  `;  
  // Set counters on the page to the right values
  AdjustWithRiskCounters();  
  // Add functionalities to dashboard
  AddGraphFunctions();  
  // Create charts
  CreatePieChart();
  CreateGraphChart();

  // Localize the static content of the dashboard
  let staticDashboardContent = [
    "issues",
    "high-risk-issues", 
    "medium-risk-issues",
    "low-risk-issues",
    "safe-issues",
    "security-stat",
    "suggested-issue", 
    "quick-fix", 
    "applications",
    "browser",
    "devices",
    "operating-system",
    "passwords",
    "other",
    "select-risks",
    "change-interval",
    "choose-issue-description",
    "bar-graph-description"
  ]
  let localizationIds = [
    "Dashboard.Issues",
    "Dashboard.HighRisk", 
    "Dashboard.MediumRisk",
    "Dashboard.LowRisk",
    "Dashboard.Safe",
    "Dashboard.SecurityStatus",
    "Dashboard.SuggestedIssue",
    "Dashboard.QuickFix",
    "Dashboard.Applications",
    "Dashboard.Browser",
    "Dashboard.Devices",
    "Dashboard.OperatingSystem",
    "Dashboard.Passwords",
    "Dashboard.Other",
    "Dashboard.SelectRisks",
    "Dashboard.ChangeInterval",
    "Dashboard.ChooseIssueDescription",
    "Dashboard.BarGraphDescription"
  ]
  for (let i = 0; i < staticDashboardContent.length; i++) {
    GetLocalization(localizationIds[i], staticDashboardContent[i])
  }
}

document.getElementById("security-dashboard-button").addEventListener("click", () => openSecurityDashboardPage());

/** Changes the risk counters to show the correct values */
function AdjustWithRiskCounters() {
  // change counters according to collected data
  document.getElementById("high-risk-counter").innerHTML = rc.lastHighRisk;
  document.getElementById("medium-risk-counter").innerHTML = rc.lastMediumRisk;
  document.getElementById("low-risk-counter").innerHTML = rc.lastLowRisk;
  document.getElementById("no-risk-counter").innerHTML = rc.lastnoRisk; 

  let securityStatus = document.getElementsByClassName("status-descriptor")[0];  
  if (rc.lastHighRisk > 1) {
    GetLocalization("Dashboard.Critical", "status-descriptor");
    // securityStatus.innerHTML = "Critical";
    securityStatus.style.backgroundColor = rc.highRiskColor;
    securityStatus.style.color = "rgb(255, 255, 255)";
  } else if (rc.lastMediumRisk > 1) {
    GetLocalization("Dashboard.MediumConcern", "status-descriptor");
    // securityStatus.innerHTML = "Medium concern";
    securityStatus.style.backgroundColor = rc.mediumRiskColor; 
    securityStatus.style.color = "rgb(255, 255, 255)";
  } else if (rc.lastLowRisk > 1) {
    GetLocalization("Dashboard.LightConcern", "status-descriptor");
    // securityStatus.innerHTML = "Light concern";
    securityStatus.style.backgroundColor = rc.lowRiskColor;
    securityStatus.style.color = "rgb(0, 0, 0)";  
  } else {
    GetLocalization("Dashboard.NoConcern", "status-descriptor");
    // securityStatus.innerHTML = "Safe";
    securityStatus.style.backgroundColor = rc.noRiskColor;
    securityStatus.style.color = "rgb(0, 0, 0)";  
  }  

  document.getElementById("graph-interval").max = rc.allNoRisks.length;
}

/** Adds eventlisteners to elements in graph-row section of the dashboard page */
function AddGraphFunctions() {
  document.getElementById("dropbtn").addEventListener("click", () => GraphDropdown());
  document.getElementById("graph-interval").addEventListener("change", () => ChangeGraph());
  document.getElementById("select-high-risk").addEventListener("change", () => ToggleRisks("high"));
  document.getElementById("select-medium-risk").addEventListener("change", () => ToggleRisks("medium"));
  document.getElementById("select-low-risk").addEventListener("change", () => ToggleRisks("low"));
  document.getElementById("select-no-risk").addEventListener("change", () => ToggleRisks("no"));
}

//#region PieChart

// Reusable snippit for other files
let pieChart;

/** Creates a pie chart for risks */
function CreatePieChart() {
  pieChart = new Chart("pieChart", {
    type: "doughnut",
    data: piechart.GetData(),
    options: piechart.GetOptions()
  });
}

//#endregion

//#region Graph

// Function to change the graph don't work when imported from another file.
// Piechart now resides here.

let graphShowHighRisks = true;
let graphShowMediumRisks = true;
let graphShowLowRisks = true;
let graphShowNoRisks = true;

let graphShowAmount = 5;

let barChart;

/** Creates a graph in the form of a bar chart for risks */
function CreateGraphChart() {
  barChart = new Chart("interval-graph", {
    type: 'bar',
    data: graph.GetData(graphShowAmount, graphShowHighRisks, graphShowMediumRisks, graphShowLowRisks, graphShowNoRisks), // The data for our dataset
    options: graph.GetOptions() // Configuration options go here
  });
}

/** Updates the graph, should be called after a change in graph properties */
function ChangeGraph() {
  graphShowAmount = document.getElementById('graph-interval').value;
  barChart.data = graph.GetData(graphShowAmount, graphShowHighRisks, graphShowMediumRisks, graphShowLowRisks, graphShowNoRisks);
  console.log(graphShowAmount);
  barChart.update();
}

/** Toggles a risks to show in the graph 
 * 
 * @param {string} category Category corresponding to risk 
 */
function ToggleRisks(category) {
  switch (category) {
    case "high":
      graphShowHighRisks = !graphShowHighRisks;
      break;
    case "medium":
      graphShowMediumRisks = !graphShowMediumRisks;
      break;
    case "low":
      graphShowLowRisks = !graphShowLowRisks;
      break;
    case "no":
      graphShowNoRisks = !graphShowNoRisks;
      break;
    default:
      break;
  }
  ChangeGraph();
}

/** toggles 'show' class on element with id:"myDropDown" */
function GraphDropdown() {
  document.getElementById("myDropdown").classList.toggle("show");
}

//#endregion 
