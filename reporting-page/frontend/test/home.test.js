import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, clickEvent, storageMock, mockChart, scanResultMock} from './mock.js';
import {RiskCounters} from '../src/js/risk-counters.js';
import data from '../src/databases/database.en-GB.json' assert { type: 'json' };

global.TESTING = true;

// Mock home page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
  <div id="logo-button" class="logo-name">
    <img id="logo" alt="logo" src="">
    <div class="header-name">
      <h1 id="title">InfoSec Agent</h1>
    </div>
  </div>
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
  LoadUserSettings: jest.fn().mockImplementation((input) => 10),
  GetImagePath: jest.fn().mockImplementation((input) => input),
  GetLighthouseState: jest.fn().mockImplementationOnce(() => 0)
    .mockImplementationOnce(() => 1)
    .mockImplementationOnce(() => 2)
    .mockImplementationOnce(() => 3)
    .mockImplementationOnce(() => 4)
    .mockImplementationOnce(() => 5)
    .mockImplementation(() => 0),
}));

// Mock LogError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
  ChangeLanguage: jest.fn(),
  ChangeScanInterval: jest.fn(),
  LogDebug: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issue.js', () => ({
  openIssuePage: jest.fn(),
}));

// Mock settings
jest.unstable_mockModule('../src/js/settings.js', () => ({
  showModal: jest.fn(),
}));

const socialMediaSizesMock = {
  facebook: {
    name: 'facebook',
    height: 315,
    width: 600,
  },
};

// Mock share
jest.unstable_mockModule('../src/js/share.js', () => ({
  setImage: jest.fn(),
  saveProgress: jest.fn(),
  shareProgress: jest.fn(),
  selectSocialMedia: jest.fn(),
  socialMediaSizes: socialMediaSizesMock,
}));

// Mock openPersonalizePage
jest.unstable_mockModule('../src/js/personalize.js', () => ({
  openPersonalizePage: jest.fn(),
  retrieveTheme: jest.fn(),
}));

describe('Home page', function() {
  it('openHomePage should add the home page to the page-contents and load the correct background', async function() {
    // Arrange
    const homepage = await import('../src/js/home.js');
    const rc = new RiskCounters(1, 1, 1, 1, 1);
    sessionStorage.setItem('RiskCounters', JSON.stringify(rc));
    const backgroundVideos = [
      'gamification/state0.mkv',
      'gamification/state1.mkv',
      'gamification/state2.mkv',
      'gamification/state3.mkv',
      'gamification/state4.mkv',
      'gamification/state0.mkv',
    ];

    // Act
    await homepage.openHomePage();
    let background = document.getElementById('lighthouse-background').src;

    // Assert
    test.value(
      document.getElementsByClassName('lang-piechart-header')[0].innerHTML).isEqualTo('Dashboard.RiskLevelDistribution',
    );
    test.value(background).isEqualTo(backgroundVideos[0]);

    // Act
    await homepage.openHomePage();
    background = document.getElementById('lighthouse-background').src;

    // Assert
    test.value(background).isEqualTo(backgroundVideos[1]);
    // Act
    await homepage.openHomePage();
    background = document.getElementById('lighthouse-background').src;

    // Assert
    test.value(background).isEqualTo(backgroundVideos[2]);
    // Act
    await homepage.openHomePage();
    background = document.getElementById('lighthouse-background').src;

    // Assert
    test.value(background).isEqualTo(backgroundVideos[3]);
    // Act
    await homepage.openHomePage();
    background = document.getElementById('lighthouse-background').src;

    // Assert
    test.value(background).isEqualTo(backgroundVideos[4]);

    // Act
    await homepage.openHomePage();
    background = document.getElementById('lighthouse-background').src;

    // Assert
    test.value(background).isEqualTo(backgroundVideos[5]);
  });
  it('suggestedIssue should open the issue page of highest risk issue', async function() {
    // Arrange
    // begin with a level 4 severity, to test that it is skipped
    const scanResultBegin = [
      {
        issue_id: 10,
        result_id: 0,
        result: [],
      },
    ];
    sessionStorage.setItem('ScanResult', JSON.stringify(scanResultBegin.concat(scanResultMock)));

    const issue = await import('../src/js/issue.js');
    const button = document.getElementById('suggested-issue');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Assert
    button.dispatchEvent(clickEvent);
    expect(openIssuePageMock).toHaveBeenCalled();
  });
  it('suggestedIssue is able to open either a security or privacy issue of highest risk', async function() {
    // Arrange
    const home = await import('../src/js/home.js');
    const scanResultBegin = [
      {
        issue_id: 4,
        result_id: 1,
        result: [],
      },
    ];
    sessionStorage.setItem('ScanResult', JSON.stringify(scanResultBegin.concat(scanResultMock)));
    const scanResult = JSON.parse(sessionStorage.getItem('ScanResult'));

    const issue = await import('../src/js/issue.js');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Act
    home.suggestedIssue('Security');
    const securityIssue = scanResult[4];

    // Assert
    expect(openIssuePageMock).toHaveBeenCalledWith(securityIssue.issue_id, securityIssue.result_id, 'home');

    // Act
    home.suggestedIssue('Privacy');
    const privacyIssue = scanResult[0];

    // Assert
    expect(openIssuePageMock).toHaveBeenCalledWith(privacyIssue.issue_id, privacyIssue.result_id, 'home');
  });
  it('getSeverity has a correct return value', async function() {
    // Arrange
    const home = await import('../src/js/home.js');

    // Act
    let severity = home.getSeverity(0, 0);

    // Assert
    test.value(severity).isUndefined();

    // Act
    severity = home.getSeverity(1, 10);

    // Assert
    test.value(severity).isUndefined();

    // Act
    severity = home.getSeverity(1, 0);

    // Assert
    test.value(severity).isEqualTo(data[1][0].Severity);
  });
  it('clicking on the share modal buttons calls the correct functions', async function() {
    // Arrange
    const settings = await import('../src/js/settings.js');
    const showModalMock = jest.spyOn(settings, 'showModal');
    const share = await import('../src/js/share.js');
    const saveProgressMock = jest.spyOn(share, 'saveProgress');
    const shareProgressMock = jest.spyOn(share, 'shareProgress');
    const selectSocialMediaMock = jest.spyOn(share, 'selectSocialMedia');

    // Act
    document.getElementById('share-progress').dispatchEvent(clickEvent);
    document.getElementById('share-save-button').dispatchEvent(clickEvent);
    document.getElementById('share-button').dispatchEvent(clickEvent);

    document.getElementById('select-facebook').dispatchEvent(clickEvent);
    document.getElementById('select-x').dispatchEvent(clickEvent);
    document.getElementById('select-linkedin').dispatchEvent(clickEvent);
    document.getElementById('select-instagram').dispatchEvent(clickEvent);

    // Assert
    expect(showModalMock).toHaveBeenCalled();
    expect(saveProgressMock).toHaveBeenCalled();
    expect(shareProgressMock).toHaveBeenCalled();
    expect(selectSocialMediaMock).toHaveBeenCalledTimes(4);
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
  it('The logo and title are loaded when the window is loaded', async function() {
    // Arrange
    localStorage.setItem('picture', 'picturePath');
    localStorage.setItem('title', 'titlePath');

    // Act
    window.dispatchEvent(new Event('load'));

    // Assert
    test.value(document.getElementById('logo').src).isEqualTo('picturePath');
    test.value(document.getElementById('title').textContent).isEqualTo('titlePath');
  });
});
