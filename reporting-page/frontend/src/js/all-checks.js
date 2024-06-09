import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {getLocalization, getLocalizationString} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {openIssuePage, scrollToElement} from './issue.js';

/**
 * Load the content of the About page
 * @param {String} area area to scroll to on the page when opened
 */
export async function openAllChecksPage(area) {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('all-checks-button');
  sessionStorage.setItem('savedPage', '8');

  document.getElementById('page-contents').innerHTML = `
  <div class="all-checks">
    <div class="all-checks-container">
      <div class="all-checks-segment all-checks-title" id="top"> <!-- title top segment -->
        <p class="lang-security-risk-areas"><p> 
      </div>
      <div class="all-checks-segment" id="applications">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-applications">Applications</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-applications-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Applications">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="securityApplications"></div>
      </div>
      <div class="all-checks-segment" id="devices"> 
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-devices">Devices</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-devices-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Devices">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="securityDevices"></div>
      </div>
      <div class="all-checks-segment" id="network">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-network">Network</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-network-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Network">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="securityNetwork"></div>
      </div>
      <div class="all-checks-segment"id="os">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-os" >Operating System</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-os-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.OS">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="securityOS"></div>
      </div>
      <div class="all-checks-segment" id="passwords">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-passwords">Passwords</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-passwords-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Passwords">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="securityPasswords" data-area="passwords"></div>
      </div>
      <div class="all-checks-segment" id="security-other">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-other">Other</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-other-security-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Security">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="securityOther" data-area="security-other"></div>
      </div>
      <div class="all-checks-segment all-checks-title"> <!-- title bottom segment -->
        <p class="lang-privacy-risk-areas"><p> 
      </div>
      <div class="all-checks-segment" id="permissions">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-permissions">Permissions</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-permissions-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Permissions">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="privacyPermissions" data-area="permissions"></div>
      </div>
      <div class="all-checks-segment" id="browser">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-browser">Browser</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-browser-text"> 
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Browser">
          Here are some checks we run regarding your browser:
        </p>
        <div class="all-checks-segment-header">
          <p>Google Chrome</p>
        </div>
        <div class="checksList" id="privacyBrowserChrome" data-area="browser"></div>
        <div class="all-checks-segment-header">
          <p>Microsoft Edge</p>
        </div>
        <div class="checksList" id="privacyBrowserEdge" data-area="browser"></div>
        <div class="all-checks-segment-header">
          <p>Mozilla Firefox</p>
        </div>
        <div class="checksList" id="privacyBrowserFirefox" data-area="browser"></div>
      </div>
      <div class="all-checks-segment" id="privacy-other">
        <div class="all-checks-segment-header">
          <p class="lang-all-checks-other">Other</p>
        </div>
        <p class="all-checks-segment-text lang-all-checks-other-privacy-text">
        </p>
        <p class="all-checks-segment-line" id="AllChecks.Privacy">
          Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="privacyOther" data-area="privacy-other"></div>
      </div>
    </div>
  </div>
  `;

  const elements = document.getElementsByClassName('checksList');
  for (let i = 0; i < elements.length; i++) {
    elements[i].innerHTML = createBulletList(elements[i].id);
    // Makes sure the grammar of the line before the bullet list is correct
    const numb = elements[i].firstChild.childElementCount;
    if (numb > 1) {
      elements[i].parentElement.childNodes[5].classList.add('lang-all-check-text-line');
    } else {
      elements[i].parentElement.childNodes[5].classList.add('lang-all-check-text-line-single');
    }
  }
  const issues = JSON.parse(sessionStorage.getItem('DataBaseData'));
  const checks = document.getElementsByClassName('all-checks-check');
  for (let i = 0; i < checks.length; i++) {
    const issue = issues.find((issue) => issue.id == checks[i].id);
    checks[i].addEventListener('click',
      () => openIssuePage(issue.jsonkey, issue.severity, checks[i].parentElement.parentElement.parentElement.id));
  }

  // Localize the static content of the about page
  const staticAboutPageConent = [
    'lang-security-risk-areas',
    'lang-privacy-risk-areas',
    'lang-all-check-text-line',
    'lang-all-check-text-line-single',
    'lang-all-checks-applications',
    'lang-all-checks-devices',
    'lang-all-checks-network',
    'lang-all-checks-os',
    'lang-all-checks-passwords',
    'lang-all-checks-other',
    'lang-all-checks-permissions',
    'lang-all-checks-browser',
    'lang-all-checks-applications-text',
    'lang-all-checks-devices-text',
    'lang-all-checks-network-text',
    'lang-all-checks-os-text',
    'lang-all-checks-passwords-text',
    'lang-all-checks-other-security-text',
    'lang-all-checks-permissions-text',
    'lang-all-checks-browser-text',
    'lang-all-checks-other-privacy-text',
  ];
  const localizationIds = [
    'Dashboard.SecurityRiskAreas',
    'Dashboard.PrivacyRiskAreas',
    'AllChecks.MultipleLines',
    'AllChecks.SingleLine',
    'AllChecks.Applications',
    'AllChecks.Devices',
    'AllChecks.Network',
    'AllChecks.OS',
    'AllChecks.Passwords',
    'AllChecks.Other',
    'AllChecks.Permissions',
    'AllChecks.Browser',
    'AllChecks.ApplicationsText',
    'AllChecks.DevicesText',
    'AllChecks.NetworkText',
    'AllChecks.OSText',
    'AllChecks.PasswordsText',
    'AllChecks.OtherSecurityText',
    'AllChecks.PermissionsText',
    'AllChecks.BrowserText',
    'AllChecks.OtherPrivacyText',
  ];
  for (let ids = 1; ids < 43; ids++) {
    staticAboutPageConent.push('lang-id'+ids.toString());
    localizationIds.push('AllChecks.Id'+ids.toString());
  }

  for (let i = 0; i < staticAboutPageConent.length; i++) {
    getLocalization(localizationIds[i], staticAboutPageConent[i]);
  }

  // Add the right suffix to the lines before the bullet lists
  const lines = document.getElementsByClassName('all-checks-segment-line');
  for (let i = 0; i < lines.length; i++) {
    const text = await getLocalizationString(lines[i].id);
    lines[i].innerHTML += text.toLowerCase() + ':';
  }

  const element = getViewedElement(area);
  scrollToElement(element);
}

/**
 * Get the node of the page you want to view when opening the page.
 * @param {String} area Area on the page to view
 * @return {HTMLElement} Node where to view the page one
 */
export function getViewedElement(area) {
  let element;
  switch (area) {
  case 'applications':
    element = document.getElementById('applications');
    break;
  case 'devices':
    element = document.getElementById('devices');
    break;
  case 'network':
    element = document.getElementById('network');
    break;
  case 'os':
    element = document.getElementById('os');
    break;
  case 'passwords':
    element = document.getElementById('passwords');
    break;
  case 'security-other':
    element = document.getElementById('security-other');
    break;
  case 'permissions':
    element = document.getElementById('permissions');
    break;
  case 'browser':
    element = document.getElementById('browser');
    break;
  case 'privacy-other':
    element = document.getElementById('privacy-other');
    break;
  default:
    element = document.getElementById('top');
  }
  return element;
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('all-checks-button').addEventListener('click', () => openAllChecksPage(''));
  } catch (error) {
    logError('Error in all-checks.js: ' + error);
  }
}

export const areaLists = {
  'securityApplications': [20],
  'securityDevices': [1, 2],
  'securityNetwork': [11, 13, 40, 41],
  'securityOS': [14, 15, 17, 18, 19, 33, 37],
  'securityPasswords': [5, 16, 38],
  'securityOther': [3, 12, 32, 34, 39, 42],
  'privacyPermissions': [6, 7, 8, 9, 10],
  'privacyBrowserChrome': [21, 23, 25, 35],
  'privacyBrowserEdge': [22, 24, 26, 36],
  'privacyBrowserFirefox': [29, 31, 30, 27, 28],
  'privacyOther': [4],
};

/**
   * create a bullet list for each entry of a security or privacy area
   * @param {string} listId id of the list to create a bullet list for
   * @return {string} returns an html list
   */
export function createBulletList(listId) {
  const list = areaLists[listId];
  let resultLine = `<ul>`;
  list.forEach((check) => {
    resultLine += `<li class="lang-id${check} all-checks-check" id="${check}"></li>`;
  });
  resultLine += `</ul>`;
  return resultLine;
}
