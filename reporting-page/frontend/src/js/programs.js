import {PieChart} from './piechart.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {scanTest} from './database.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {openIssuePage} from './issue.js';
import {saveProgress, shareProgress, selectSocialMedia} from './share.js';
import data from '../databases/database.en-GB.json' assert { type: 'json' };
import {showModal} from './settings.js';

/** Load the content of the Home page */
export function openProgramsPage() {
    retrieveTheme();
    closeNavigation(document.body.offsetWidth);
    markSelectedNavigationItem('programs-button');
    sessionStorage.setItem('savedPage', 1);

    document.getElementById('page-contents').innerHTML = `
    <div class="issues-data">
        <div class="table-container">
        <h2 class="lang-issue-table"></h2>
        <table class="issues-table" id="issues-table">
            <thead>
            <tr>
            <th class="issue-column">
                <span class="table-header lang-name"></span>
                <span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span>
            </th>
            <th class="type-column">
                <span class="table-header lang-type"></span>
                <span class="material-symbols-outlined" id="sort-on-type">swap_vert</span>
            </th>
            <th class="risk-column">
                <span class="table-header lang-risk"></span>
                <span class="material-symbols-outlined" id="sort-on-risk">swap_vert</span>
            </th>
            </tr>
            </thead>
            <tbody>
            </tbody>
        </table>
        </div>
    </div>
    `;

    for (let i = 0; i < staticHomePageContent.length; i++) {
        getLocalization(localizationIds[i], staticHomePageContent[i]);
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

