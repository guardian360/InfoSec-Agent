import {closeNavigation} from './navigation-menu';
import {markSelectedNavigationItem} from './navigation-menu';
import {retrieveTheme} from './personalize';

/** Load the content of the About page */
function openAboutPage() {
  closeNavigation();
  markSelectedNavigationItem('about-button');

  document.getElementById('page-contents').innerHTML = `
  <div class="dashboard-data"></div>
  `;
  document.onload = retrieveTheme();
}

document.getElementById('about-button').addEventListener('click', () => openAboutPage());
