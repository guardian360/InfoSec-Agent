import data from "../database.json"; //get data from database
import "../css/issues.css";
import { openIssuesPage } from "./issues.js";


var solutionStepCounter = 0;
function nextSolutionStep(solution) {
    if (solutionStepCounter < solution.length - 1) {
        solutionStepCounter++;
    }

    var solutionStep = document.getElementById("solution-text");
    solutionStep.innerHTML = `${solution[solutionStepCounter]}`;
}
function previousSolutionStep(solution){
    if (solutionStepCounter > 0){
        solutionStepCounter--;
    }

    var solutionStep = document.getElementById("solution-text");
    solutionStep.innerHTML = `${solution[solutionStepCounter]}`;
}

export function openIssuePage(issueId) {
    console.log("opened issue page: " + issueId);

    const currentIssue = data.find((element) => element.Name = issueId);
    document.getElementById("page-contents").innerHTML = `
        <h1 id="issue-name">${currentIssue.Name}</h1>
        <div id="issue-information">
            <h2>Information</h2>
            <p>${currentIssue.Information}</p>
            <h2>Solution</h2>
            <div id="issue-solution">
                <p id="solution-text">${currentIssue.Solution[solutionStepCounter]}</p>
                <div id="solution-buttons">
                    <div id="previous-button" class="step-button">&laquo; Previous step</div>
                    <div id="next-button" class="step-button">Next step &raquo;</div>
                </div>
            </div>
        </div>
        <div id="back-button">Back to issues overview</div>
    `;

    document.getElementById("next-button").addEventListener("click", () => nextSolutionStep(currentIssue.Solution));
    document.getElementById("previous-button").addEventListener("click", () => previousSolutionStep(currentIssue.Solution));

    document.getElementById("back-button").addEventListener("click", () => openIssuesPage());
}