import {jest} from '@jest/globals';

export function mockPageFunctions() {
  // Mock LogError
  jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
    LogError: jest.fn(),
  }));

  // Mock Navigation
  jest.unstable_mockModule('../src/js/navigation-menu.js', () => ({
    closeNavigation: jest.fn(),
    markSelectedNavigationItem: jest.fn(),
    loadPersonalizeNavigation: jest.fn(),
  }));

  // Mock retrieveTheme
  jest.unstable_mockModule('../src/js/personalize.js', () => ({
    retrieveTheme: jest.fn(),
  }));  
}

/** Mock of getLocalization function with no functionality
 *
 * @param {string} messageID - The ID of the message to be localized.
 * @return {string} The messageID passed as the argument.
 */
export function mockGetLocalization(messageID) {
  const myPromise = new Promise(function(myResolve, myReject) {
     if (messageID !== '') myResolve(messageID) 
     else myReject(new Error('error'));
    });
  return myPromise;
}

// Create mock mouse events
export const clickEvent = new window.MouseEvent('click');
export const beginHover = new window.MouseEvent('mouseenter');
export const endHover = new window.MouseEvent('mouseleave');

// Mock global storage
export const storageMock = (() => {
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
