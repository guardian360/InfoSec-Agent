import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, clickEvent, storageMock, mockChart, scanResultMock} from './mock.js';
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
  GetImagePath: jest.fn(),
  GetLighthouseState: jest.fn(),
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
    sessionStorage.setItem('ScanResult', JSON.stringify(scanResultMock));

    const issue = await import('../src/js/issue.js');
    const button = document.getElementById('suggested-issue');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Assert
    button.dispatchEvent(clickEvent);
    expect(openIssuePageMock).toHaveBeenCalled();
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
});
