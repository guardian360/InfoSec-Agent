import {getLocalization, getLocalizationString} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the Home page */
export async function openProgramsPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('programs-button');
  sessionStorage.setItem('savedPage', 5);
  const searchTranslation = await getLocalizationString('Programs.Search');

  // Find the result of the programs check
  const programsFromStorage = sessionStorage.getItem('ProgramList');
  if (!programsFromStorage) {
    document.getElementById('page-contents').innerHTML = `
    <div class="reload-container">
      <div class="reload-segment">
        <p class="lang-reload-text"><p>
        <div class="button lang-reload-button" id="reload-button"></div>
      </div>
    </div>
    `;
    document.getElementById('reload-button').addEventListener('click', () => openProgramsPage());
  } else {
    const programsJson = JSON.parse(programsFromStorage);
    logError('Programs JSON: ' + JSON.stringify(programsJson));

    document.getElementById('page-contents').innerHTML = `
    <div class="program-data">
        <div class="program-container">
        <h2 class="lang-program-table"></h2>
        <input type="text" id="search-input" placeholder="${searchTranslation}">
        <table class="program-table" id="program-table">
            <thead>
            <tr>
            <th class="program-column">
                <span class="table-header lang-name"></span>
            </th>
            <th class="version-column">
                <span class="table-header lang-version"></span>
            </th>
            </tr>
            </thead>
            <tbody>
            </tbody>
        </table>
        </div>
    </div>
    `;

    const issueTableHtml = document.getElementById('program-table').querySelector('tbody');
    fillProgamTable(issueTableHtml, programsJson.result);

    // Add event listeners for sorting and searching
    document.getElementById('search-input').addEventListener('input', function(event) {
      const query = event.target.value;
      searchTable(issueTableHtml, query);
    });
  }

  // Translate the page contents
  const tableHeaders = [
    'lang-program-table',
    'lang-name',
    'lang-version',
    'lang-reload-text',
    'lang-reload-button',
  ];
  const localizationIds = [
    'Programs.ProgramTable',
    'Programs.Name',
    'Programs.Version',
    'Programs.ReloadText',
    'Programs.ReloadButton',
  ];
  for (let i = 0; i < tableHeaders.length; i++) {
    getLocalization(localizationIds[i], tableHeaders[i]);
  }
}

/**
 * Fills the program table with the given programs
 *
 * @param {HTMLTableSectionElement} tbody Table body to be filled
 * @param {Array} programs Programs to be added to the table
 * @return {void}
 * */
export function fillProgamTable(tbody, programs) {
  programs.forEach((program) => {
    // Parse program data
    const name = program.split(' | ')[0];
    const version = program.split(' | ')[1];

    // Create a new row for the program
    const row = document.createElement('tr');
    row.innerHTML = `
            <td class="issue-link">${name}</td>
            <td>${version}</td>
        `;
    tbody.appendChild(row);
  });

  // Sort the table after filling
  sortProgramTable(tbody, 'ascending');
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

  // Sort the table after filtering
  sortProgramTable(tbody, 'ascending');
}

/** Sorts the rows of the issues table
 *
 * @param {HTMLTableSectionElement} tbody Table to be sorted
 * @param {string} direction Direction to sort in
 */
export function sortProgramTable(tbody, direction) {
  // Sort the table rows
  const rows = Array.from(tbody.rows);
  rows.sort((a, b) => {
    const nameA = a.cells[0].textContent.toLowerCase();
    const nameB = b.cells[0].textContent.toLowerCase();
    // Sort on program name
    if (direction === 'ascending') {
      return nameA.localeCompare(nameB);
    } else {
      return nameB.localeCompare(nameA);
    }
  });

  // Clear the table and refill it with the sorted rows
  while (tbody.rows.length > 0) {
    tbody.deleteRow(0);
  }
  rows.forEach((row) => {
    tbody.appendChild(row);
  });
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('programs-button').addEventListener('click', () => openProgramsPage());
  } catch (error) {
    logError('Error in programs.js: ' + error);
  }
}
