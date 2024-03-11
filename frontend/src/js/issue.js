import data from "../database.json" assert { type: "json" };
import { openIssuesPage } from "./issues.js";

let stepCounter = 0;

// Update contents of solution guide
function updateSolutionStep(solution, screenshots, stepCounter) {
    const solutionStep = document.getElementById("solution-text");
    solutionStep.innerHTML = solution[stepCounter];
    document.getElementById("step-screenshot").src = screenshots[stepCounter];
}

// Go to next step of solution guide
function nextSolutionStep(solution, screenshots) {
    if (stepCounter < solution.length - 1) {
        stepCounter++;
        updateSolutionStep(solution, screenshots, stepCounter);
    }
}

// Go to previous step of solution guide
function previousSolutionStep(solution, screenshots) {
    if (stepCounter > 0) {
        stepCounter--;
        updateSolutionStep(solution, screenshots, stepCounter);
    }
}

// Load the content of the issue page
export function openIssuePage(issueId) {
    stepCounter = 0;
    const currentIssue = data.find((element) => element.Name === issueId);
    const pageContents = document.getElementById("page-contents");
    pageContents.innerHTML = `
        <h1 id="issue-name">${currentIssue.Name}</h1>
        <div id="issue-information">
            <h2>Information</h2>
            <p>${currentIssue.Information}</p>
            <h2>Solution</h2>
            <div id="issue-solution">
                <p id="solution-text">${currentIssue.Solution[stepCounter]}</p>
                <img style='display:block; width:500px;height:auto' id="step-screenshot"></img>
                <div id="solution-buttons">
                    <div id="button-box">
                        <div id="previous-button" class="step-button">&laquo; Previous step</div>
                        <div id="next-button" class="step-button">Next step &raquo;</div>
                    </div>
                </div>
            </div>
        </div>
        <div id="back-button">Back to issues overview</div>
    `;

    try {
        document.getElementById("step-screenshot").src = currentIssue.Screenshots[stepCounter];
    } catch (error) { }

    document.getElementById("next-button").addEventListener("click", () => nextSolutionStep(currentIssue.Solution, currentIssue.Screenshots));
    document.getElementById("previous-button").addEventListener("click", () => previousSolutionStep(currentIssue.Solution, currentIssue.Screenshots));
    document.getElementById("back-button").addEventListener("click", () => openIssuesPage());
}
