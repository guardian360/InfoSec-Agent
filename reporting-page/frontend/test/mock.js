import {jest} from '@jest/globals';

/** Mocks common function used inside openPage functions */
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

export function mockChangeLanguage(bool) {
  const myPromise = new Promise(function(myResolve, myReject) {
    if (bool) myResolve() 
    else myReject(new Error('error'));
   });
 return myPromise;
}

/** Mock of Chart constructor and update function from chart.js */
export function mockChart() {
  // Mock Chart constructor
  jest.unstable_mockModule('chart.js/auto', () => ({
    Chart: jest.fn().mockImplementation((context, config) => {
      return {
      // properties
        type: config?.type || '',
        data: config?.data || {},
        options: config?.options || {},
        // functions
        update: jest.fn(),
      };
    })
  }));
}

/** Mock of graph class */
export function mockGraph() {
  // Mock Chart constructor
  jest.unstable_mockModule('../src/js/graph.js', () => ({
    Graph: jest.fn().mockImplementation((context, riskcount) => {
      return {
        // properties
        graphShowHighRisks : true,
        graphShowMediumRisks : true,
        graphShowLowRisks : true,
        graphShowNoRisks : true,
        graphShowInfoRisks : true,
        
        graphShowAmount : 1,
        
        barChart : {},
        rc : riskcount,
        // functions
        createGraphChart: jest.fn(),
        changeGraph: jest.fn(),
        toggleRisks: jest.fn(),
        graphDropdown: jest.fn(),
        getData: jest.fn(),
        getOptions: jest.fn(),
      };
    })
  }));
}

// Create mock mouse events
export const clickEvent = new window.MouseEvent('click');
export const beginHover = new window.MouseEvent('mouseenter');
export const endHover = new window.MouseEvent('mouseleave');
export const changeEvent = new Event('change');
export const resizeEvent = new Event('resize');

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
    removeItem: (key) => {
      delete store[key];
    }
  };
})();

/** Mocks of openPage functions from:
 *  home, security-dashboard, privacy-dashboard,
 *  issues, integration, about and personalize.
 */
export function mockOpenPageFunctions() {
  // Mock openHomePage
  jest.unstable_mockModule('../src/js/Home.js', () => ({
    openHomePage: jest.fn(),
  }));

  // Mock openSecurityDashboardPage
  jest.unstable_mockModule('../src/js/security-dashboard.js', () => ({
    openSecurityDashboardPage: jest.fn(),
  }));

  // Mock openPrivacyDashboardPage
  jest.unstable_mockModule('../src/js/privacy-dashboard.js', () => ({
    openPrivacyDashboardPage: jest.fn(),
  }));

  // Mock openIssuesPage
  jest.unstable_mockModule('../src/js/issues.js', () => ({
    openIssuesPage: jest.fn(),
  }));

  // Mock openIntegrationPage
  jest.unstable_mockModule('../src/js/integration.js', () => ({
    openIntegrationPage: jest.fn(),
  }));

  // Mock openAboutPage
  jest.unstable_mockModule('../src/js/about.js', () => ({
    openAboutPage: jest.fn(),
  }));

  // Mock openPersonalizePage
  jest.unstable_mockModule('../src/js/personalize.js', () => ({
    openPersonalizePage: jest.fn(),
    retrieveTheme: jest.fn(),
  }));
}
