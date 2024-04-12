import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {updateSolutionStep} from '../src/js/issue.js';
import {nextSolutionStep} from '../src/js/issue.js';
import {previousSolutionStep} from '../src/js/issue.js';

global.TESTING = true;

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
  const solution = ['Step 1', 'Step 2', 'Step 3'];
  const screenshots = ['screenshot1.jpg', 'screenshot2.jpg', 'screenshot3.jpg'];
  const solutionStep = document.createElement('p');
  const screenshot = document.createElement('img');
  solutionStep.id = 'solution-text';
  screenshot.id = 'step-screenshot';
  document.body.appendChild(solutionStep);
  solutionStep.innerHTML = 'Step 1';
  document.body.appendChild(screenshot);
  screenshot.src = 'screenshot1.jpg';

  it('updateSolutionStep should update the solution step', function() {
    // Arrange
    stepCounter = 0;

    // Act
    updateSolutionStep(solutionStep, screenshot, solution, screenshots, stepCounter);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo('Step 1');
    test.value(screenshot.src).isEqualTo('screenshot1.jpg');
  });

  it('nextSolutionStep should update the current step and screenshot', function() {
    // Arrange
    stepCounter = 0;

    // Act
    nextSolutionStep(solutionStep, screenshot, solution, screenshots);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo('Step 2');
    test.value(screenshot.src).isEqualTo('screenshot2.jpg');
  });

  it('previousSolutionStep should update the current step and screenshot', function() {
    // Arrange
    stepCounter = 1;

    // Act
    previousSolutionStep(solutionStep, screenshot, solution, screenshots);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo('Step 1');
    test.value(screenshot.src).isEqualTo('screenshot1.jpg');
  });
});
