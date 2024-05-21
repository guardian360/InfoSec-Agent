import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockGetLocalization, clickEvent, resizeEvent} from './mock.js';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
  <input type="file" class="personalize-input-invisible" id="faviconInput" accept=".png,.ico"> 
  <div class="header">
    <div class="header-hamburger container">
      <span id="header-hamburger" class="header-hamburger material-symbols-outlined">menu</span>
    </div>
    <div class="header-logo">
      <div id="logo-button" class="logo-name">
        <img id="logo" alt="logo" src="./src/assets/images/logoTeamA-transformed.png"> 
        <div class="header-name">
          <h1 id="title">Little Brother</h1><!-- Use id to dynamically change title -->
        </div>
      </div>
    </div>
    <div class="header-settings">
      <div class="nav-link settings-button">
        <span><span class="material-symbols-outlined">settings</span></span>
        <div class="dropdown-content">
          <a id="personalize-button" class="lang-personalize-page">Personalize page</a>
          <a id="language-button" class="lang-change-language">Change Language</a>
        </div>
      </div>
    </div>
  </div> 
  <div class="left-nav">
    <div id="home-button" class="nav-link">
      <p><span class="material-symbols-outlined">home</span><span class="nav-item lang-home">Home</span></p>
    </div>
    <div id="security-dashboard-button" class="nav-link">
      <p><span class="material-symbols-outlined">security</span>
      <span class="nav-item lang-security-dashboard">Security Dashboard</span></p>
    </div>
    <div id="privacy-dashboard-button" class="nav-link">
      <p><span class="material-symbols-outlined">lock</span>
      <span class="nav-item lang-privacy-dashboard">Privacy Dashboard</span></p>
    </div>
    <div id="issues-button" class="nav-link">
      <p><span class="material-symbols-outlined">checklist</span><span class="nav-item lang-issues">Issues</span></p>
    </div>
    <div id="integration-button" class="nav-link">
      <p><span class="material-symbols-outlined">integration_instructions</span>
      <span class="nav-item lang-integration">Integration</span></p>
    </div>
    <div id="about-button" class="nav-link">
      <p><span class="material-symbols-outlined">info</span><span class="nav-item lang-about" >About</span></p>
    </div>
  </div>
  <div id="page-contents"></div>
  <div class="page-contents"></div>
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

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
}));

// Test cases
describe('Navigation menu', function() {
  it('markNavigationItem should give the current navigation item the correct background color', async function() {
    // Arrange
    const navigationMenu = await import('../src/js/navigation-menu.js');
    document.documentElement.style.setProperty('--background-color-left-nav', 'red');
    document.documentElement.style.setProperty('--background-nav-hover', 'blue');

    // Act
    navigationMenu.markSelectedNavigationItem('integration-button');

    // Assert
    test.value(document.getElementById('issues-button').style.backgroundColor).isEqualTo('red');
    test.value(document.getElementById('integration-button').style.backgroundColor).isEqualTo('blue');

    // Act
    navigationMenu.markSelectedNavigationItem('settings-button');

    // Assert
    test.value(document.getElementById('integration-button').style.backgroundColor).isEqualTo('red');
    test.value(document.getElementsByClassName('settings-button')[0].style.backgroundColor).isNotEqualTo('blue');
  });
  it('closeNavigation should close the navigation menu if screen size is small', async function() {
    // Arrange
    const navigationMenu = await import('../src/js/navigation-menu.js');
    const leftNav = document.getElementsByClassName('left-nav')[0];

    // Act
    navigationMenu.closeNavigation(799);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('hidden');

    // Act
    leftNav.style.visibility = 'visible';
    navigationMenu.closeNavigation(800);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');
  });
  it('toggleNavigationHamburger should open or close the navigation menu', async function() {
    // Arrange
    const navigationMenu = await import('../src/js/navigation-menu.js');
    const leftNav = document.getElementsByClassName('left-nav')[0];

    // Act
    leftNav.style.visibility = 'visible';
    // click button which calls toggleNavigationHamburger
    document.getElementById('header-hamburger').dispatchEvent(clickEvent);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('hidden');

    // Act
    leftNav.style.visibility = 'hidden';
    // click button which calls toggleNavigationHamburger
    document.getElementById('header-hamburger').dispatchEvent(clickEvent);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');

    // Act
    leftNav.style.visibility = 'visible';
    navigationMenu.toggleNavigationHamburger(800);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');
  });
  it('toggleNavigationResize should open or close the navigation menu when the screen is resized', async function() {
    // Arrange
    const navigationMenu = await import('../src/js/navigation-menu.js');
    const leftNav = document.getElementsByClassName('left-nav')[0];

    // Act
    // resize window which calls toggleNavigationResize with a appwidth of 0
    window.dispatchEvent(resizeEvent);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('hidden');

    // Act
    navigationMenu.toggleNavigationResize(800);

    // Assert
    test.value(window.getComputedStyle(leftNav).getPropertyValue('visibility')).isEqualTo('visible');
  });
  it('items are localized', async function() {
    // Arrange
    await import('../src/js/navigation-menu.js');
    const navbarItems = [
      'lang-home',
      'lang-security-dashboard',
      'lang-privacy-dashboard',
      'lang-issues',
      'lang-integration',
      'lang-about',
      'lang-personalize-page',
      'lang-change-language',
    ];
    const expectedNames = [
      'Navigation.Home',
      'Navigation.SecurityDashboard',
      'Navigation.PrivacyDashboard',
      'Navigation.Issues',
      'Navigation.Integration',
      'Navigation.About',
    ];

    // Assert
    expectedNames.forEach((name, index) => {
      test.value(document.getElementsByClassName(navbarItems[index])[0].innerHTML).isEqualTo(name);
    });
  });
});
