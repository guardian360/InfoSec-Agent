import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Give the selected navigation item a different color
 * @param {string} item - The navigation item that is selected
*/
export function markSelectedNavigationItem(item) {
  localize();
  const navItems = document.getElementsByClassName('nav-link');
  const stylesheet = getComputedStyle(document.documentElement);
  for (let i = 1; i < navItems.length; i++) {
    navItems[i].style.backgroundColor = stylesheet.getPropertyValue('--background-color-left-nav');
  }

  if (item === 'issue-button' || item === 'personalize-button') {
    return;
  }

  document.getElementById(item).style.backgroundColor = stylesheet.getPropertyValue('--background-nav-hover');
}

/** Close the navigation menu when a navigation item is clicked, only when screen size is less than 800px
 * @param {int} appWidth - The width of the application screen
*/
export function closeNavigation(appWidth) {
  if (appWidth < 800) {
    document.getElementsByClassName('left-nav')[0].style.visibility = 'hidden';
  }
}

/** Open or close the navigation menu when user clicks on hamburger menu
 * @param {int} appWidth - The width of the application screen
*/
export function toggleNavigationHamburger(appWidth) {
  if (appWidth < 800) {
    if (document.getElementsByClassName('left-nav')[0].style.visibility === 'visible') {
      document.getElementsByClassName('left-nav')[0].style.visibility = 'hidden';
    } else {
      document.getElementsByClassName('left-nav')[0].style.visibility = 'visible';
    }
  }
}

/** Open or close the navigation menu when user resizes the screen
 * @param {int} appWidth - The width of the application screen
*/
export function toggleNavigationResize(appWidth) {
  if (appWidth > 799) {
    document.getElementsByClassName('left-nav')[0].style.visibility = 'visible';
  } else {
    document.getElementsByClassName('left-nav')[0].style.visibility = 'hidden';
  }
}
/** Localizes the navigation menu and sets up event listeners for responsive behavior. */
function localize() {
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
  const localizationIds = [
    'Navigation.Home',
    'Navigation.SecurityDashboard',
    'Navigation.PrivacyDashboard',
    'Navigation.Issues',
    'Navigation.Integration',
    'Navigation.About',
    'Navigation.Personalize',
    'Navigation.ChangeLanguage',
  ];
  for (let i = 0; i < navbarItems.length; i++) {
    getLocalization(localizationIds[i], navbarItems[i]);
  }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    const header = document.getElementById('header-hamburger');
    header.addEventListener('click', () => toggleNavigationHamburger(document.body.offsetWidth));
    document.body.onresize = () => toggleNavigationResize(document.body.offsetWidth);
  } catch (error) {
    logError('Error in navigation-menu.js: ' + error);
  }
}
