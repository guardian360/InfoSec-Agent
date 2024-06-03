import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import data from '../src/databases/database.en-GB.json' assert { type: 'json' };
import {mockPageFunctions, mockGetLocalization, clickEvent, storageMock} from './mock.js';

global.TESTING = true;

/** removes html elements form a string
 *
 * @param {string} input string with html elements
 * @return {string} input without html element
 */
function htmlDecode(input) {
  const e = document.createElement('div');
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

// Mock scantest function
jest.unstable_mockModule('../src/js/database.js', () => ({
  scanTest: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issues.js', () => ({
  openIssuesPage: jest.fn(),
  getUserSettings: jest.fn().mockImplementationOnce(() => 1)
    .mockImplementationOnce(() => 1)
    .mockImplementationOnce(() => 0)
    .mockImplementationOnce(() => 1)
    .mockImplementationOnce(() => 2)
    .mockImplementationOnce(() => 3)
    .mockImplementationOnce(() => 4)
    .mockImplementationOnce(() => 5)
    .mockImplementationOnce(() => 6)
    .mockImplementationOnce(() => 7)
    .mockImplementation(() => 1),
}));

// Mock sessionStorage
global.sessionStorage = storageMock;

describe('Issue page', function() {
  it('openIssuesPage should not add solutions for a non-issue to the page-contents', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    const nonIssueID = 161;
    const currentIssue = data[nonIssueID];

    // Act
    await issue.openIssuePage(nonIssueID, 0);
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
  const severity = 2;
  let currentIssue = data[issueID];

  it('openIssuesPage should add the right info about the issue to the page-contents', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    // Act
    await issue.openIssuePage(issueID, severity);
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
    issue.updateSolutionStep(solutionText, solutionScreenshot,
      currentIssue, stepCounter);

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


  it('openIssuePage should open the right localized issue database with localizations', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    const localizations =[
      {path: '../src/databases/database.de.json', lang: 0},
      {path: '../src/databases/database.en-GB.json', lang: 1},
      {path: '../src/databases/database.en-US.json', lang: 2},
      {path: '../src/databases/database.es.json', lang: 3},
      {path: '../src/databases/database.fr.json', lang: 4},
      {path: '../src/databases/database.nl.json', lang: 5},
      {path: '../src/databases/database.pt.json', lang: 6},
      {path: '../src/databases/database.en-GB.json', lang: 999},
    ];
    for (const localization of localizations) {
      // Act
      const data = await import(localization.path, {assert: {type: 'json'}});
      const currentIssue = data.default[issueID];
      await issue.openIssuePage(issueID, severity);

      const name = document.getElementsByClassName('issue-name')[0].innerHTML;
      const description = document.getElementById('information').nextElementSibling.innerHTML;
      const solution = document.getElementById('solution-text').innerHTML;

      // Assert
      test.value(name).isEqualTo(currentIssue.Name);
      test.value(description).isEqualTo(currentIssue.Information);
      test.value(htmlDecode(solution)).isEqualTo('1. ' + currentIssue.Solution[0]);
    }
  });

  // Mock scan results for the parseShowResult function
  const issueResultIds = [
    [1, 1],
    [2, 1],
    [6, 0],
    [7, 0],
    [8, 0],
    [9, 0],
    [10, 0],
    [11, 0],
    [16, 0],
    [17, 3],
    [20, 1],
    [23, 0],
    [27, 1],
    [31, 1],
    [32, 0],
    [35, 1],
    [36, 1]];
  // Mock scan results
  const mockResult = [];
  issueResultIds.forEach((ir) => {
    mockResult.push({
      issue_id: ir[0],
      result_id: ir[1],
      result: [
        'process: p, port: 1, 2, 3',
        'SYSTEM',
        'CIS registry 1',
        'SOFTWARE',
        'CIS registry 2',
      ],
    });
  });

  it('Should fill page with parseShowResult with a set of resultIDs', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    sessionStorage.setItem('ScanResult', JSON.stringify(mockResult));

    mockResult.forEach(async (result, index) => {
      const jsonkey = result.issue_id.toString() + result.result_id.toString();
      currentIssue = data[jsonkey];

      // Act
      await issue.openIssuePage(jsonkey, severity);
      const name = document.getElementsByClassName('issue-name')[0].innerHTML;
      const description = document.getElementById('information').nextElementSibling.innerHTML;
      const solution = document.getElementById('solution-text').innerHTML;

      // Assert
      test.value(name).isEqualTo(currentIssue.Name);
      test.value(description).isEqualTo(currentIssue.Information);
      test.value(htmlDecode(solution)).isEqualTo('1. ' + currentIssue.Solution[0]);
    });
  });
  it('parseShowResult fills the page with the correct structure for specific results', async function() {
    // Arrange
    // expectedFindings should be changed if the structure for specific results is changed in the code
    const expectedFindings = [
      '<li>process: p, port: 1, 2, 3</li><li>SYSTEM</li><li>CIS registry 1</li>' +
      '<li>SOFTWARE</li><li>CIS registry 2</li>',
      '<thead><tr><th>Process</th><th>Port(s)</th></tr></thead><tbody><tr><td style="width: 30%">p</td>\n' +
      '        <td style="width: 30%">1<br>2<br>3</td></tr></tbody>',
      'You changed your password on: process: p, port: 1, 2, 3You changed your password on: SYSTEM' +
      'You changed your password on: CIS registry 1You changed your password on: ' +
      'SOFTWAREYou changed your password on: CIS registry 2',
      '<tbody><tr><td style="width: 30%; word-break: break-all">SYSTEM</td></tr><tr>' +
      '<td style="width: 30%; word-break: break-all">' +
      'SOFTWARE</td></tr><tr><td style="width: 30%; word-break: break-all">undefined</td></tr></tbody>',
      '<tbody><tr><td style="width: 30%; word-break: break-all">SYSTEM</td>\n' +
      '        <td>CIS registry 1</td></tr><tr><td style="width: 30%; word-break: break-all">SOFTWARE</td>\n' +
      '        <td>CIS registry 2</td></tr></tbody>',
      '<tbody><tr><td style="width: 30%; word-break: break-all">SYSTEM</td></tr><tr><td style="width: 30%;' +
      ' word-break: break-all">SOFTWARE</td></tr><tr><td style="width: 30%; word-break: break-all">' +
      'undefined</td></tr></tbody>',
    ];

    // Assert
    await testParseShowResult('11', expectedFindings[0]);
    await testParseShowResult('21', expectedFindings[0]);
    await testParseShowResult('60', expectedFindings[0]);
    await testParseShowResult('70', expectedFindings[0]);
    await testParseShowResult('80', expectedFindings[0]);
    await testParseShowResult('90', expectedFindings[0]);
    await testParseShowResult('100', expectedFindings[0]);
    await testParseShowResult('110', expectedFindings[1]);
    await testParseShowResult('160', expectedFindings[2]);
    await testParseShowResult('173', expectedFindings[0]);
    await testParseShowResult('201', expectedFindings[0]);
    await testParseShowResult('230', expectedFindings[0]);
    await testParseShowResult('271', expectedFindings[3]);
    await testParseShowResult('311', expectedFindings[0]);
    await testParseShowResult('320', expectedFindings[4]);
    await testParseShowResult('351', expectedFindings[5]);
    await testParseShowResult('361', expectedFindings[5]);
  });

  /** helper function for testing the correct structure of parseShowResult
   *
   * @param {*} jsonkey key of issue being tested
   * @param {string} expectedFinding expected result found in the resultline part of parseShowResult
   */
  async function testParseShowResult(jsonkey, expectedFinding) {
    // Arrange
    const issue = await import('../src/js/issue.js');
    sessionStorage.setItem('ScanResult', JSON.stringify(mockResult));

    // Act
    await issue.openIssuePage(jsonkey, severity);
    let findings = '';
    if (jsonkey == 60 || jsonkey == 70 || jsonkey == 80 || jsonkey == 90 || jsonkey == 100) {
      findings = document.getElementById('description').nextElementSibling.innerHTML;
      expectedFinding = 'Issues.Permissions';
    } else if (jsonkey == 110) {
      findings = document.getElementById('description').nextElementSibling.innerHTML;
      expectedFinding = 'Issues.Port';
    } else if (jsonkey == 160) {
      findings = document.getElementById('description').nextElementSibling.innerHTML;
      expectedFinding = 'Issues.Password';
    } else if (jsonkey == 271 || jsonkey == 351|| jsonkey == 361) {
      findings = document.getElementById('description').nextElementSibling.innerHTML;
      expectedFinding = 'Issues.Cookies';
    } else if (jsonkey == 320) {
      findings = document.getElementsByClassName('issues-table')[0].innerHTML;
    } else {
      findings = document.getElementById('description').nextElementSibling.innerHTML;
    }
    // Assert
    test.value(findings).isEqualTo(expectedFinding);
  }

  it('parseShowResult keeps findings empty if the issueID is not in the issuesWithResultsShow list', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    // Mock scan result
    const mockResult = [
      {
        issue_id: 1,
        result_id: 0,
        result: [
          'findings',
        ],
      },
    ];
    sessionStorage.setItem('ScanResult', JSON.stringify(mockResult));
    const jsonkey = mockResult[0].issue_id.toString() + mockResult[0].result_id.toString();
    currentIssue = data[jsonkey];
    const pageContents = document.getElementById('page-contents');

    // Act
    pageContents.innerHTML = issue.parseShowResult(jsonkey, currentIssue);
    const findings = document.getElementById('description').innerHTML;

    // Assert
    test.value(findings).isEqualTo('');
  });
  it('checkShowResult should check if an issue name contains "applications with"', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');

    // Act
    const checked = issue.checkShowResult(data[60]);

    // Assert
    test.value(checked).isEqualTo(true);
  });
  it('getVersionScreenshot returns the right screenshot for the detected windows version', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    let testIssue = data['11'];

    // Act
    // clear sessionstorage
    sessionStorage.removeItem('WindowsVersion');
    let result = issue.getVersionScreenshot(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.Screenshots[0]);

    // Act
    sessionStorage.setItem('WindowsVersion', '10');
    result = issue.getVersionScreenshot(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.ScreenshotsWindows10[0]);

    // Act
    sessionStorage.setItem('WindowsVersion', '11');
    result = issue.getVersionScreenshot(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.Screenshots[0]);

    // Act
    sessionStorage.setItem('WindowsVersion', '10');
    testIssue = data['30'];
    result = issue.getVersionScreenshot(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.Screenshots[0]);

    // Act
    testIssue = data['310'];
    result = issue.getVersionScreenshot(testIssue, 0);

    // Assert
    test.value(result).isEqualTo('');
  });
  it('getVersionSolution returns the right solution for the detected windows version', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    let testIssue = data['11'];

    // Act
    // clear sessionstorage
    sessionStorage.removeItem('WindowsVersion');
    let result = issue.getVersionSolution(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.Solution[0]);

    // Act
    sessionStorage.setItem('WindowsVersion', '10');
    result = issue.getVersionSolution(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.SolutionWindows10[0]);

    // Act
    sessionStorage.setItem('WindowsVersion', '11');
    result = issue.getVersionSolution(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.Solution[0]);

    // Act
    sessionStorage.setItem('WindowsVersion', '10');
    testIssue = data['30'];
    result = issue.getVersionSolution(testIssue, 0);

    // Assert
    test.value(result).isEqualTo(testIssue.Solution[0]);

    // Act
    testIssue = data['310'];
    result = issue.getVersionSolution(testIssue, 0);

    // Assert
    test.value(result).isEqualTo('');
  });
});


