import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, storageMock} from './mock.js';

global.TESTING = true;

// Mock issue page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div id="page-contents">
        <div id="share-modal" class="modal">
            <div class="modal-content">
              <div class="modal-header">
                <span id="close-share-modal" class="close">&times;</span>
                <p>Select where to share your progress, Save and download it, then share it with others!</p>
              </div>
              <div id="share-node" class="modal-body">
                <img class="api-key-image" src="https://placehold.co/600x315" alt="Step 1 Image">
              </div>
              <div id="share-buttons" class="modal-body">
                <a id="share-save-button" class="modal-button share-button">Save</a>
                <a class="share-button-break">|</a>
                <a id="select-facebook" class="select-button selected">Facebook</a>
                <a id="select-x" class="select-button">X</a>
                <a id="select-linkedin" class="select-button">LinkedIn</a>
                <a id="select-instagram" class="select-button">Instagram</a>
                <a class="share-button-break">|</a>
                <a id="share-button" class="modal-button share-button">Share</a>
              </div>
            </div>
        </div>
    </div>
</body>
</html>
`);
global.document = dom.window.document;
global.window = dom.window;

// mock createObjectURL
window.URL.createObjectURL = jest.fn().mockImplementation((input) => input);

// mock window.open
window.open = jest.fn();

// Mock sessionStorage
global.sessionStorage = storageMock;
global.localStorage = storageMock;

// Mock often used page functions
mockPageFunctions();

// Mock Chart constructor
jest.unstable_mockModule('html-to-image', () => ({
  toBlob: jest.fn().mockImplementation((node, config) => {
    return node.innerHTML + '_' + config.width.toString() + '_' + config.height.toString();
  }),
}));

jest.unstable_mockModule('browser-image-compression', () => ({
  imageCompression: jest.fn().mockImplementation((input, i) => input),
  default: jest.fn().mockImplementation((input) => input),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issues.js', () => ({
  getUserSettings: jest.fn().mockImplementationOnce(() => 2),
}));


describe('share functions', function() {
  beforeAll(() => {
    jest.useFakeTimers('modern');
    jest.setSystemTime(new Date(2000, 5, 1));
  });

  afterAll(() => {
    jest.useRealTimers();
  });


  it('getImage should return a url of the passed node converted to an image', async function() {
    // Arrange
    const share = await import('../src/js/share.js');

    // Act
    const node = document.getElementById('share-node');
    const url = await share.getImage(node, 600, 315);

    // Assert
    test.value(url).isEqualTo(node.innerHTML + '_600_315');
  });
  it('saveProgress should get the image from the html node passed and download it', async function() {
    // Arrange
    const share = await import('../src/js/share.js');

    const linkElement = {
      download: '',
      href: '',
      click: jest.fn(),
    };
    jest.spyOn(document, 'createElement').mockImplementation(() => linkElement);

    // Act
    const node = document.getElementById('share-node');
    await share.saveProgress(node);

    // Assert
    test.value(linkElement.download).isEqualTo('Info-Sec-Agent_6-1-2000_facebook.png');
    test.value(linkElement.href).isEqualTo(node.innerHTML + '_600_315');
    expect(linkElement.click).toHaveBeenCalled();
  });
  it('shareProgress should call window.open to the selected social media page', async function() {
    // Arrange
    const share = await import('../src/js/share.js');
    jest.spyOn(window, 'open');

    // Act
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(1);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['x']));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(2);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['linkedin']));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(3);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(share.socialMediaSizes['instagram']));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(4);

    // Act
    sessionStorage.setItem('ShareSocial', JSON.stringify(''));
    share.shareProgress();

    // Assert
    expect(window.open).toHaveBeenCalledTimes(4);
  });
  it('selectSocialMedia should select the right social media and set it in the session storage', async function() {
    // Arrange
    const share = await import('../src/js/share.js');
    const socialMedias = ['facebook', 'x', 'linkedin', 'instagram'];

    socialMedias.forEach((social) => {
      // Act
      share.selectSocialMedia(social);

      // Assert
      // The right social media box is selected
      socialMedias.forEach((social2) => {
        if (social == social2) {
          test.value(document.getElementById('select-' + social2).classList.contains('selected')).isTrue();
        } else test.value(document.getElementById('select-' + social2).classList.contains('selected')).isFalse();
      });

      const inStorage = JSON.parse(sessionStorage.getItem('ShareSocial'));
      test.value(inStorage.name).isEqualTo(share.socialMediaSizes[social].name);
      test.value(inStorage.height).isEqualTo(share.socialMediaSizes[social].height);
      test.value(inStorage.width).isEqualTo(share.socialMediaSizes[social].width);
    });

    // Act

    // Assert
  });
});
