import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
  <div class="left-nav">
    <div id="home-button" class="nav-link">
      <p><span class="material-symbols-outlined">home</span><span class="nav-item home">Home</span></p>
    </div>
    <div id="security-dashboard-button" class="nav-link">
      <p><span class="material-symbols-outlined">security</span><span class="nav-item security-dashboard">Security Dashboard</span></p>
    </div>
    <div id="privacy-dashboard-button" class="nav-link">
      <p><span class="material-symbols-outlined">lock</span><span class="nav-item privacy-dashboard">Privacy Dashboard</span></p>
    </div>
    <div id="issues-button" class="nav-link">
      <p><span class="material-symbols-outlined">checklist</span><span class="nav-item issues">Issues</span></p>
    </div>
    <div id="integration-button" class="nav-link">
      <p><span class="material-symbols-outlined">integration_instructions</span><span class="nav-item integration">Integration</span></p>
    </div>
    <div id="about-button" class="nav-link">
      <p><span class="material-symbols-outlined">info</span><span class="nav-item about" >About</span></p>
    </div>
  </div>
`, {
  url: 'http://localhost',
});
global.document = dom.window.document;
global.window = dom.window;

// Mock LogError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
}));

// Test cases
describe('Navigation menu', function() {
  it('closeNavigation should close the navigation menu if screen size is small', async function() {
    const navigationMenu = await import('../src/js/navigation-menu.js');
    const leftNav = document.getElementsByClassName('left-nav')[0];

    navigationMenu.closeNavigation(799);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('hidden');

    leftNav.style.visibility = 'visible';
    navigationMenu.closeNavigation(800);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');
  });
  it('toggleNavigationHamburger should open or close the navigation menu', async function() {
    const navigationMenu = await import('../src/js/navigation-menu.js');
    const leftNav = document.getElementsByClassName('left-nav')[0];

    leftNav.style.visibility = 'visible';
    navigationMenu.toggleNavigationHamburger(799);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('hidden');

    leftNav.style.visibility = 'hidden';
    navigationMenu.toggleNavigationHamburger(799);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');

    leftNav.style.visibility = 'visible';
    navigationMenu.toggleNavigationHamburger(800);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');
  });
  it('toggleNavigationResize should open or close the navigation menu when the screen is resized', async function() {
    const navigationMenu = await import('../src/js/navigation-menu.js');
    const leftNav = document.getElementsByClassName('left-nav')[0];

    navigationMenu.toggleNavigationResize(799);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('hidden');

    navigationMenu.toggleNavigationResize(800);
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');
  });
});
