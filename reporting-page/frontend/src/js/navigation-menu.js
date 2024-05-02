import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Give the selected navigation item a different color
 * @param {string} item - The navigation item that is selected
*/
export function markSelectedNavigationItem(item) {
  const navItems = document.getElementsByClassName('nav-link');
  const stylesheet = getComputedStyle(document.documentElement);
  for (let i = 1; i < navItems.length; i++) {
    navItems[i].style.backgroundColor = stylesheet.getPropertyValue('--background-color-left-nav');
  }
  if (item === 'settings-button' ) {
    return;
  }
  document.getElementById(item).style.backgroundColor = stylesheet.getPropertyValue('--background-nav-hover');
}

/**
 * Loads personalized navigation by applying background color to navigation links.
 * Background color is retrieved from CSS variables.
 */
export function loadPersonalizeNavigation() {
  const navItems = document.getElementsByClassName('nav-link');
  const stylesheet = getComputedStyle(document.documentElement);
  for (let i = 1; i < navItems.length; i++) {
    navItems[i].style.backgroundColor = stylesheet.getPropertyValue('--background-color-left-nav');
  }
}

/** Close the navigation menu when a navigation item is clicked, only when screen size is less than 800px */
export function closeNavigation() {
  if (document.body.offsetWidth < 800) {
    document.getElementsByClassName('left-nav')[0].style.visibility = 'hidden';
  }
}

/** Open or close the navigation menu when user clicks on hamburger menu */
export function toggleNavigationHamburger() {
  if (document.body.offsetWidth < 800) {
    if (document.getElementsByClassName('left-nav')[0].style.visibility === 'visible') {
      document.getElementsByClassName('left-nav')[0].style.visibility = 'hidden';
    } else {
      document.getElementsByClassName('left-nav')[0].style.visibility = 'visible';
    }
  }
}

/** Open or close the navigation menu when user resizes the screen */
export function toggleNavigationResize() {
  if (document.body.offsetWidth > 799) {
    document.getElementsByClassName('left-nav')[0].style.visibility = 'visible';
  } else {
    document.getElementsByClassName('left-nav')[0].style.visibility = 'hidden';
  }
}

if (typeof document !== 'undefined') {
  try {
    document.getElementById('header-hamburger').addEventListener('click', () => toggleNavigationHamburger());
    document.body.onresize = () => toggleNavigationResize();

    const navbarItems = [
      'home',
      'security-dashboard',
      'privacy-dashboard',
      'issues',
      'integration',
      'about',
    ];
    const localizationIds = [
      'Navigation.Home',
      'Navigation.SecurityDashboard',
      'Navigation.PrivacyDashboard',
      'Navigation.Issues',
      'Navigation.Integration',
      'Navigation.About',
    ];
    for (let i = 0; i < navbarItems.length; i++) {
      getLocalization(localizationIds[i], navbarItems[i]);
    }
  } catch (error) {
    logError('Error in navigation-menu.js: ' + error);
  }
}
