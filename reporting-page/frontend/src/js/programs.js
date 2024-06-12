import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {sortTable} from './issues.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the Home page */
export function openProgramsPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('programs-button');
  sessionStorage.setItem('savedPage', 5);

  document.getElementById('page-contents').innerHTML = `
    <div class="program-data">
        <div class="program-container">
        <h2 class="lang-program-table"></h2>
        <input type="text" id="search-input" placeholder="Search software...">
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

  const programsJson = JSON.parse(sessionStorage.getItem('ScanResult'));
  const issueTableHtml = document.getElementById('program-table').querySelector('tbody');
  const foundObject = programsJson.find((obj) => obj.issue_id === 37);

  // Check if the object was found
  if (foundObject) {
    fillProgamTable(issueTableHtml, foundObject.result);
  } else {
    console.log(`Object with ID ${targetId} not found.`);
  }

  document.getElementById('sort-on-issue').addEventListener('click', () => sortTable(issueTableHtml, 0));

  document.getElementById('search-input').addEventListener('input', function(event) {
    const query = event.target.value;
    searchTable(issueTableHtml, query);
  });

  const tableHeaders = [
    'lang-program-table',
    'lang-name',
  ];
  const localizationIds = [
    'Programs.ProgramTable',
    'Programs.Name',
    'Programs.Search',
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

/**
 * Fills the program table with the given programs
 * @param {HTMLTableSectionElement} tbody Table body to be filled
 * @param {Array} programs Programs to be added to the table
 * @return {void}
 * */
export function fillProgamTable(tbody, programs) {
  programs.forEach((program) => {
    const row = document.createElement('tr');
    const name= program.split(' | ')[0];
    const version = program.split(' | ')[1];

    row.innerHTML = `
            <td class="issue-link">${name}</td>
            <td>${version}</td>
        `;
    tbody.appendChild(row);
  });
}

/**
 * Filters the table based on the search query
 *
 * @param {HTMLTableSectionElement} tbody Table body to be filtered
 * @param {string} query Search query
 */
export function searchTable(tbody, query) {
  const rows = tbody.querySelectorAll('tr'); // Alternative to tbody.rows
  const lowerCaseQuery = query.toLowerCase();

  rows.forEach((row) => {
    const nameCell = row.cells[0]; // Assuming the name is in the first column
    const nameText = nameCell.textContent.toLowerCase();

    if (nameText.includes(lowerCaseQuery)) {
      row.style.display = ''; // Show row
    } else {
      row.style.display = 'none'; // Hide row
    }
  });
}

