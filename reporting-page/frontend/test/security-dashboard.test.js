import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, mockChart,
  mockGraph, clickEvent, changeEvent, storageMock} from './mock.js';
import {RiskCounters} from '../src/js/risk-counters.js';

global.TESTING = true;

// Mock issue page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div id="page-contents"></div>
    <button id="security-button-applications"></button>
    <button id="security-button-devices"></button>
    <button id="security-button-network"></button>
    <button id="security-button-os"></button>
    <button id="security-button-passwords"></button>
    <button id="security-button-other"></button>
    <div id="dropbtn">Dropdown Button</div>
    <div id="myDropdown" class="dropdown-content show">Dropdown Content</div>
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

// Mock graph
mockGraph();

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

// Mock Tray
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
  ChangeLanguage: jest.fn(),
  ChangeScanInterval: jest.fn(),
  LogDebug: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issue.js', () => ({
  openIssuePage: jest.fn(),
  scrollToElement: jest.fn(),
}));

// Mock suggestedIssue
jest.unstable_mockModule('../src/js/home.js', () => ({
  suggestedIssue: jest.fn(),
}));

// Mock openPersonalizePage
jest.unstable_mockModule('../src/js/personalize.js', () => ({
  openPersonalizePage: jest.fn(),
  retrieveTheme: jest.fn(),
}));

// Mock openAllChecksPage
jest.unstable_mockModule('../src/js/all-checks.js', () => ({
  openAllChecksPage: jest.fn(),
}));

// test cases
describe('Security dashboard', function() {
  it('openSecurityDashboardPage should add the dashboard to the page-contents with graph functions', async function() {
    // Arrange
    const graph = await import('../src/js/graph.js');
    const Graph = new graph.Graph();

    const changeGraphSpy = jest.spyOn(Graph, 'changeGraph');
    const toggleRisksSpy = jest.spyOn(Graph, 'toggleRisks');
    const graphDropdownSpy = jest.spyOn(Graph, 'graphDropdown');

    const dashboard = await import('../src/js/security-dashboard.js');
    const rc = new RiskCounters(1, 1, 1, 1, 1);
    sessionStorage.setItem('SecurityRiskCounters', JSON.stringify(rc));

    // Act
    await dashboard.openSecurityDashboardPage();
    dashboard.addGraphFunctions(Graph);

    // Assert
    const status = document.getElementsByClassName('lang-security-stat')[0].innerHTML;
    test.value(status).isEqualTo('Dashboard.SecurityStatus');

    // Act
    document.getElementById('dropbtn').dispatchEvent(clickEvent);

    // Assert
    expect(graphDropdownSpy).toHaveBeenCalled();

    // Act
    document.getElementById('graph-interval').dispatchEvent(changeEvent);

    // Assert
    expect(changeGraphSpy).toHaveBeenCalled();

    // Act
    document.getElementById('select-high-risk').dispatchEvent(clickEvent);

    // Assert
    expect(toggleRisksSpy).toHaveBeenCalled();

    // Act
    document.getElementById('select-medium-risk').dispatchEvent(clickEvent);

    // Assert
    expect(toggleRisksSpy).toHaveBeenCalled();

    // Act
    document.getElementById('select-low-risk').dispatchEvent(clickEvent);

    // Assert
    expect(toggleRisksSpy).toHaveBeenCalled();

    // Act
    document.getElementById('select-info-risk').dispatchEvent(clickEvent);

    // Assert
    expect(toggleRisksSpy).toHaveBeenCalled();

    // Act
    document.getElementById('select-no-risk').dispatchEvent(clickEvent);

    // Assert
    expect(toggleRisksSpy).toHaveBeenCalled();
  });

  it('adjustWithRiskCounters should show data from risk counters', async function() {
    // arrange
    const mockRiskCounters = {
      lastInfoRisk: 0,
      lastHighRisk: 2,
      lastMediumRisk: 3,
      lastLowRisk: 4,
      lastNoRisk: 5,
      count: 5,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    // act
    dashboard.adjustWithRiskCounters(mockRiskCounters, global.document, false);

    // assert
    test.value(document.getElementById('high-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastHighRisk);
    test.value(document.getElementById('medium-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastMediumRisk);
    test.value(document.getElementById('low-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastLowRisk);
    test.value(document.getElementById('no-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastNoRisk);
    test.value(document.getElementById('info-risk-counter').innerHTML).isEqualTo(mockRiskCounters.lastInfoRisk);
  });

  it('Should display the right security status', async function() {
    // arrange
    const expectedColors = ['rgb(255, 255, 255)', 'rgb(255, 255, 255)', 'rgb(0, 0, 0)', 'rgb(0, 0, 0)'];

    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',

      lastHighRisk: 10,
      lastMediumRisk: 10,
      lastLowRisk: 10,
      lastNoRisk: 10,
      lastInfoRisk: 10,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    expectedColors.forEach((element, index) => {
      // act
      dashboard.adjustWithRiskCounters(mockRiskCounters, dom.window.document, false);

      // assert
      test.value(dom.window.document.getElementById('high-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastHighRisk);
      test.value(dom.window.document.getElementById('medium-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastMediumRisk);
      test.value(dom.window.document.getElementById('low-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastLowRisk);
      test.value(dom.window.document.getElementById('no-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastNoRisk);
      test.value(dom.window.document.getElementById('info-risk-counter').innerHTML)
        .isEqualTo(mockRiskCounters.lastInfoRisk);
    });
  });

  it('adjustWithRiskCounters should display the right security status', async function() {
    // Arrange
    const expectedText = [
      'Dashboard.Critical',
      'Dashboard.MediumConcern',
      'Dashboard.LowConcern',
      'Dashboard.NoConcern',
      'Dashboard.NoConcern',
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
    sessionStorage.setItem('SecurityRiskCounters', JSON.stringify(mockRiskCounters));

    const dashboard = await import('../src/js/security-dashboard.js');

    dashboard.openSecurityDashboardPage();
    const securityStatus = document.getElementsByClassName('status-descriptor');
    expectedText.forEach(async (element, index) => {
      // Act
      if (index == 1) mockRiskCounters.lastHighRisk = 0;
      if (index == 2) mockRiskCounters.lastMediumRisk = 0;
      if (index == 3) mockRiskCounters.lastLowRisk = 0;
      await dashboard.adjustWithRiskCounters(mockRiskCounters, dom.window.document, false);

      // Assert
      test.value(securityStatus[0].innerHTML)
        .isEqualTo(element);
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

  it('setMaxInterval should set the max value of the graph interval to the maximum amount of data', async function() {
    // arrange
    const mockRiskCounters = {
      count: 5,
    };

    const dashboard = await import('../src/js/security-dashboard.js');

    // act
    dashboard.setMaxInterval(mockRiskCounters, dom.window.document);

    // assert
    test.value(dom.window.document.getElementById('graph-interval').max).isEqualTo(mockRiskCounters.count);
  });

  it('suggestedIssue should open the issue page of highest risk security issue', async function() {
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
    const suggestedIssueMockMock = jest.spyOn(home, 'suggestedIssue');

    // Assert
    button.dispatchEvent(clickEvent);
    expect(suggestedIssueMockMock).toHaveBeenCalled();
  });

  it('Clicking the buttons should call openAllChecksPage on the right place', async function() {
    // Arrange
    const buttonApp = document.getElementById('security-button-applications');
    const buttonDevices = document.getElementById('security-button-devices');
    const buttonNet = document.getElementById('security-button-network');
    const buttonOS = document.getElementById('security-button-os');
    const buttonPass = document.getElementById('security-button-passwords');
    const buttonOther = document.getElementById('security-button-other');
    const allChecks = await import('../src/js/all-checks.js');

    // Act
    buttonApp.dispatchEvent(clickEvent);
    buttonDevices.dispatchEvent(clickEvent);
    buttonNet.dispatchEvent(clickEvent);
    buttonOS.dispatchEvent(clickEvent);
    buttonPass.dispatchEvent(clickEvent);
    buttonOther.dispatchEvent(clickEvent);

    // Assert
    expect(allChecks.openAllChecksPage).toHaveBeenCalledWith('applications');
    expect(allChecks.openAllChecksPage).toHaveBeenCalledWith('devices');
    expect(allChecks.openAllChecksPage).toHaveBeenCalledWith('network');
    expect(allChecks.openAllChecksPage).toHaveBeenCalledWith('os');
    expect(allChecks.openAllChecksPage).toHaveBeenCalledWith('passwords');
    expect(allChecks.openAllChecksPage).toHaveBeenCalledWith('security-other');
  });
  it('Clicking outside the dropdown should close it', async function() {
    // Arrange
    const dropdown = document.getElementById('myDropdown');
    dropdown.classList.add('show'); // Make sure the dropdown is initially open

    // Act
    document.body.dispatchEvent(new dom.window.Event('click', {bubbles: true}));

    // Assert
    expect(dropdown.classList.contains('show')).toBe(false);
  });
  it('Clicking inside the dropdown should not close it', async function() {
    // Arrange
    const dropdown = document.getElementById('myDropdown');
    const dropbtn = document.getElementById('dropbtn');
    dropdown.classList.add('show'); // Make sure the dropdown is initially open

    // Act
    dropbtn.dispatchEvent(new dom.window.Event('click', {bubbles: true}));

    // Assert
    expect(dropdown.classList.contains('show')).toBe(true);
  });
});
