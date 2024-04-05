import 'jsdom-global/register.js';
import test from 'unit.js';
import { JSDOM } from "jsdom";
import { MarkSelectedNavigationItem, CloseNavigation, ToggleNavigationHamburger, ToggleNavigationResize } from '../src/js/navigation-menu.js';

// Mock page
const dom = new JSDOM(`
`);
global.document = dom.window.document
global.window = dom.window

// Test cases
describe('Navigation menu', function() {
  it('markSelectedNavigationItem should mark the item in the navigation menu that is selected', function() {
    // Arrange

    // Act

    // Assert
  });

  it('closeNavigation should close the navigation menu if a menu item is selected and the screen size is smaller than 800 px', function() {
    // Arrange

    // Act

    // Assert
  });

  it('toggleNavigationHamburger should open or close the navigation menu when the hamburger menu is clicked', function() {
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
