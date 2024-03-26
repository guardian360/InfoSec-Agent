import {PieChart} from "./piechart";
import {ScanNow} from '../../wailsjs/go/main/Tray';
import {RiskCounters} from "./risk-counters";

/** Load the content of the Home page */
function openHomePage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="container-data">       
    <div class="data-column risk-counters">     
      <div class="data-column piechart">
        <canvas id="pieChart"></canvas>
      </div>
    </div>
    <div class="data-column issue-buttons">
      <H2>You have some issues you can fix. 
        To start resolving a issue either navigate to the issues page, or pick a suggested issue below
      </H2>
      <a class="issue-button">Suggested Issue</a>
      <a class="issue-button">Quick Fix</a>
      <a class="issue-button" id="scan-button">Scan Now</a>
    </div>
  </div>
  <h2 class="title-medals">Medals</h2>
  <div class="container">  
    <div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 1</p>
      <p> 01-03-2024</p>
    </div>
    <div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 2</p>
      <p> 01-03-2024</p>
    </div>
    <div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 3</p>
      <p> 01-03-2024</p>
    </div><div class="medal-layout">
      <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
      <p class="medal-name"> Medal 1</p>
      <p> 01-03-2024</p>
    </div>
  </div>
  `;  
    let rc = new RiskCounters();
    new PieChart("pieChart",rc);

    document.getElementById("scan-button").addEventListener("click", () => scanNow());
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

document.onload = openHomePage();