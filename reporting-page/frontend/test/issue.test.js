import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';

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

  // Mock LogError
  jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
    LogError: jest.fn(),
  }));

  it('updateSolutionStep should update the solution step', async function() {
    // Arrange
    stepCounter = 0;

    const issue = await import('../src/js/issue.js');

    // Act
    issue.updateSolutionStep(solutionStep, screenshot, solution, screenshots, stepCounter);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo('1. Step 1');
    test.value(screenshot.src).isEqualTo('screenshot1.jpg');
  });

  it('nextSolutionStep should update the current step and screenshot', async function() {
    // Arrange
    stepCounter = 0;

    const issue = await import('../src/js/issue.js');

    // Act
    issue.nextSolutionStep(solutionStep, screenshot, solution, screenshots);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo('2. Step 2');
    test.value(screenshot.src).isEqualTo('screenshot2.jpg');
  });

  it('previousSolutionStep should update the current step and screenshot', async function() {
    // Arrange
    stepCounter = 1;

    const issue = await import('../src/js/issue.js');

    // Act
    issue.previousSolutionStep(solutionStep, screenshot, solution, screenshots);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo('1. Step 1');
    test.value(screenshot.src).isEqualTo('screenshot1.jpg');
  });
});
