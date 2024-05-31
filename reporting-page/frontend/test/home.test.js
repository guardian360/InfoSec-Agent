import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, clickEvent, storageMock, mockChart} from './mock.js';
import {RiskCounters} from '../src/js/risk-counters.js';
// import {suggestedIssue} from '../src/js/home.js';

global.TESTING = true;

// Mock home page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div id="page-contents"></div>
</body>
</html>
`);
global.document = dom.window.document;
global.window = dom.window;

// Mock sessionStorage
global.sessionStorage = storageMock;
global.localStorage = storageMock;

// Mock often used page functions
mockPageFunctions();

// Mock chart constructor
mockChart();

// Mock scanTest
jest.unstable_mockModule('../src/js/database.js', () => ({
  scanTest: jest.fn(),
}));

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
  LoadUserSettings: jest.fn(),
}));

// Mock LogError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
  ChangeLanguage: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issue.js', () => ({
  openIssuePage: jest.fn(),
}));

// Mock openPersonalizePage
jest.unstable_mockModule('../src/js/personalize.js', () => ({
  openPersonalizePage: jest.fn(),
  retrieveTheme: jest.fn(),
}));

describe('Home page', function() {
  it('openHomePage should add the home page to the page-contents', async function() {
    // Arrange
    const homepage = await import('../src/js/home.js');
    const rc = new RiskCounters(1, 1, 1, 1, 1);
    sessionStorage.setItem('RiskCounters', JSON.stringify(rc));

    // Act
    await homepage.openHomePage();

    // Assert
    test.value(
      document.getElementsByClassName('lang-piechart-header')[0].innerHTML).isEqualTo('Dashboard.RiskLevelDistribution',
    );
  });
  it('suggestedIssue should open the issue page of highest risk issue', async function() {
    // Arrange
    let issues = [];
    issues = [
      {id: 1, severity: 4, jsonkey: 10},
      {id: 5, severity: 1, jsonkey: 51},
      {id: 15, severity: 0, jsonkey: 150},
      {id: 4, severity: 2, jsonkey: 41},
    ];
    sessionStorage.setItem('DataBaseData', JSON.stringify(issues));

    const issue = await import('../src/js/issue.js');
    const button = document.getElementById('suggested-issue');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Assert
    button.dispatchEvent(clickEvent);
    expect(openIssuePageMock).toHaveBeenCalled();
  });
});
