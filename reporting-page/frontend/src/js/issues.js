import data from '../database.json' assert { type: 'json' };
import {openIssuePage} from './issue.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

/** Load the content of the Issues page */
export function openIssuesPage() {
  closeNavigation();
  markSelectedNavigationItem('issues-button');

  const pageContents = document.getElementById('page-contents');
  pageContents.innerHTML = `
  <div class="issues-data">
    <div class="table-container">
      <h2>Issue table</h2>
      <table class="issues-table" id="issues-table">
        <thead>
          <tr>
          <th class="issue-column">
            <span class="table-header name">Name</span>
            <span class="material-symbols-outlined" id="sort-on-issue">swap_vert</span>
          </th>
          <th class="type-column">
            <span class="table-header type">Type</span>
            <span class="material-symbols-outlined" id="sort-on-type">swap_vert</span>
          </th>
          <th class="risk-column">
            <span class="table-header risk">Risk level</span>
            <span class="material-symbols-outlined" id="sort-on-risk">swap_vert</span>
          </th>
          </tr>
        </thead>
        <tbody>
        </tbody>
      </table>
    </div>
    <div class="dropdown-container">
      <button id="dropbtn-table" class="dropbtn-table"><span class="select-risks">Select Risks</span></button>
      <div class="dropdown-selector-table" id="myDropdown-table">
        <p><input type="checkbox" checked="true" value="true" id="select-high-risk-table">
          <label for="select-high-risk" class="high-risk-issues"> High risks</label><br>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-medium-risk-table">
          <label for="select-medium-risk" class="medium-risk-issues"> Medium risks</label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-low-risk-table">
          <label for="select-low-risk" class="low-risk-issues"> Low risks</label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-info-risk-table">
          <label for="select-info-risk" class="info-risk-issues"> Informative</label>
        </p>
      </div>
    </div>
    <div class="table-container">
      <h2>Non issue table</h2>
      <table class="issues-table" id="non-issues-table">
        <thead>
          <tr>
          <th class="issue-column">
            <span class="table-header name">Name</span>
            <span class="material-symbols-outlined" id="sort-on-issue2">swap_vert</span>
          </th>
          <th class="type-column">
            <span class="table-header type">Type</span>
            <span class="material-symbols-outlined" id="sort-on-type2">swap_vert</span>
          </th>
          </tr>
        </thead>
        <tbody>
        </tbody>
      </table>
    </div>


  </div>
  `;

  const tableHeaders = ['name', 'type', 'risk'];
  const localizationIds = ['Issues.Name', 'Issues.Type', 'Issues.Risk'];
  for (let i = 0; i < tableHeaders.length; i++) {
    getLocalization(localizationIds[i], tableHeaders[i]);
  }

  // retrieve issues from tray application
  const issues = JSON.parse(sessionStorage.getItem('DataBaseData'));

  const issueTable = document.getElementById('issues-table').querySelector('tbody');
  fillTable(issueTable, issues, true);

  const nonIssueTable = document.getElementById('non-issues-table').querySelector('tbody');
  fillTable(nonIssueTable, issues, false);

  const myDropdownTable = document.getElementById('myDropdown-table');
  document.getElementById('dropbtn-table').addEventListener('click', () => myDropdownTable.classList.toggle('show'));
  document.getElementById('select-high-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-medium-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-low-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-info-risk-table').addEventListener('change', changeTable);

  document.onload = retrieveTheme();
}
/**
 * Returns the risk level based on the given numeric level.
 * @param {number} level - The numeric representation of the risk level.
 * @return {string} The risk level corresponding to the numeric input:
 */
function toRiskLevel(level) {
  switch (level) {
  case 0:
    return 'Acceptable';
  case 1:
    return 'Low';
  case 2:
    return 'Medium';
  case 3:
    return 'High';
  case 4:
    return 'Info';
  }
}

/** Fill the table with issues
 *
 * @param {HTMLTableSectionElement} tbody Table to be filled
 * @param {Issue} issues Issues to be filled in
 * @param {Bool} isIssue True for issue table, false for non issue table
 */
export function fillTable(tbody, issues, isIssue) {
  issues.forEach((issue) => {
    const currentIssue = data[issue.jsonkey];

    if (isIssue) {
      if (currentIssue) {
        if (issue.severity != '0') {
          const row = document.createElement('tr');
          row.innerHTML = `
              <td class="issue-link">${currentIssue.Name}</td>
              <td>${currentIssue.Type}</td>
              <td>${toRiskLevel(issue.severity)}</td>
            `;
          row.cells[0].id = issue.jsonkey;
          tbody.appendChild(row);
        }
      }
    } else {
      if (currentIssue) {
        if (issue.severity == '0') {
          const row = document.createElement('tr');
          row.innerHTML = `
              <td class="issue-link">${currentIssue.Name}</td>
              <td>${currentIssue.Type}</td>
            `;
          row.cells[0].id = issue.jsonkey;
          tbody.appendChild(row);
        }
      }
    }
  });

  // Add links to issue information pages
  const issueLinks = document.querySelectorAll('.issue-link');
  issueLinks.forEach((link) => {
    link.addEventListener('click', () => openIssuePage(link.id));
  });

  // Add buttons to sort on columns
  if (isIssue) {
    document.getElementById('sort-on-issue').addEventListener('click', () => sortTable(tbody, 0));
    document.getElementById('sort-on-type').addEventListener('click', () => sortTable(tbody, 1));
    document.getElementById('sort-on-risk').addEventListener('click', () => sortTable(tbody, 2));
  } else {
    document.getElementById('sort-on-issue2').addEventListener('click', () => sortTable(tbody, 0));
    document.getElementById('sort-on-type2').addEventListener('click', () => sortTable(tbody, 1));
  }
}

/** Sorts the table
 *
 * @param {HTMLTableSectionElement} tbody Table to be sorted
 * @param {number} column Column to sort the table on
 */
export function sortTable(tbody, column) {
  const table = tbody.closest('table');
  let direction = table.getAttribute('data-sort-direction');
  direction = direction === 'ascending' ? 'descending' : 'ascending';
  const rows = Array.from(tbody.rows);
  rows.sort((a, b) => {
    if (column !== 2) {
      // Alphabetical sorting for other columns
      const textA = a.cells[column].textContent.toLowerCase();
      const textB = b.cells[column].textContent.toLowerCase();
      if (direction === 'ascending') {
        return textA.localeCompare(textB);
      } else {
        return textB.localeCompare(textA);
      }
    } else {
      // Custom sorting for the last column
      const order = {'high': 1, 'medium': 2, 'low': 3, 'acceptable': 4, 'info': 5};
      const textA = a.cells[column].textContent.toLowerCase();
      const textB = b.cells[column].textContent.toLowerCase();
      if (direction === 'ascending') {
        return order[textA] - order[textB];
      } else {
        return order[textB] - order[textA];
      }
    }
  });
  while (tbody.rows.length > 0) {
    tbody.deleteRow(0);
  }
  rows.forEach((row) => {
    tbody.appendChild(row);
  });
  table.setAttribute('data-sort-direction', direction);
}

if (typeof document !== 'undefined') {
  try {
    document.getElementById('issues-button').addEventListener('click', () => openIssuesPage());
  } catch (error) {
    logError('Error in issues.js: ' + error);
  }
}

/**
 * Updates the displayed issues table based on the selected risk levels.
 * Retrieves issues data from session storage, filters it based on selected risk levels,
 * and updates the table with the filtered data.
 */
function changeTable() {
  const selectedHigh = document.getElementById('select-high-risk-table').checked;
  const selectedMedium = document.getElementById('select-medium-risk-table').checked;
  const selectedLow = document.getElementById('select-low-risk-table').checked;
  const selectedInfo = document.getElementById('select-info-risk-table').checked;

  const issues = JSON.parse(sessionStorage.getItem('DataBaseData'));

  const issueTable = document.getElementById('issues-table').querySelector('tbody');

  // Filter issues based on the selected risk levels
  const filteredIssues = issues.filter((issue) => {
    return (
      (selectedHigh && issue.severity === 3) ||
      (selectedMedium && issue.severity === 2) ||
      (selectedLow && issue.severity === 1) ||
      (selectedInfo && issue.severity === 4)
    );
  });

  // Clear existing table rows
  issueTable.innerHTML = '';

  // Refill tables with filtered issues
  fillTable(issueTable, filteredIssues, true);
}
