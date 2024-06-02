import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the Home page */
export function openProgramsPage() {
    retrieveTheme();
    closeNavigation(document.body.offsetWidth);
    markSelectedNavigationItem('programs-button');
    sessionStorage.setItem('savedPage', 1);

    document.getElementById('page-contents').innerHTML = `
    <div class="program-data">
        <div class="program-container">
        <h2 class="lang-program-table"></h2>
        <table class="program-table" id="program-table">
            <thead>
            <tr>
            <th class="program-column">
                <span class="table-header lang-name"></span>
                <span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span>
            </th>
            <th class="version-column">
                <span class="table-header"></span>
            </th>
            </tr>
            </thead>
            <tbody>
            </tbody>
        </table>
        </div>
    </div>
    `;

    const tableHeaders = [
        'lang-program-table',
        'lang-name',
      ];
      const localizationIds = [
        'Programs.ProgramTable',
        'Programs.Name',
      ];
      for (let i = 0; i < tableHeaders.length; i++) {
        getLocalization(localizationIds[i], tableHeaders[i]);
      }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
    try {
    document.getElementById('programs-button').addEventListener('click', () => openProgramsPage());
    } catch (error) {
    logError('Error in programs.js: ' + error);
    }
}

