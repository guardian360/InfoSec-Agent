import { handlePictureChange, handleTitleChange } from '../src/js/personalize.js'; // Assuming the function is in picture.js
import { JSDOM } from 'jsdom';
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
global.FileReader = class {
  readAsDataURL() {
    this.onload({ target: { result: 'data:image/jpeg;base64' } });
  }
};

// Mock localStorage
const localStorageMock = (() => {
  let store = {};

  return {
    getItem: key => store[key],
    setItem: (key, value) => {
      store[key] = value.toString();
    },
    clear: () => {
      store = {};
    }
  };
})();
global.localStorage = localStorageMock;


// Test cases
describe('handlePictureChange', () => {
  it('should change the logo image', () => {
    // Arrange
    const fileInput = document.getElementById('picture-input');
    const logo = document.getElementById('logo');

    // Act
    handlePictureChange({ target: { files: [new Blob(['dummy'], { type: 'image/jpeg' })] } });

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpeg;base64');
  
    })
  it('should save in localstorage', () => {
    //Arrange
    const fileInput = document.getElementById('picture-input');
    const logo = document.getElementById('logo');
    // Act
    handlePictureChange({ target: { files: [new Blob(['dummy'], { type: 'image/jpeg' })] } });

    // Assert
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/jpeg;base64');
  })

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