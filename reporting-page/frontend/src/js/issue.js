import dataDe from '../databases/database.de.json' assert { type: 'json' };
import dataEnGB from '../databases/database.en-GB.json' assert { type: 'json' };
import dataEnUS from '../databases/database.en-US.json' assert { type: 'json' };
import dataEs from '../databases/database.es.json' assert { type: 'json' };
import dataFr from '../databases/database.fr.json' assert { type: 'json' };
import dataNl from '../databases/database.nl.json' assert { type: 'json' };
import dataPt from '../databases/database.pt.json' assert { type: 'json' };

import {openIssuesPage, getUserSettings} from './issues.js';
import {getLocalization} from './localize.js';
import {retrieveTheme} from './personalize.js';

let stepCounter = 0;
const issuesWithResultsShow =
    ['11', '21', '60', '70', '80', '90', '100', '110', '160', '173', '201', '230', '271', '311', '320'];

/** Update contents of solution guide
 *
 * @param {HTMLParagraphElement} solutionText Element in which textual solution step is shown
 * @param {HTMLImageElement} solutionScreenshot Element in which screenshot of solution step is shown
 * @param {[string]} solution List of textual solution steps
 * @param {[image]} screenshots List of images of solution steps
 * @param {int} stepCounter Counter specifying the current step
 */
export function updateSolutionStep(solutionText, solutionScreenshot, solution, screenshots, stepCounter) {
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
 * @param {string} severity severity of the issue to open
 */
export async function openIssuePage(issueId, severity) {
  retrieveTheme();
  stepCounter = 0;
  
  //to reload on correct page
  sessionStorage.setItem('savedPage', '8');
  sessionStorage.setItem('issueId', issueId);
  sessionStorage.getItem('severity', severity);

  const language = await getUserSettings();
  let currentIssue;
  switch (language) {
  case 0:
    currentIssue = dataDe[issueId];
    break;
  case 1:
    currentIssue = dataEnGB[issueId];
    break;
  case 2:
    currentIssue = dataEnUS[issueId];
    break;
  case 3:
    currentIssue = dataEs[issueId];
    break;
  case 4:
    currentIssue = dataFr[issueId];
    break;
  case 5:
    currentIssue = dataNl[issueId];
    break;
  case 6:
    currentIssue = dataPt[issueId];
    break;
  default:
    currentIssue = dataEnGB[issueId];
  }

  // Check if the issue has no screenshots, if so, display that there is no issue (acceptable)
  if (severity == 0) {
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
          <h2 id="information" class="lang-information"></h2>
          <p>${currentIssue.Information}</p>
          <h2 id="solution" class="lang-solution"></h2>
          <div class="issue-solution">
            <p id="solution-text">${stepCounter +1}. ${currentIssue.Solution[stepCounter]}</p>
            <img style='display:block; width:750px;height:auto' id="step-screenshot"></img>
            <div class="solution-buttons">
              <div class="button-box">
                <div id="previous-button" class="lang-previous-button button"></div>
                <div id="next-button" class="lang-next-button button">;</div>
              </div>
            </div>
          </div>
          <div class="lang-back-button button" id="back-button"></div>
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

  const texts = ['lang-information', 'lang-solution', 'lang-previous-button', 'lang-next-button', 'lang-back-button'];
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
    generateBulletList(issues, 1);
    break;
  case '21':
    generateBulletList(issues, 2);
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
    resultLine += `The following processes are currently running on your device on the following ports: <br>`;
    const portTable = processPortsTable(issues.find((issue) => issue.issue_id === 11).result);
    resultLine += `<table class = "issues-table">`;
    resultLine += `<thead><tr><th>Process</th><th>Port(s)</th></tr></thead>`;
    portTable.forEach((entry) => {
      resultLine += `<tr><td style="width: 30%">${entry.portProcess}</td>
        <td style="width: 30%">${entry.ports.join('<br>')}</td></tr>`;
    });
    resultLine += '</table>';
    break;
  case '160':
    issues.find((issue) => issue.issue_id === 16).result.forEach((issue) => {
      resultLine += `You changed your password on: ${issue}`;
    });
    break;
  case '173':
    generateBulletList(issues, 17);
    break;
  case '201':
    generateBulletList(issues, 20);
    break;
  case '230':
    generateBulletList(issues, 23);
    break;
  case '271':
    resultLine += '(Possible) tracking cookies have been found from the following websites:';
    resultLine += cookiesTable(issues.find((issue) => issue.issue_id === 27).result);
    break;
  case '311':
    generateBulletList(issues, 31);
    break;
  case '320':
    const cisTable = cisregistryTable(issues.find((issue) => issue.issue_id === 32).result);
    resultLine += `<table class = "issues-table">`;
    cisTable.forEach((entry) => {
      resultLine += `<tr><td style="width: 30%; word-break: break-all">${entry.registryKey}</td>
        <td>${entry.values.join('<br>')}</td></tr>`;
    });
    resultLine += '</table>';
    break;
  default:
    break;
  }

  /**
   * Generate a bullet list for each entry of a result of certain issues
   * @param {string} issues to generate a bullet list for
   * @param {int} issueId of the issue
   * @return {string} html tags for a bullet list
   */
  function generateBulletList(issues, issueId) {
    resultLine += `<ul>`;
    issues.find((issue) => issue.issue_id === issueId).result.forEach((issue) => {
      resultLine += `<li>${issue}</li>`;
    });
    resultLine += `</ul>`;
    return resultLine;
  }

  /**
   *
   * @param {string} issues with the permission results
   * @return {string} resultLine with the permission results
   */
  function permissionShowResults(issues) {
    let applications = '<ul>';
    issues.forEach((issue) => {
      if (issue.issue_id.toString() + issue.result_id.toString() === issueId.toString()) {
        const issueResult = issue.result;
        issueResult.forEach((application) => {
          applications += `<li>${application}</li>`;
        });
      }
    });
    applications += '</ul>'; // Close the list
    resultLine = `The following applications currently have been given permission:<br>${applications}`;
    return resultLine;
  }

  /**
   * Create a table for the CIS registry issues
   * @param {string} issues list of incorrect registry keys
   * @return {*[]} table with registry keys and values
   */
  function cisregistryTable(issues) {
    const table = [];
    let currentKey = null;
    let currentValues = [];

    issues.forEach((issue) => {
      if (issue.includes('SYSTEM') || issue.includes('SOFTWARE')) {
        if (currentKey) {
          table.push({registryKey: currentKey, values: currentValues});
        }
        currentKey = issue;
        currentValues = [];
      } else if (currentKey) {
        currentValues.push(issue);
      }
    });

    if (currentKey) {
      table.push({registryKey: currentKey, values: currentValues});
    }

    return table;
  }

  /**
   * Create a table for the process ports
   * @param {string} issues list of processes and ports
   * @return {*[]} table with process names and ports
   */
  function processPortsTable(issues) {
    const table = [];
    issues.forEach((issue) => {
      const parts = issue.split(/[ ,]+/); // Split on space and comma
      const processIndex = parts.indexOf('process:');
      const portIndex = parts.indexOf('port:');

      if (processIndex !== -1 && portIndex !== -1) {
        const processName = parts[processIndex + 1];
        const ports = new Set(parts.slice(portIndex + 1));
        table.push({portProcess: processName, ports: Array.from(ports)});
      }
    });

    return table;
  }

  /**
   * Create a table to display found (possible) tracking cookies
   * @param {string} issues list of cookies and their host
   * @return {string} HTML table with cookies and their host
   */
  function cookiesTable(issues) {
    const cookiesByHost = {};
    for (let i = 0; i < issues.length; i += 2) {
      const host = issues[i+1];

      if (!cookiesByHost[host]) {
        cookiesByHost[host] = true;
      }
    }

    // Generate HTML for table
    let tableHTML = '<table class="issues-table">';
    for (const host in cookiesByHost) {
      if (cookiesByHost.hasOwnProperty(host)) {
        tableHTML += `<tr><td style="width: 30%; word-break: break-all">${host}</td></tr>`;
      }
    }
    tableHTML += '</table>';

    return tableHTML;
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
