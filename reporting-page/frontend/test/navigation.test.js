import 'jsdom-global/register.js';
// import test from 'unit.js';
import {JSDOM} from 'jsdom';
// import {
//   markSelectedNavigationItem,
//   closeNavigation,
//   toggleNavigationHamburger,
//   toggleNavigationResize
// } from '../src/js/navigation-menu.js';

global.TESTING = true

// Mock page
const dom = new JSDOM(`
`);
global.document = dom.window.document;
global.window = dom.window;

// Test cases
describe('Navigation menu', function() {
  it('markSelectedNavigationItem should mark the item in the navigation menu that is selected', function() {
    // Arrange

    // Act

    // Assert
  });
  it('closeNavigation should close the navigation menu if screen size is small', function() {
    // Arrange

    // Act

    // Assert
  });
  it('toggleNavigationHamburger should open or close the navigation menu', function() {
    // Arrange

    // Act

    // Assert
  });
  it('toggleNavigationResize should open or close the navigation menu when the screen is resized', function() {
    // Arrange

    // Act

    // Assert
  });
});
