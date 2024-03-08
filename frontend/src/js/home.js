import "../css/home.css";
import "../css/color-palette.css";

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
        </div>
    </div>
    <h2 id="title-medals">Medals</h2>
    <div class="container">  
        <div class="medal-layout">
            <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
            <p id="medal-name"> Medal 1</p>
            <p> 01-03-2024</p>
        </div>
        <div class="medal-layout">
            <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
            <p id="medal-name"> Medal 2</p>
            <p> 01-03-2024</p>
        </div>
        <div class="medal-layout">
            <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
            <p id="medal-name"> Medal 3</p>
            <p> 01-03-2024</p>
        </div><div class="medal-layout">
            <img src="src/assets/images/img_medal1.jpg" alt="Photo of medal">
            <p id="medal-name"> Medal 1</p>
            <p> 01-03-2024</p>
        </div>
    </div>
    `;
}

document.onload = openHomePage();
document.getElementById("logo-button").addEventListener("click", () => openHomePage());
document.getElementById("home-button").addEventListener("click", () => openHomePage());