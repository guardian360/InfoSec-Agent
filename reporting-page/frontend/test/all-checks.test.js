import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, storageMock, clickEvent} from './mock.js';

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

// Mock often used page functions
mockPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
  LoadUserSettings: jest.fn(),
  GetImagePath: jest.fn(),
}));

// Mock runtime functions
jest.unstable_mockModule('../wailsjs/runtime/runtime.js', () => ({
  LogPrint: jest.fn(),
  WindowMaximise: jest.fn(),
  WindowShow: jest.fn(),
}));

// Mock LogError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
  LogDebug: jest.fn(),
  ChangeLanguage: jest.fn(),
  ChangeScanInterval: jest.fn(),
  ScanNow: jest.fn(),
}));

// Mock personalize
jest.unstable_mockModule('../src/js/personalize.js', () => ({
  retrieveTheme: jest.fn(),
  openPersonalizePage: jest.fn(),
}));

// Mock issue.js
jest.unstable_mockModule('../src/js/issue.js', () => ({
  scrollToElement: jest.fn(),
  openIssuePage: jest.fn(),
}));

// Mock sessionStorage
global.sessionStorage = storageMock;

describe('Checks page', function() {
  it('openAllChecksPage opens the checks page with all checks', async function() {
    // Arrange
    const allChecks = await import('../src/js/all-checks.js');
    const mockScanResult = [{issue_id: 3, result_id: 1, result: []},
      {issue_id: 4, result_id: 1, result: []},
      {issue_id: 5, result_id: 1, result: []}];

    // Act
    sessionStorage.setItem('ScanResult', JSON.stringify(mockScanResult));
    await allChecks.openAllChecksPage();
    const foundChecks = document.getElementsByClassName('all-checks-check');

    // Assert
    test.value(foundChecks.length).isEqualTo(42);
  });
  it('Checks found in the database have a link to an openIssuePage function', async function() {
    // Arrange
    const allChecks = await import('../src/js/all-checks.js');
    const mockScanResult = [{issue_id: 3, result_id: 1, result: []},
      {issue_id: 4, result_id: 1, result: []},
      {issue_id: 5, result_id: 1, result: []}];
    const issue = await import('../src/js/issue.js');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Act
    sessionStorage.setItem('ScanResult', JSON.stringify(mockScanResult));
    await allChecks.openAllChecksPage();

    mockScanResult.forEach((issue) => {
      document.getElementById(issue.issue_id).dispatchEvent(clickEvent);
    });

    // Assert
    expect(openIssuePageMock).toHaveBeenCalledTimes(mockScanResult.length);
  });
  it('getViewedElement returns the right node', async function() {
    // Arrange
    const allChecks = await import('../src/js/all-checks.js');
    const list = [
      'applications',
      'devices',
      'network',
      'os',
      'passwords',
      'security-other',
      'permissions',
      'browser',
      'privacy-other',
      'top',
    ];

    list.forEach((area) => {
      // Act
      const node = allChecks.getViewedElement(area);

      // Assert
      test.value(node.id).isEqualTo(area);
    });
  });
});
