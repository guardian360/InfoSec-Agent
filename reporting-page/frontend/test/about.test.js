import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, storageMock} from './mock.js';

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
}));

// Mock sessionStorage
global.sessionStorage = storageMock;

describe('About page', function() {
  it('openAboutPage opens the about page', async function() {
    // Arrange
    const about = await import('../src/js/about.js');
    const classNames = [
      'lang-about-title',
      'lang-about-info',
      'lang-summary-title',
      'lang-summary-info',
      'lang-affiliations-title',
      'lang-affiliations-info',
      'lang-contributing-title',
      'lang-contributing-info',
    ];
    const expectedTexts = [
      'About.AboutTitle',
      'About.AboutInfo',
      'About.SummaryTitle',
      'About.SummaryInfo',
      'About.AffiliationsTitle',
      'About.AffiliationsInfo',
      'About.ContributingTitle',
      'About.ContributingInfo',
    ];

    // Act
    await about.openAboutPage();

    // Assert
    expectedTexts.forEach((expected, index) => {
      test.value(document.getElementsByClassName(classNames[index])[0].innerHTML).isEqualTo(expected);
    });
  });
});
