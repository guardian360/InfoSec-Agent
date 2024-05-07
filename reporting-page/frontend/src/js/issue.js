import data from '../database.json' assert { type: 'json' };
import {openIssuesPage} from './issues.js';
import {getLocalization} from './localize.js';
import {retrieveTheme} from './personalize.js';

let stepCounter = 0;
const issuesWithResultsShow = ['11', '60', '70', '80', '90', '100', '110', '160'];

/** Update contents of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 * @param {int} stepCounter Counter specifying the current step
 */
export function updateSolutionStep(solutionText, solutionScreenshot, solution, screenshots, stepCounter) {
  console.log(stepCounter);
  solutionText.innerHTML = `${stepCounter + 1}. ${solution[stepCounter]}`;
  solutionScreenshot.src = screenshots[stepCounter].toString();
}

/** Go to next step of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 */
export function nextSolutionStep(solutionText, solutionScreenshot, solution, screenshots) {
  if (stepCounter < solution.length - 1) {
    stepCounter++;
    updateSolutionStep(solutionText, solutionScreenshot, solution, screenshots, stepCounter);
  }
}

/** Go to previous step of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 */
export function previousSolutionStep(solutionText, solutionScreenshot, solution, screenshots) {
  if (stepCounter > 0) {
    stepCounter--;
    updateSolutionStep(solutionText, solutionScreenshot, solution, screenshots, stepCounter);
  }
}


/** Load the content of the issue page
 *
 * @param {string} issueId Id of the issue to open
 */
export function openIssuePage(issueId) {
  stepCounter = 0;
  const currentIssue = data[issueId];
  // Check if the issue has no screenshots, if so, display that there is no issue (acceptable)
  if (currentIssue.Screenshots.length === 0) {
    const pageContents = document.getElementById('page-contents');
    pageContents.innerHTML = `
      <h1 class="issue-name">${currentIssue.Name}</h1>
      <div class="issue-information">
        <h2 id="information">Information</h2>
        <p id="description">${currentIssue.Information}</p>
        <h2 id="solution">Acceptable</h2>
        <div class="issue-solution">
          <p id="solution-text">${currentIssue.Solution[stepCounter]}</p>
        </div>
        <div class="button" id="back-button">Back to issues overview</div>
      </div>
    `;
  } else { // Issue has screenshots, display the solution guide
    const pageContents = document.getElementById('page-contents');
    if (issuesWithResultsShow.includes(issueId)) {
      pageContents.innerHTML = parseShowResult(issueId, currentIssue);
    } else {
      pageContents.innerHTML = `
        <h1 class="issue-name">${currentIssue.Name}</h1>
        <div class="issue-information">
          <h2 id="information">Information</h2>
          <p>${currentIssue.Information}</p>
          <h2 id="solution">Solution</h2>
          <div class="issue-solution">
            <p id="solution-text">${stepCounter +1}. ${currentIssue.Solution[stepCounter]}</p>
            <img style='display:block; width:750px;height:auto' id="step-screenshot"></img>
            <div class="solution-buttons">
              <div class="button-box">
                <div id="previous-button" class="button">&laquo; Previous step</div>
                <div id="next-button" class="button">Next step &raquo;</div>
              </div>
            </div>
          </div>
          <div class="button" id="back-button">Back to issues overview</div>
        </div>
      `;
    }
    try {
      document.getElementById('step-screenshot').src = currentIssue.Screenshots[stepCounter];
    } catch (error) { }

    // Add functions to page for navigation
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');
    document.getElementById('next-button').addEventListener('click', () =>
      nextSolutionStep(solutionText, solutionScreenshot, currentIssue.Solution, currentIssue.Screenshots));
    document.getElementById('previous-button').addEventListener('click', () =>
      previousSolutionStep(solutionText, solutionScreenshot, currentIssue.Solution, currentIssue.Screenshots));
  }

  document.onload = retrieveTheme();
  const texts = ['information', 'solution', 'previous-button', 'next-button', 'back-button'];
  const localizationIds = ['Issues.Information', 'Issues.Solution', 'Issues.Previous', 'Issues.Next', 'Issues.Back'];
  for (let i = 0; i < texts.length; i++) {
    getLocalization(localizationIds[i], texts[i]);
  }
  document.getElementById('back-button').addEventListener('click', () => openIssuesPage());
}

/** Check if the issue is a show result issue
 *
 * @param {string} issue checks if the issue is a show result issue
 * @return {boolean} if the issue is a show result issue
 */
export function checkShowResult(issue) {
  return issue.Name.includes('Applications with');
}

/** Parse the show result of an issue
 *
 * @param {string} issueId of the issue
 * @param {string} currentIssue of the issue we are looking at
 * @return {string} result of the show result
 */
export function parseShowResult(issueId, currentIssue) {
  let issues = [];
  issues = JSON.parse(sessionStorage.getItem('ScanResult'));
  let resultLine = '';

  switch (issueId) {
  case '11':
    resultLine += `The following devices are or have been connected via bluetooth: <br>`;
    issues.find((issue) => issue.issue_id === 1).result.forEach((issue) => {
      resultLine += `${issue} <br> `;
    });
    break;
  case '60':
    resultLine = permissionShowResults(issues);
    break;
  case '70':
    resultLine = permissionShowResults(issues);
    break;
  case '80':
    resultLine = permissionShowResults(issues);
    break;
  case '90':
    resultLine = permissionShowResults(issues);
    break;
  case '100':
    resultLine = permissionShowResults(issues);
    break;
  case '110':
    resultLine += `The following ports are open: <br>`;
    issues.find((issue) => issue.issue_id === 11).result.forEach((issue) => {
      resultLine += `${issue} <br> `;
    });
    break;
  case '160':
    issues.find((issue) => issue.issue_id === 16).result.forEach((issue) => {
      resultLine += `You changed your password on: ${issue}`;
    });
    break;
  default:
    break;
  }

  /**
   *
   * @param {string} issues with the permission results
   * @return {string} resultLine with the permission results
   */
  function permissionShowResults(issues) {
    let applications = '';
    issues.forEach((issue) => {
      if (issue.issue_id.toString() + issue.result_id.toString() === issueId.toString()) {
        const issueResult = issue.result;
        issueResult.forEach((application) => {
          applications += `${application}, `;
        });
      }
    });
    applications = applications.slice(0, -2);
    resultLine = `The following applications currently have been given permission: ${applications}.`;
    return resultLine;
  }

  const result = `
  <h1 class="issue-name">${currentIssue.Name}</h1>
  <div class="issue-information">
    <h2 id="information">Information</h2>
    <p>${currentIssue.Information}</p>
    <h2 id="information">Findings</h2>
    <p id="description">${resultLine}</p>
    <h2 id="solution">Solution</h2>
    <div class="issue-solution">
      <p id="solution-text">${stepCounter +1}. ${currentIssue.Solution[stepCounter]}</p>
      <img style='display:block; width:750px;height:auto' id="step-screenshot"></img>
      <div class="solution-buttons">
        <div class="button-box">
          <div id="previous-button" class="button">&laquo; Previous step</div>
          <div id="next-button" class="button">Next step &raquo;</div>
        </div>
      </div>
    </div>
    <div class="button" id="back-button">Back to issues overview</div>
  </div>
`;
  return result;
}
