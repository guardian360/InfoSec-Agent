import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, mockChart, clickEvent, storageMock} from './mock.js';
import {RiskCounters} from '../src/js/risk-counters.js';

global.TESTING = true;

// Mock issue page
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
  GetImagePath: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issue.js', () => ({
  openIssuePage: jest.fn(),
  scrollToElement: jest.fn(),
}));

// Mock openPersonalizePage
jest.unstable_mockModule('../src/js/personalize.js', () => ({
  openPersonalizePage: jest.fn(),
  retrieveTheme: jest.fn(),
}));

// Mock suggestedIssue
jest.unstable_mockModule('../src/js/home.js', () => ({
  suggestedIssue: jest.fn(),
}));

// Mock Tray
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
  ChangeLanguage: jest.fn(),
  ChangeScanInterval: jest.fn(),
  LogDebug: jest.fn(),
}));

describe('Privacy dashboard page', function() {
  it('openPrivacyDashboardPage should add the dashboard to the page-contents', async function() {
    // Arrange
    const dashboard = await import('../src/js/privacy-dashboard.js');
    const rc = new RiskCounters(1, 1, 1, 1, 1);
    sessionStorage.setItem('PrivacyRiskCounters', JSON.stringify(rc));

    // Act
    await dashboard.openPrivacyDashboardPage();

    // Assert
    test.value(document.getElementsByClassName('lang-privacy-stat')[0].innerHTML).isEqualTo('Dashboard.PrivacyStatus');
  });
  it('adjustWithRiskCounters should show the correct style', async function() {
    // arrange
    const expectedBackgroundColors = [
      'rgb(0, 255, 255)',
      'rgb(0, 0, 255)',
      'rgb(255, 0, 0)',
      'rgb(255, 255, 255)',
      'rgb(255, 255, 255)',
    ];
    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      infoColor: 'rgb(255, 255, 255)',
      noRiskColor: 'rgb(255, 255, 255)',

      lastHighRisk: 10,
      lastMediumRisk: 10,
      lastLowRisk: 10,
      lastInfoRisk: 10,
      lastnoRisk: 10,

      allHighRisks: [10],
      allMediumRisks: [10],
      allLowRisks: [10],
      allNoRisks: [10],
      allInfoRisks: [10],
    };
    sessionStorage.setItem('PrivacyRiskCounters', JSON.stringify(mockRiskCounters));

    const pDashboard = await import('../src/js/privacy-dashboard.js');
    const sDashboard = await import('../src/js/security-dashboard.js');

    pDashboard.openPrivacyDashboardPage();
    const securityStatus = document.getElementsByClassName('status-descriptor')[0];
    expectedBackgroundColors.forEach(async (element, index) => {
      // Act
      if (index == 1) mockRiskCounters.lastHighRisk = 0;
      if (index == 2) mockRiskCounters.lastMediumRisk = 0;
      if (index == 3) mockRiskCounters.lastLowRisk = 0;
      if (index == 4) mockRiskCounters.noRiskColor = 0;
      sDashboard.adjustWithRiskCounters(mockRiskCounters, dom.window.document, false);

      // Assert
      test.value(securityStatus.style.backgroundColor)
        .isEqualTo(expectedBackgroundColors[index]);
    });
  });
  it('Clicking the scan-now button should call scanTest', async function() {
    // Arrange
    const database = await import('../src/js/database.js');
    const scanTestMock = jest.spyOn(database, 'scanTest');
    const scanButton = document.getElementById('scan-now');

    // Act
    scanButton.dispatchEvent(clickEvent);

    // Assert
    expect(scanTestMock).toHaveBeenCalled();
  });
  it('suggestedIssue should open the issue page of highest risk privacy issue', async function() {
    // Arrange
    let issues = [];
    issues = [
      {id: 1, severity: 4, jsonkey: 10},
      {id: 5, severity: 1, jsonkey: 51},
      {id: 15, severity: 0, jsonkey: 150},
      {id: 4, severity: 2, jsonkey: 41},
    ];
    sessionStorage.setItem('DataBaseData', JSON.stringify(issues));

    const home = await import('../src/js/home.js');
    const button = document.getElementById('suggested-issue');
    const suggestedIssueMock = jest.spyOn(home, 'suggestedIssue');

    // Assert
    button.dispatchEvent(clickEvent);

    expect(suggestedIssueMock).toHaveBeenCalled();
  });
});
