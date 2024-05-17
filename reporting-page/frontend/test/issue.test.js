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

// Mock often used page functions
mockPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issues.js', () => ({
  openIssuesPage: jest.fn(),
}));

// Mock sessionStorage
global.sessionStorage = storageMock;

describe('Issue page', function() {
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
  it('clicking on the back button should call openIssuesPage', async function() {
    // Arrange
    const issues = await import('../src/js/issues.js');
    const button = document.getElementById('back-button');
    const openIssuesPageMock = jest.spyOn(issues, 'openIssuesPage');

    // Act
    button.dispatchEvent(clickEvent);

    // Assert
    expect(openIssuesPageMock).toHaveBeenCalled();
  });

  // from here on issueID 160 is used for tests up to parseShowResults tests
  const issueID = 160;
  let currentIssue = data[issueID];

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

    // Act
    // calls previousSolutionStep
    document.getElementById('previous-button').dispatchEvent(clickEvent);

    // Assert
    test.value(htmlDecode(solutionText.innerHTML)).isEqualTo('1. ' + currentIssue.Solution[0]);
    test.value(solutionScreenshot.src).isEqualTo(currentIssue.Screenshots[0]);
  });

  it('clicking previous step button on first step should not update the current step and screenshot', async function() {
    // Arrange
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');

    // Act
    // calls previousSolutionStep
    document.getElementById('previous-button').dispatchEvent(clickEvent);

    // Assert
    test.value(htmlDecode(solutionText.innerHTML)).isEqualTo('1. ' + currentIssue.Solution[0]);
    test.value(solutionScreenshot.src).isEqualTo(currentIssue.Screenshots[0]);  
  });
  it('clicking next step button at last step should not update the current step and screenshot', async function() {
    // Arrange
    const solutionText = document.getElementById('solution-text');
    const solutionScreenshot = document.getElementById('step-screenshot');

    // Act
    // starts on step 1. 
    // calls nextSolutionStep
    document.getElementById('next-button').dispatchEvent(clickEvent);
    document.getElementById('next-button').dispatchEvent(clickEvent);
    // At step 3 now, click next one more time.
    document.getElementById('next-button').dispatchEvent(clickEvent);

    // Assert
    test.value(solutionText.innerHTML).isEqualTo('3. ' + currentIssue.Solution[2]);
    test.value(solutionScreenshot.src).isEqualTo(currentIssue.Screenshots[2]);
  });
  
  // Mock scan results for the parseShowResult function 
  const issueResult_ids = [
    [1,1],
    [2,1],
    [6,0],
    [7,0],
    [8,0],
    [9,0],
    [10,0],
    [11,0],
    [16,0],
    [17,3],
    [20,1],
    [23,0],
    [27,1],
    [31,1],
    [32,0]];
  // Mock scan results
  const mockResult = []
  issueResult_ids.forEach((ir) => {
    mockResult.push({
      issue_id: ir[0],
      result_id: ir[1],
      result: [
        'process: p, port: 1, 2, 3',
        'SYSTEM',
        'CIS registry 1',
        'SOFTWARE',
        'CIS registry 2',      
      ]
    })
  });

  it('Should fill page with parseShowResult with a set of resultIDs', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    sessionStorage.setItem('ScanResult', JSON.stringify(mockResult));

    mockResult.forEach((result, index) => {
      let jsonkey = result.issue_id.toString() + result.result_id.toString();
      currentIssue = data[jsonkey];

      // Act
      issue.openIssuePage(jsonkey);
      const name = document.getElementsByClassName('issue-name')[0].innerHTML;
      const description = document.getElementById('information').nextElementSibling.innerHTML;
      const solution = document.getElementById('solution-text').innerHTML;
  
      // Assert
      test.value(name).isEqualTo(currentIssue.Name);
      test.value(description).isEqualTo(currentIssue.Information);
      test.value(htmlDecode(solution)).isEqualTo('1. ' + currentIssue.Solution[0]);
    })
  });
  it('parseShowResult fills the page with the correct structure for specific results', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    // expectedFindings should be changed if the structure for specific results is changed in the code
    const expectedFindings = [
      '<li>process: p, port: 1, 2, 3</li><li>SYSTEM</li><li>CIS registry 1</li><li>SOFTWARE</li><li>CIS registry 2</li>',
      '<thead><tr><th>Process</th><th>Port(s)</th></tr></thead><tbody><tr><td style="width: 30%">p</td>\n' +
      '        <td style="width: 30%">1<br>2<br>3</td></tr></tbody>',
      'You changed your password on: process: p, port: 1, 2, 3You changed your password on: SYSTEM' +
      'You changed your password on: CIS registry 1You changed your password on: SOFTWAREYou changed your password on: CIS registry 2',
      '<tbody><tr><td style="width: 30%; word-break: break-all">SYSTEM</td></tr><tr><td style="width: 30%; word-break: break-all">' +
      'SOFTWARE</td></tr><tr><td style="width: 30%; word-break: break-all">undefined</td></tr></tbody>',
      '<tbody><tr><td style="width: 30%; word-break: break-all">SYSTEM</td>\n' +
      '        <td>CIS registry 1</td></tr><tr><td style="width: 30%; word-break: break-all">SOFTWARE</td>\n' +
      '        <td>CIS registry 2</td></tr></tbody>',
    ]

    sessionStorage.setItem('ScanResult', JSON.stringify(mockResult));

    mockResult.forEach((result, index) => {
      let jsonkey = result.issue_id.toString() + result.result_id.toString();
      currentIssue = data[jsonkey];

      // Act
      issue.openIssuePage(jsonkey);

      if (index < 7 || (index > 8 && index < 12) || index == 13) {
        // called to generateBulletList and permissionShowResults
        const findings = document.getElementById('description').nextElementSibling.innerHTML;
  
        // Assert
        test.value(findings).isEqualTo(expectedFindings[0]);
      } else if (index == 7) {
        // called to processPortsTable
        const findings = document.getElementById('description').nextElementSibling.innerHTML;
  
        // Assert
        test.value(findings).isEqualTo(expectedFindings[1]);
      } else if (index == 8) {
        const findings = document.getElementById('description').innerHTML;
  
        // Assert
        test.value(findings).isEqualTo(expectedFindings[2]);
      } else if (index == 12) {
        // called to cookiesTable
        const findings = document.getElementById('description').nextElementSibling.innerHTML;

        // Assert
        // console.log(findings);
        test.value(findings).isEqualTo(expectedFindings[3]);
      } else {
        // called to cisregristryTable
        const findings = document.getElementsByClassName('issues-table')[0].innerHTML;

        // Assert
        test.value(findings).isEqualTo(expectedFindings[4]);
      }
    })
  })
  it('parseShowResult keeps findings empty if the issueID is not in the issuesWithResultsShow list', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    // Mock scan result
    const mockResult = [
      {
        issue_id: 1,
        result_id: 0,
        result: [
          "findings",        
        ]
      },
    ]
    sessionStorage.setItem('ScanResult', JSON.stringify(mockResult));
    let jsonkey = mockResult[0].issue_id.toString() + mockResult[0].result_id.toString();
    currentIssue = data[jsonkey];
    const pageContents = document.getElementById('page-contents');

    // Act
    pageContents.innerHTML = issue.parseShowResult(jsonkey, currentIssue);
    const findings = document.getElementById('description').innerHTML;

    // Assert
    test.value(findings).isEqualTo(''); 
  })
  it('checkShowResult should check if an issue name contains "applications with"', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    // Act
    const checked = issue.checkShowResult(data[60])

    // Assert
    test.value(checked).isEqualTo(true);
  });
});
