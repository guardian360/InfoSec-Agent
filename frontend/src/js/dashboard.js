import * as rc from "./risk-counters";
import * as graph from "./graph";
import * as piechart from "./piechart";

function openDashboardPage() {
    document.getElementById("page-contents").innerHTML = `
    <div class="dashboard-data">
            <div class="data-column risk-counters">
                <div class="security-status">
                    <div><p>Security status</p></div>
                    <div><p id="security-status">Critical</p></div>
                </div>
                <div class="risk-counter high-risk">
                    <div><p>High risk issues</p></div>
                    <div><p id="high-risk-counter">4</p></div>
                </div>
                <div class="risk-counter medium-risk">
                    <div><p>Medium risk issues</p></div>
                    <div><p id="medium-risk-counter">4</p></div>
                </div>
                <div class="risk-counter low-risk">
                    <div><p>Low risk issues</p></div>
                    <div><p id="low-risk-counter">4</p></div>
                </div>
                <div class="risk-counter no-risk">
                    <div><p>Safe issues</p></div>
                    <div><p id="no-risk-counter">4</p></div>
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
                <H2>In this graph you are able to see the distribution of different issues we have found over the past 5 times we ran a check.
                </H2>
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
      </div>
      
    `;

    // Set counters on the page to the right values
    AdjustWithRiskCounters();

    AddGraphFunctions();

    // Create charts
    CreatePieChart();
    CreateGraphChart();
}

document.getElementById("dashboard-button").addEventListener("click", () => openDashboardPage());

function AdjustWithRiskCounters() {
    // change counters according to collected data
    document.getElementById("high-risk-counter").innerHTML = rc.lastHighRisk;
    document.getElementById("medium-risk-counter").innerHTML = rc.lastMediumRisk;
    document.getElementById("low-risk-counter").innerHTML = rc.lastLowRisk;
    document.getElementById("no-risk-counter").innerHTML = rc.lastnoRisk;

    var securityStatus = document.getElementById("security-status");

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

    document.getElementById("graph-interval").max = rc.allNoRisks.length;
}

function AddGraphFunctions() {
    document.getElementById("dropbtn").addEventListener("click", () => GraphDropdown());
    document.getElementById("graph-interval").addEventListener("change", () => ChangeGraph())
    document.getElementById("select-high-risk").addEventListener("change", () => ToggleHighRisks())
    document.getElementById("select-medium-risk").addEventListener("change", () => ToggleMediumRisks())
    document.getElementById("select-low-risk").addEventListener("change", () => ToggleLowRisks())
    document.getElementById("select-no-risk").addEventListener("change", () => ToggleNoRisks())
}

//#region PieChart

// Reusable snippit for other files
let pieChart;
function CreatePieChart() {
    pieChart = new Chart("pieChart", {
        type: "doughnut",
        data: piechart.GetData()
        ,
        options: piechart.GetOptions()
      });
}

//#endregion

//#region Graph

// Function to change the graph don't work when imported from another file.
// Piechart now resides here.

var graphShowHighRisks = true;
var graphShowMediumRisks = true;
var graphShowLowRisks = true;
var graphShowNoRisks = true;

var graphShowAmount = 5;

let barChart;
function CreateGraphChart() {
    barChart = new Chart("interval-graph", {
        type: 'bar',
        // The data for our dataset
        data: graph.GetData(graphShowAmount, graphShowHighRisks, graphShowMediumRisks, graphShowLowRisks, graphShowNoRisks),
    
        // Configuration options go here
        options: graph.GetOptions()
    });
}


// Function to change the graph
function ChangeGraph() {
    graphShowAmount = document.getElementById('graph-interval').value;
    barChart.data = graph.GetData(graphShowAmount, graphShowHighRisks, graphShowMediumRisks, graphShowLowRisks, graphShowNoRisks);
    console.log(graphShowAmount);
    barChart.update();
}

function ToggleHighRisks() {
    graphShowHighRisks = !graphShowHighRisks;
    ChangeGraph();
}

function ToggleMediumRisks() {
    graphShowMediumRisks = !graphShowMediumRisks;
    ChangeGraph();
}

function ToggleLowRisks() {
    graphShowLowRisks = !graphShowLowRisks;
    ChangeGraph();
}

function ToggleNoRisks() {
    graphShowNoRisks = !graphShowNoRisks;
    ChangeGraph();
}

function GraphDropdown() {
    document.getElementById("myDropdown").classList.toggle("show");
}

//#endregion 
