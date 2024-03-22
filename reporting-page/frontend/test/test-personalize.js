import { changePicture } from '../src/js/personalize.js'; // Assuming the function is in picture.js
import { JSDOM } from 'jsdom';
import test from 'unit.js';

// Mock page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
  <input type="file" id="picture-input">
  <img id="logo" src="">
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
describe('changePicture', () => {
  it('should change the logo image and save in localStorage', () => {
    // Arrange
    const fileInput = document.getElementById('picture-input');
    const logo = document.getElementById('logo');

    // Act
    changePicture({ target: { files: [new Blob(['dummy'], { type: 'image/jpeg' })] } });

    // Assert
    test.value(logo.src).isEqualTo('data:image/jpeg;base64');
    test.value(localStorageMock.getItem('picture')).isEqualTo('data:image/jpeg;base64');
  
    })

});

