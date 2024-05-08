import 'jsdom-global/register.js';
import {JSDOM} from 'jsdom';
import test from 'unit.js';
import {jest} from '@jest/globals';
import {fireEvent} from '@testing-library/dom';
global.TESTING = true;

// Mock page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
  <div id="page-contents"></div>
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
   */
  readAsDataURL() {
    this.onload({target: {result: 'data:image/x-icon'}});
  }
};
const fileReaderPng = global.FileReader = class { /**
 * Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
 * This method is used for testing purposes to mimic the behavior of FileReader.
 */
  readAsDataURL() {
    this.onload({target: {result: 'data:image/png'}});
  }
};
const fileReaderJpg = global.FileReader = class { /**
 * Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
 * This method is used for testing purposes to mimic the behavior of FileReader.
 */
  readAsDataURL() {
    this.onload({target: {result: 'data:image/jpg'}});
  }
};
const fileReaderJpeg = global.FileReader = class { /**
 * Simulates the behavior of the FileReader's readAsDataURL method by triggering the onload event with a mock result.
 * This method is used for testing purposes to mimic the behavior of FileReader.
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

// Mock sessionStorage
const sessionStorageMock = (() => {
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
global.sessionStorage = sessionStorageMock;

// Mock logError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
}));

describe('openPersonalizePage function', () => {
  beforeEach(async () => {
    // Reset DOM and mocks before each test
    document.getElementById('page-contents').innerHTML = ``;
    sessionStorageMock.clear();
    localStorageMock.clear();
    jest.clearAllMocks();

    const personalize = await import('../src/js/personalize.js');
    personalize.openPersonalizePage();
  });

  it('should populate page contents when opened', () => {
    // Arrange
    const pageContents = document.getElementById('page-contents').innerHTML;

    // Act

    // Assert
    expect(pageContents).toContain('personalize-container');
    expect(pageContents).toContain('class="personalize-container"');
    expect(pageContents).toContain('Change icon');
    expect(pageContents).toContain('Change logo');
    expect(pageContents).toContain('Change title');
    expect(sessionStorage.getItem('savedPage')).toBe('7');
  });

  it('should display custom modal when change title button is clicked', () => {
    // Arrange
    // Act
    fireEvent.click(document.querySelector('.title-button'));

    // Assert
    expect(document.getElementById('custom-modal').style.display).toBe('block');
  });
  it('should save the new title and close modal on save', () => {
    // Arrange
    const input = document.getElementById('new-title-input');

    // Act
    fireEvent.click(document.querySelector('.title-button'));
    fireEvent.input(input, {target: {value: 'New Title'}});
    fireEvent.click(document.getElementById('saveTitleButton'));

    // Assert
    expect(document.getElementById('custom-modal').style.display).toBe('none');
  });
  it('should not close modal or save when input is empty', () => {
    // Arrange
    const input = document.getElementById('new-title-input');

    // Act
    fireEvent.click(document.querySelector('.title-button'));
    fireEvent.input(input, {target: {value: ''}});
    fireEvent.click(document.getElementById('saveTitleButton'));

    // Assert
    expect(document.getElementById('custom-modal').style.display).toBe('block');
  });

  it('should attach a click event listener to the icon button', () => {
    // Arrange
    const iconButton = document.querySelector('.icon-button');
    const mockClick = jest.fn();

    // Act
    expect(iconButton).not.toBeNull();
    iconButton.addEventListener('click', mockClick);
    iconButton.click();

    // Assert
    expect(mockClick).toHaveBeenCalled();
  });

  it('should attach a click event listener to the logo button', () => {
    // Arrange
    const logoButton = document.querySelector('.logo-button');
    const mockClick = jest.fn();

    // Act
    logoButton.addEventListener('click', mockClick);
    logoButton.click();

    // Assert
    expect(logoButton).not.toBeNull();
    expect(mockClick).toHaveBeenCalled();
  });

  it('should attach a click event listener to the title button', () => {
    // Arrange
    const titleButton = document.querySelector('.title-button');
    const mockClick = jest.fn();

    // Act
    titleButton.addEventListener('click', mockClick);
    titleButton.click();

    // Arrange
    expect(titleButton).not.toBeNull();
    expect(mockClick).toHaveBeenCalled();
  });
  it('should call resetSettings function on reset button click', async () => {
    // Arrange
    localStorage.setItem('favicon', 'custom.ico');
    localStorage.setItem('title', 'tempTitle');
    const resetButton = document.querySelector('.reset-button');

    document.getElementById('page-contents').innerHTML += `
      <img id="logo" src="">
      <input type="text" id="newTitle" value="New Page Title">
      <h1 id="title">Old Title</h1>
      <link rel="icon" id="favicon" href="">
    `;
    const logo = document.getElementById('logo');
    const title = document.getElementById('title');
    const favicon = document.getElementById('favicon');
    const setItemMock = jest.spyOn(localStorage, 'setItem').mockImplementation(() => {});
    
    // Act
    resetButton.click();

    setItemMock.mockRestore();    
    // Arrange
    expect(localStorage.getItem('favicon')) === null;
    expect(localStorage.getItem('title')) === null;
    expect(logo).not.toBeNull();
    expect(logo.src).toContain('frontend/src/assets/images/logoTeamA-transformed.png');
    expect(title).not.toBeNull();
    expect(title.textContent).toBe('Little Brother');
    expect(favicon).not.toBeNull();
    expect(favicon.href).toContain('frontend/src/assets/images/logoTeamA-transformed.png');
  });

  it('should save selected theme to localStorage on theme change', () => {
    // Arrange
    const darkThemeRadio = document.getElementById('dark');

    // Act
    darkThemeRadio.click();

    // Assert

    expect(localStorage.getItem('theme')).toBe('dark');
    expect(document.getElementById('dark').checked).toBe(true);
    expect(document.getElementById('normal').checked).toBe(false);
  });
  it('should check the correct theme radio button based on localStorage value', async () => {
    // Arrange

    localStorage.setItem('theme', 'dark');
    const personalize = await import('../src/js/personalize.js');

    // act
    personalize.openPersonalizePage();

    expect(document.getElementById('dark').checked).toBe(true);
    expect(document.documentElement.className).toBe('dark');
    expect(document.getElementById('normal').checked).toBe(false);
  });
});

describe('handleFaviconSelect', () => {
  it('should do nothing if no file is selected', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');
    const iconEvent = {target: {files: []}}; // No file is selected

    // Spy on FileReader to check if it's called
    const spyFileReader = jest.spyOn(window, 'FileReader');

    // Act
    personalize.handleFaviconChange(iconEvent);

    // Assert
    expect(spyFileReader).not.toHaveBeenCalled(); // FileReader should not be instantiated
    expect(document.querySelector('link[rel="icon"]')).toBeNull(); // No favicon link element should be added
  });
  it('should change the favicon when a valid .ico file is selected', async () => {
    // Arrange
    const head = global.document.querySelector('head');
    const personalize = await import('../src/js/personalize.js');

    // Act
    FileReader = fileReaderIco;
    const blob = new Blob(['dummy'], {type: 'image/x-icon'});
    personalize.handleFaviconChange({target: {files: [blob]}});

    // Assert
    const newFavicon = head.querySelector('link[rel="icon"]');
    test.value(newFavicon.href).isEqualTo('data:image/x-icon');
  });

  it('should change the favicon when a valid .png file is selected', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');
    const head = document.querySelector('head');

    // Act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    personalize.handleFaviconChange({target: {files: [blob]}});

    // Assert
    const newFavicon = head.querySelector('link[rel="icon"]');
    test.value(newFavicon.href).isEqualTo('data:image/png');
  });
  it('saves valid .ico favicon in localStorage', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');
    // Act
    FileReader = fileReaderIco;
    const blob = new Blob(['dummy'], {type: 'image/x-icon'});
    personalize.handleFaviconChange({target: {files: [blob]}});

    // Assert
    const favicon = localStorageMock.getItem('favicon');
    const expectedValue = 'data:image/x-icon';
    test.value(favicon).isEqualTo(expectedValue);
  });
  it('saves valid .png favicon in localStorage', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');
    // Act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    personalize.handleFaviconChange({target: {files: [blob]}});

    // Assert
    test.value(localStorageMock.getItem('favicon')).isEqualTo('data:image/png');
  });
});

describe('handlePictureChange', () => {
  it('changes navigation picture with valid .png file', async () => {
    // Arrange
    const logo = document.getElementById('logo');
    const personalize = await import('../src/js/personalize.js');

    // Act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    personalize.handlePictureChange({target: {files: blob}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/png');
  });
  it('changes navigation picture with valid .jpg file', async () => {
    // Arrange
    const logo = document.getElementById('logo');
    const personalize = await import('../src/js/personalize.js');

    // Act
    FileReader = fileReaderJpg;
    const blob = new Blob(['dummy'], {type: 'image/jpg'});
    personalize.handlePictureChange({target: {files: blob}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpg');
  });
  it('changes navigation picture with valid .jpeg file', async () => {
    // Arrange
    const logo = document.getElementById('logo');
    const personalize = await import('../src/js/personalize.js');

    // Act
    FileReader = fileReaderJpeg;
    const blob = new Blob(['dummy'], {type: 'image/jpeg'});
    personalize.handlePictureChange({target: {files: blob}});

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpeg');
  });
  it('saves valid .png file in localStorage', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');

    // act
    FileReader = fileReaderPng;
    const blob = new Blob(['dummy'], {type: 'image/png'});
    personalize.handlePictureChange({target: {files: blob}});

    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/png');
  });
  it('saves valid .jpg file in localStorage', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');

    // act
    FileReader = fileReaderJpg;

    const blob = new Blob(['dummy'], {type: 'image/jpg'});
    personalize.handlePictureChange({target: {files: blob}});
    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/jpg');
  });
  it('saves .jpeg file in localStorage', async () => {
    // Arrange
    const personalize = await import('../src/js/personalize.js');

    // act
    FileReader = fileReaderJpeg;
    const blob = new Blob(['dummy'], {type: 'image/jpeg'});
    personalize.handlePictureChange({target: {files: blob}});

    // Assert
    const localStoragePicture = localStorageMock.getItem('picture');
    const expectedValue = 'data:image/jpeg';
    test.value(localStoragePicture).isEqualTo(expectedValue);
  });
});

describe('handleTitleChange', () => {
  it('should change the title of the page', async () => {
    // Arrange
    const newTitleInput = document.getElementById('newTitle');
    const titleElement = document.getElementById('title');
    const personalize = await import('../src/js/personalize.js');

    // Act
    personalize.handleTitleChange(newTitleInput.value);

    // Assert
    test.value(titleElement.textContent).isEqualTo(newTitleInput.value);
  });

  it('should save the new title to localStorage', async () => {
    // Arrange
    const newTitleInput = document.getElementById('newTitle');
    const personalize = await import('../src/js/personalize.js');

    // Act
    personalize.handleTitleChange(newTitleInput.value);

    // Assert
    const localStorageTitle = localStorageMock.getItem('title');
    const expectedValue = newTitleInput.value;
    test.value(localStorageTitle).isEqualTo(expectedValue);
  });
});

describe('retrieveTheme', () => {
  it('should apply the stored theme class to the document root', async () => {
    // Arrange
    const expectedThemeClass = 'dark-theme';
    localStorage.setItem('theme', expectedThemeClass);
    const personalize = await import('../src/js/personalize.js');

    // Act
    personalize.retrieveTheme();

    // Assert
    const appliedThemeClass = document.documentElement.className;
    test.value(appliedThemeClass).isEqualTo(expectedThemeClass);
  });
});
