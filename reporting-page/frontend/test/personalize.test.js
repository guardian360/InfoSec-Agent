import {handleFaviconChange, handlePictureChange, handleTitleChange} from '../src/js/personalize.js'; // Assuming the function is in picture.js
import {JSDOM} from 'jsdom';
import test from 'unit.js';

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
  readAsDataURL() {
    this.onload({target: {result: 'data:image/x-icon'}});
  }
};
const fileReaderPng = global.FileReader = class {
  readAsDataURL() {
    this.onload({target: {result: 'data:image/png'}});
  }
};
const fileReaderJpg = global.FileReader = class {
  readAsDataURL() {
    this.onload({target: {result: 'data:image/jpg'}});
  }
};
const fileReaderJpeg = global.FileReader = class {
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
    handleFaviconChange({target: {files: [new Blob(['dummy'], {type: 'image/x-icon'})]}});

    // Assert
    const newFavicon = head.querySelector('link[rel="icon"]');
    test.value(newFavicon.href).isEqualTo('data:image/x-icon');
  });

  it('should change the favicon when a valid .png file is selected', () => {
    // Arrange
    const head = document.querySelector('head');

    // Act
    FileReader = fileReaderPng;
    handleFaviconChange({target: {files: [new Blob(['dummy'], {type: 'image/png'})]}});

    // Assert
    const newFavicon = head.querySelector('link[rel="icon"]');
    test.value(newFavicon.href).isEqualTo('data:image/png');
  });
  it('should save the favicon when a valid .ico file is selected in localstorage', () => {
    // Arrange
    // Act
    FileReader = fileReaderIco;
    handleFaviconChange({target: {files: [new Blob(['dummy'], {type: 'image/x-icon'})]}});

    // Assert
    test.value(localStorageMock.getItem('favicon')).isEqualTo('data:image/x-icon');
  });
  it('should save the favicon when a valid .png file is selected in localstorage', () => {
    // Arrange
    // Act
    FileReader = fileReaderPng;
    handleFaviconChange({target: {files: [new Blob(['dummy'], {type: 'image/png'})]}});

    // Assert
    test.value(localStorageMock.getItem('favicon')).isEqualTo('data:image/png');
  });
});

describe('handlePictureChange', () => {
  it('should change the navigation picture when a valid .png file is selected', () => {
    // Arrange
    const logo = document.getElementById('logo');

    // Act
    FileReader = fileReaderPng;
    handlePictureChange({target: {files: [new Blob(['dummy'], {type: 'image/png'})]}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/png');
  });
  it('should change the navigation picture when a valid .jpg file is selected', () => {
    // Arrange
    const logo = document.getElementById('logo');

    // Act
    FileReader = fileReaderJpg;
    handlePictureChange({target: {files: [new Blob(['dummy'], {type: 'image/jpg'})]}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpg');
  });
  it('should change the navigation picture when a valid .jpeg file is selected', () => {
    // Arrange
    const logo = document.getElementById('logo');

    // Act
    FileReader = fileReaderJpeg;
    handlePictureChange({target: {files: [new Blob(['dummy'], {type: 'image/jpeg'})]}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpeg');
  });
  it('should save the navigation picture when a valid .png file is selected in localstorage', () => {
    // Arrange
    // act
    FileReader = fileReaderPng;
    handlePictureChange({target: {files: [new Blob(['dummy'], {type: 'image/png'})]}});

    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/png');
  });
  it('should save the navigation picture when a valid .jpg file is selected in localstorage', () => {
    // Arrange
    // act
    FileReader = fileReaderJpg;
    handlePictureChange({target: {files: [new Blob(['dummy'], {type: 'image/jpg'})]}});

    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/jpg');
  });
  it('should save the navigation picture when a valid .jpeg file is selected in localstorage', () => {
    // Arrange
    // act
    FileReader = fileReaderJpeg;
    handlePictureChange({target: {files: [new Blob(['dummy'], {type: 'image/jpeg'})]}});

    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/jpeg');
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
    test.value(localStorageMock.getItem('title')).isEqualTo(newTitleInput.value);
  });
});
