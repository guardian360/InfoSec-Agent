import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions,
  mockGetLocalization,
  storageMock,
  mockScanNowGo,
  mockRiskCounters,
  mockGetDataBaseData,
  scanResultMock,
  mockOpenPageFunctions} from './mock.js';

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

// Mock RiskCounters
mockRiskCounters();

// Mock open page functions
mockOpenPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
}));

// Mock runtime functions
jest.unstable_mockModule('../wailsjs/runtime/runtime.js', () => ({
  WindowShow: jest.fn(),
  WindowMaximise: jest.fn(),
  LogPrint: jest.fn(),
}));

// Mock scanNowGo
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  ScanNow: jest.fn().mockImplementationOnce(() => mockScanNowGo(true))
    .mockImplementationOnce(() => mockScanNowGo(true))
    .mockImplementation(() => mockScanNowGo(false)),
  LogError: jest.fn(),
}));

// Mock session and localStorage
global.sessionStorage = storageMock;
global.localStorage = storageMock;

describe('database functions', function() {
  it('scanTest is called which calls scanNowGo and fills sessionstorage with data', async function() {
    // Arrange
    await import('../src/js/database.js');

    // Act
    const scanResult = JSON.parse(sessionStorage.getItem('ScanResult'));
    const called = sessionStorage.getItem('scanTest');
    const rc = JSON.parse(sessionStorage.getItem('RiskCounters'));
    const src = JSON.parse(sessionStorage.getItem('SecurityRiskCounters'));
    const prc = JSON.parse(sessionStorage.getItem('PrivacyRiskCounters'));

    // Assert
    test.array(scanResult).is(scanResultMock);
    test.value(called).isEqualTo('called');

    // risk counters have the right values
    // if amount is not correct, first change scanResultMock in mock.js
    // to correctly return issues covering every severity once
    test.array(rc.high).is([1]);
    test.array(rc.medium).is([1]);
    test.array(rc.low).is([1]);
    test.array(rc.info).is([1]);
    test.array(rc.acceptable).is([1]);

    test.array(src.high).is([1]);
    test.array(src.medium).is([0]);
    test.array(src.low).is([1]);
    test.array(src.info).is([0]);
    test.array(src.acceptable).is([0]);

    test.array(prc.high).is([0]);
    test.array(prc.medium).is([1]);
    test.array(prc.low).is([0]);
    test.array(prc.info).is([1]);
    test.array(prc.acceptable).is([1]);
  });
  it('calling scanTest again fill sessionstorage with additional data', async function() {
    // Arrange
    const database = await import('../src/js/database.js');
    sessionStorage.setItem('isScanning', 'false');
    // Act
    await database.scanTest(false);
    const rc = JSON.parse(sessionStorage.getItem('RiskCounters'));
    const src = JSON.parse(sessionStorage.getItem('SecurityRiskCounters'));
    const prc = JSON.parse(sessionStorage.getItem('PrivacyRiskCounters'));

    // Assert
    test.array(rc.high).is([1, 1]);
    test.array(rc.medium).is([1, 1]);
    test.array(rc.low).is([1, 1]);
    test.array(rc.info).is([1, 1]);
    test.array(rc.acceptable).is([1, 1]);

    test.array(src.high).is([1, 1]);
    test.array(src.medium).is([0, 0]);
    test.array(src.low).is([1, 1]);
    test.array(src.info).is([0, 0]);
    test.array(src.acceptable).is([0, 0]);

    test.array(prc.high).is([0, 0]);
    test.array(prc.medium).is([1, 1]);
    test.array(prc.low).is([0, 0]);
    test.array(prc.info).is([1, 1]);
    test.array(prc.acceptable).is([1, 1]);
  });
  it('errors are caught when calling scanNowGo', async function() {
    // Arrange
    const database = await import('../src/js/database.js');
    const tray = await import('../wailsjs/go/main/Tray.js');
    const logErrorMock = jest.spyOn(tray, 'LogError');

    // Act
    await database.scanTest(false);

    // Assert
    expect(logErrorMock).toHaveBeenCalled();

    // remove the called flag from sessionstorage
    sessionStorage.clear();
  });
  it('if scanTest is already called, loading the database.js file should not call scanTest', async function() {
    // Arrange
    sessionStorage.setItem('scanTest', 'called');
    await import('../src/js/database.js');

    // Act
    const scanResult = sessionStorage.getItem('ScanResult');
    const called = sessionStorage.getItem('scanTest');
    const rc = sessionStorage.getItem('RiskCounters');

    // Assert
    test.value(scanResult).isUndefined();
    test.value(called).isEqualTo('called');
    test.value(rc).isUndefined();
  });
});
