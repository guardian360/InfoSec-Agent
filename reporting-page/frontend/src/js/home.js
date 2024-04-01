import {PieChart} from "./piechart";
import {ScanNow} from '../../wailsjs/go/main/Tray';
import { GetLocalization } from './localize.js';
import { CloseNavigation } from "./navigation-menu.js";
import { MarkSelectedNavigationItem } from "./navigation-menu.js";
import medal from '../assets/images/img_medal1.jpg';
import { retrieveTheme } from "./personalize";

/** Load the content of the Home page */
export function openHomePage() {
  CloseNavigation();
  MarkSelectedNavigationItem("home-button");
  
  document.getElementById("page-contents").innerHTML = `
  <div class="home-data">
    <div class="container-data home-column-one"> 
      <div class="data-column risk-counters">     
        <div class="data-column data-segment piechart">
          <div class="data-segment-header">
            <p class="piechart-header">Risk level distribution</p>
          </div>
          <div class="piechart-container">
            <canvas id="pieChart"></canvas>
          </div>
        </div>
      </div>
      <div class="data-column data-segment issue-buttons">
        <div class="data-segment-header">
          <p class="choose-issue-description"></p>
        </div>
        <a class="issue-button suggested-issue">Suggested Issue</a>
        <a class="issue-button quick-fix">Quick Fix</a>
        <a class="issue-button scan-now">Scan Now</a>
      </div>
    </div>
    <div class="data-segment progress">  
      <div class="data-segment-header">
        <p class="title-medals"></p>
      </div>
      <div class="medals">
        <div class="medal-layout">
          <img id="medal" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 1</p>
          <p> 01-04-2024</p>
        </div>
        <div class="medal-layout">
          <img id="medal2" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 2</p>
          <p> 01-04-2024</p>
        </div>
        <div class="medal-layout">
          <img id="medal3" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 3</p>
          <p> 01-04-2024</p>
        </div><div class="medal-layout">
          <img id="medal4" alt="Photo of medal"></img>
          <p class="medal-name"> Medal 4</p>
          <p> 01-04-2024</p>
        </div>
      </div>
    </div>
  </div>
  `;  

    document.getElementById('medal').src = medal;
    document.getElementById('medal2').src = medal;
    document.getElementById('medal3').src = medal;
    document.getElementById('medal4').src = medal;

    let rc = JSON.parse(sessionStorage.getItem("RiskCounters"));
    new PieChart("pieChart",rc);

  // Localize the static content of the home page
  let staticHomePageConent = [
    "suggested-issue", 
    "quick-fix", 
    "scan-now", 
    "title-medals", 
    "security-status",
    "high-risk-issues", 
    "medium-risk-issues",
    "low-risk-issues",
    "safe-issues",
    "choose-issue-description"
    ]
    let localizationIds = [
      "Dashboard.SuggestedIssue", 
      "Dashboard.QuickFix", 
      "Dashboard.ScanNow", 
      "Dashboard.Medals", 
      "Dashboard.SecurityStatus",
      "Dashboard.HighRisk", 
      "Dashboard.MediumRisk",
      "Dashboard.LowRisk",
      "Dashboard.Safe",
      "Dashboard.ChooseIssueDescription"
    ]
    for (let i = 0; i < staticHomePageConent.length; i++) {
        GetLocalization(localizationIds[i], staticHomePageConent[i])
  }

    document.getElementsByClassName("scan-now")[0].addEventListener("click", () => scanNow());
    document.getElementById("home-button").addEventListener("click", () => openHomePage());
    document.getElementById("logo").innerHTML = localStorage.getItem("picture");

    document.onload = retrieveTheme();
}

document.getElementById("logo-button").addEventListener("click", () => openHomePage());
document.getElementById("home-button").addEventListener("click", () => openHomePage());

function scanNow() {
    ScanNow()
    .then((result) => {
    })
    .catch((err) => {
        console.error(err);
    });
}

//document.onload = openHomePage();

window.onload = function() {
    let savedImage = localStorage.getItem('picture');
    let savedText = localStorage.getItem('title');
    let savedIcon = localStorage.getItem('favicon');
    if (savedImage) {
      let logo = document.getElementById('logo');
      logo.src = savedImage;
    }
    if (savedText) {
      let title = document.getElementById('title');
      title.textContent = savedText;
    }
    if(savedIcon){
      let favicon = document.getElementById('favicon');
      favicon.href = savedIcon;
    }
  };