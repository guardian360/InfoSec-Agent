import "../css/home.css";
import "../css/dashboard.css";
import "../css/color-palette.css";

// import "./graph"
// import "./piechart"
// import "./risk-counters"

function openDashboardPage() {
    document.getElementById("page-contents").innerHTML = `
    <div class="dashboard-data">
        <div class="data-column risk-counters">
            <div class="security-status">
                <p>Security status</p>
                <p>Critical</p>
            </div>
            <div class="risk-counter high-risk">
                <p>High risk issues</p>
                <p id="high-risk-counter">4</p>
            </div>
            <div class="risk-counter medium-risk">
                <p>Medium risk issues</p>
                <p id="medium-risk-counter">4</p>
            </div>
            <div class="risk-counter low-risk">
                <p>Low risk issues</p>
                <p id="low-risk-counter">4</p>
            </div>
            <div class="risk-counter no-risk">
                <p>Safe issues</p>
                <p id="no-risk-counter">4</p>
            </div>
        </div>
        <div class="data-column piechart">
            <canvas id="pieChart"></canvas>
        </div>
        <div class="data-column issue-buttons">
            <H2>You have some issues you can fix. 
                To start resolving a issue either navigate to the issues page, or pick a suggested issue below
            </H2>
            <a class="issue-button">Suggested Issue</a>
            <a class="issue-button">Quick Fix</a>
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
            <button class="dropbtn" onclick="GraphDropdown()">Select Risks</button>
            <div class="dropdown-selector" id="myDropdown">
                <p><input type="checkbox" checked="true" value="true" id="select-high-risk" onchange="ToggleHighRisks()">
                    <label for="select-high-risk"> High risks</label><br>
                </p>
                <p><input type="checkbox" checked="true" value="true" id="select-medium-risk" onchange="ToggleMediumRisks()">
                    <label for="select-medium-risk"> Medium risks</label>
                </p>
                <p><input type="checkbox" checked="true" value="true" id="select-low-risk" onchange="ToggleLowRisks()">
                    <label for="select-low-risk"> Low risks</label>
                </p>
                <p><input type="checkbox" checked="true" value="true" id="select-no-risk" onchange="ToggleNoRisks()">
                    <label for="select-no-risk"> Safe</label>
                </p>
            </div>
        </div>
        <a class="interval-button"><p>Change interval</p><input type="number" value="5" id="graph-interval" onchange="ChangeGraph()"></a>
    </div>
    <div class="graph-column issues-graph">
        <canvas id="graph"></canvas>
    </div>
  </div>
    `;
}

document.getElementById("dashboard-button").addEventListener("click", () => openDashboardPage());