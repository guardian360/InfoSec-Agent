// import {GetLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';

/** Load the content of the Privacy Dashboard page */
function openPrivacyDashboardPage() {
  closeNavigation();
  markSelectedNavigationItem('privacy-dashboard-button');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard-data"></div>
  `;
  document.onload = retrieveTheme();
}

document.getElementById('privacy-dashboard-button').addEventListener('click', () => openPrivacyDashboardPage());
