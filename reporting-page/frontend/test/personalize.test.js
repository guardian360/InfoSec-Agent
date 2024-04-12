import 'jsdom-global/register.js';
import {
  handleFaviconChange,
  handlePictureChange,
  handleTitleChange} from '../src/js/personalize.js';
import {JSDOM} from 'jsdom';
import test from 'unit.js';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
  <input type="file" id="picture-input">
  <img id="logo" src="">
  <input type="text" id="newTitle" value="New Page Title">
  <h1 id="title">Old Title</h1>
</body>
</html>
`);
global.document = dom.window.document;
global.window = dom.window;

// Mock FileReader
const fileReaderIco = global.FileReader = class {
  /**
 * Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
 * This method is used for testing purposes to mimic the behavior of FileReader.
 * @param {Function} onload - The onload event handler to be triggered with a mock result.
 */
  readAsDataURL() {
    this.onload({target: {result: 'data:image/x-icon'}});
  }
};
const fileReaderPng = global.FileReader = class { /**
* Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
* This method is used for testing purposes to mimic the behavior of FileReader.
* @param {Function} onload - The onload event handler to be triggered with a mock result.
*/
  readAsDataURL() {
    this.onload({target: {result: 'data:image/png'}});
  }
};
const fileReaderJpg = global.FileReader = class { /**
* Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
* This method is used for testing purposes to mimic the behavior of FileReader.
* @param {Function} onload - The onload event handler to be triggered with a mock result.
*/
  readAsDataURL() {
    this.onload({target: {result: 'data:image/jpg'}});
  }
};
const fileReaderJpeg = global.FileReader = class { /**
* Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
* This method is used for testing purposes to mimic the behavior of FileReader.
* @param {Function} onload - The onload event handler to be triggered with a mock result.
*/
  readAsDataURL() {
    this.onload({target: {result: 'data:image/jpeg'}});
  }
};

// Mock localStorage
const localStorageMock = (() => {
  let store = {};

  return {
    getItem: (key) => store[key],
    setItem: (key, value) => {
      store[key] = value.toString();
    },
    clear: () => {
      store = {};
    },
  };
})();
global.localStorage = localStorageMock;


// Test cases
describe('handleFaviconSelect', () => {
  it('should change the favicon when a valid .ico file is selected', () => {
    // Arrange
    const head = document.querySelector('head');

    // Act
    FileReader = fileReaderIco;
    const blob = new Blob(['dummy'], {type: 'image/x-icon'});
    handleFaviconChange({target: {files: blob}});

    // Assert
    const newFavicon = head.querySelector('link[rel="icon"]');
    test.value(newFavicon.href).isEqualTo('data:image/x-icon');
  });

  it('should change the favicon when a valid .png file is selected', () => {
    // Arrange
    const head = document.querySelector('head');

    // Act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    handleFaviconChange({target: {files: blob}});

    // Assert
    const newFavicon = head.querySelector('link[rel="icon"]');
    test.value(newFavicon.href).isEqualTo('data:image/png');
  });
  it('saves valid .ico favicon in localStorage', () => {
    // Arrange
    // Act
    FileReader = fileReaderIco;
    const blob = new Blob(['dummy'], {type: 'image/x-icon'});
    handleFaviconChange({target: {files: blob}});

    // Assert
    const favicon = localStorageMock.getItem('favicon');
    const expectedValue = 'data:image/x-icon';
    test.value(favicon).isEqualTo(expectedValue);
  });
  it('saves valid .png favicon in localStorage', () => {
    // Arrange
    // Act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    handleFaviconChange({target: {files: blob}});

    // Assert
    test.value(localStorageMock.getItem('favicon')).isEqualTo('data:image/png');
  });
});

describe('handlePictureChange', () => {
  it('changes navigation picture with valid .png file', () => {
    // Arrange
    const logo = document.getElementById('logo');

    // Act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    handlePictureChange({target: {files: blob}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/png');
  });
  it('changes navigation picture with valid .jpg file', () => {
    // Arrange
    const logo = document.getElementById('logo');

    // Act
    FileReader = fileReaderJpg;
    const blob = new Blob(['dummy'], {type: 'image/jpg'});
    handlePictureChange({target: {files: blob}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpg');
  });
  it('changes navigation picture with valid .jpeg file', () => {
    // Arrange
    const logo = document.getElementById('logo');

    // Act
    FileReader = fileReaderJpeg;
    const blob = new Blob(['dummy'], {type: 'image/jpeg'});
    handlePictureChange({target: {files: blob}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpeg');
  });
  it('saves valid .png file in localStorage', () => {
    // Arrange
    // act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    handlePictureChange({target: {files: blob}});

    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/png');
  });
  it('saves valid .jpg file in localStorage', () => {
    // Arrange
    // act
    FileReader = fileReaderJpg;

    const blob = new Blob(['dummy'], {type: 'image/jpg'});
    handlePictureChange({target: {files: blob}});
    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/jpg');
  });
  it('saves .jpeg file in localStorage', () => {
    // Arrange
    // act
    FileReader = fileReaderJpeg;
    const blob = new Blob(['dummy'], {type: 'image/jpeg'});
    handlePictureChange({target: {files: blob}});

    // Assert
    const localStoragePicture = localStorageMock.getItem('picture');
    const expectedValue = 'data:image/jpeg';
    test.value(localStoragePicture).isEqualTo(expectedValue);
  });
});

describe('handleTitleChange', () => {
  it('should change the title of the page', () => {
    // Arrange
    const newTitleInput = document.getElementById('newTitle');
    const titleElement = document.getElementById('title');

    // Act
    handleTitleChange();

    // Assert
    test.value(titleElement.textContent).isEqualTo(newTitleInput.value);
  });

  it('should save the new title to localStorage', () => {
    // Arrange
    const newTitleInput = document.getElementById('newTitle');

    // Act
    handleTitleChange();

    // Assert
    const localStorageTitle = localStorageMock.getItem('title');
    const expectedValue = newTitleInput.value;
    test.value(localStorageTitle).isEqualTo(expectedValue);
  });
});
