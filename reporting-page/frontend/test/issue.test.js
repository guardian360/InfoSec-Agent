import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import data from '../src/database.json' assert { type: 'json' };
import {mockPageFunctions,mockGetLocalization,clickEvent,storageMock} from './mock.js';

global.TESTING = true;

function htmlDecode(input){
  var e = document.createElement('div');
  e.innerHTML = input;
  return e.childNodes[0].nodeValue;
}

// Mock sessionStorage
global.sessionStorage = storageMock;

describe('Issue page', function() {
  // Mock issue page
  const issuesDOM = new JSDOM(`
  <!DOCTYPE html>
  <html>
  <body>
      <div id="page-contents"></div>
  </body>
  </html>
  `);
  global.document = issuesDOM.window.document;
  global.window = issuesDOM.window;

  let stepCounter = 0;
  // const solution = ['Step 1', 'Step 2', 'Step 3'];
  // const screenshots = ['screenshot1.jpg', 'screenshot2.jpg', 'screenshot3.jpg'];
  // const solutionStep = document.createElement('p');
  // const screenshot = document.createElement('img');
  // solutionStep.id = 'solution-text';
  // screenshot.id = 'step-screenshot';
  // document.body.appendChild(solutionStep);
  // solutionStep.innerHTML = 'Step 1';
  // document.body.appendChild(screenshot);
  // screenshot.src = 'screenshot1.jpg';

  // Mock often used page functions
  mockPageFunctions();

  // Mock Localize function
  jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
    Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
  }));

  it('openIssuesPage should not add solutions for a non-issue to the page-contents', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    const nonIssueID = 161
    const currentIssue = data[nonIssueID];

    // Act
    await issue.openIssuePage(nonIssueID);
    const name = document.getElementsByClassName('issue-name')[0].innerHTML;
    const description = document.getElementById('description').innerHTML;
    const solution = document.getElementById('solution-text').innerHTML;

    // Assert
    test.value(name).isEqualTo(currentIssue.Name);
    test.value(description).isEqualTo(currentIssue.Information);
    test.value(solution).isEqualTo(currentIssue.Solution[0]);

  });

  // from here on issueID 160 is used for testing
  const issueID = 160;
  const currentIssue = data[issueID];

  it('openIssuesPage should add the right info about the issue to the page-contents', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    // Act
    await issue.openIssuePage(issueID);
    const name = document.getElementsByClassName('issue-name')[0].innerHTML;
    const description = document.getElementById('information').nextElementSibling.innerHTML;
    const solution = document.getElementById('solution-text').innerHTML;

    // Assert
    test.value(name).isEqualTo(currentIssue.Name);
    test.value(description).isEqualTo(currentIssue.Information);
    test.value(htmlDecode(solution)).isEqualTo('1. ' + currentIssue.Solution[0]);
  });
  it('updateSolutionStep should update the solution step', async function() {
    // Arrange
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');
    stepCounter = 0;

    const issue = await import('../src/js/issue.js');

    // Act
    issue.updateSolutionStep(solutionText, solutionScreenshot, currentIssue.Solution, currentIssue.Screenshots, stepCounter);

    // Assert
    test.value(htmlDecode(solutionText.innerHTML)).isEqualTo('1. ' + currentIssue.Solution[0]);
    test.value(solutionScreenshot.src).isEqualTo(currentIssue.Screenshots[0]);
  });

  it('nextSolutionStep should update the current step and screenshot', async function() {
    // Arrange
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');

    const issue = await import('../src/js/issue.js');

    // Act
    // calls nextSolutionStep
    document.getElementById('next-button').dispatchEvent(clickEvent);

    // Assert
    test.value(solutionText.innerHTML).isEqualTo('2. ' + currentIssue.Solution[1]);
    test.value(solutionScreenshot.src).isEqualTo(currentIssue.Screenshots[1]);
  });

  it('previousSolutionStep should update the current step and screenshot', async function() {
    // Arrange
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');

    const issue = await import('../src/js/issue.js');

    // Act
    // calls previousSolutionStep
    document.getElementById('previous-button').dispatchEvent(clickEvent);

    // Assert
    test.value(htmlDecode(solutionText.innerHTML)).isEqualTo('1. ' + currentIssue.Solution[0]);
    test.value(solutionScreenshot.src).isEqualTo(currentIssue.Screenshots[0]);
  });

  it('Should fill page with parseShowResult with a set of resultIDs', async function() {

  })
});
