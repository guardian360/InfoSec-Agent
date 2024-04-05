import 'jsdom-global/register.js';
import test from 'unit.js';
import { JSDOM } from "jsdom";
import { openIssuePage } from "../src/js/issue.js";
import { updateSolutionStep } from "../src/js/issue.js";
import { nextSolutionStep } from "../src/js/issue.js";
import { previousSolutionStep } from "../src/js/issue.js";
import data from "../src/database.json" assert { type: "json" };

// Mock issue page
const issuesDOM = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div id="page-contents"></div>
</body>
</html>
`);
global.document = issuesDOM.window.document
global.window = issuesDOM.window

var stepCounter = 0;
const solution = ["Step 1", "Step 2", "Step 3"];
const screenshots = ["screenshot1.jpg", "screenshot2.jpg", "screenshot3.jpg"];
const solutionStep = document.createElement("p");
const screenshot = document.createElement("img");
solutionStep.id = "solution-text";
screenshot.id = "step-screenshot";
document.body.appendChild(solutionStep);
document.body.appendChild(screenshot);

// unit test for updateSolutionStep
describe('updateSolutionStep', function() {
  it('should update the solution step', function() {
    // Arrange
    stepCounter = 0;

    // Act
    updateSolutionStep(solution, screenshots, stepCounter);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo("Step 1");
    test.value(screenshot.src).isEqualTo("screenshot1.jpg");
  });
});

// unit test for nextSolutionStep
describe('nextSolutionStep', function() {
  it('should go to the next step', function() {
    // Arrange
    stepCounter = 0;

    // Act
    nextSolutionStep(solution, screenshots);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo("Step 2");
    test.value(screenshot.src).isEqualTo("screenshot2.jpg");
  });
});

// unit test for previousSolutionStep
describe('previousSolutionStep', function() {
  it('should go to the previous step', function() {
    // Arrange
    stepCounter = 1;

    // Act
    previousSolutionStep(solution, screenshots);

    // Assert
    test.value(solutionStep.innerHTML).isEqualTo("Step 1");
    test.value(screenshot.src).isEqualTo("screenshot1.jpg");
  });
});