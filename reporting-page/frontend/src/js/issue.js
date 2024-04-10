import data from '../database.json' assert { type: 'json' };
import {openIssuesPage} from './issues.js';
import {getLocalization} from './localize.js';
import {retrieveTheme} from './personalize.js';

let stepCounter = 0;

/** Update contents of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 * @param {int} stepCounter Counter specifying the current step
 */
export function updateSolutionStep(solutionText, solutionScreenshot, solution, screenshots, stepCounter) {
  solutionText.innerHTML = solution[stepCounter];
  solutionScreenshot.src = screenshots[stepCounter];
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
  const pageContents = document.getElementById('page-contents');
  pageContents.innerHTML = `
    <h1 class="issue-name">${currentIssue.Name}</h1>
    <div class="issue-information">
      <h2 id="information">Information</h2>
      <p>${currentIssue.Information}</p>
      <h2 id="solution">Solution</h2>
      <div class="issue-solution">
        <p id="solution-text">${currentIssue.Solution[stepCounter]}</p>
        <img style='display:block; width:750px;height:auto' id="step-screenshot"></img>
        <div class="solution-buttons">
          <div class="button-box">
            <div id="previous-button" class="step-button">&laquo; Previous step</div>
            <div id="next-button" class="step-button">Next step &raquo;</div>
          </div>
        </div>
      </div>
    </div>
    <div id="back-button">Back to issues overview</div>
  `;

  const texts = ['information', 'solution', 'previous-button', 'next-button', 'back-button'];
  const localizationIds = ['Issues.Information', 'Issues.Solution', 'Issues.Previous', 'Issues.Next', 'Issues.Back'];
  for (let i = 0; i < texts.length; i++) {
    getLocalization(localizationIds[i], texts[i]);
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
  document.getElementById('back-button').addEventListener('click', () => openIssuesPage());

  document.onload = retrieveTheme();
}
