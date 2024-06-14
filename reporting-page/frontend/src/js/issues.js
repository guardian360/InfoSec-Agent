import dataDe from '../databases/database.de.json' assert { type: 'json' };
import dataEnGB from '../databases/database.en-GB.json' assert { type: 'json' };
import dataEnUS from '../databases/database.en-US.json' assert { type: 'json' };
import dataEs from '../databases/database.es.json' assert { type: 'json' };
import dataFr from '../databases/database.fr.json' assert { type: 'json' };
import dataNl from '../databases/database.nl.json' assert { type: 'json' };
import dataPt from '../databases/database.pt.json' assert { type: 'json' };

import {openIssuePage} from './issue.js';
import {getLocalization} from './localize.js';
import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {LoadUserSettings as loadUserSettings} from '../../wailsjs/go/main/App.js';

/** Load the content of the Issues page */
export function openIssuesPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('issues-button');
  sessionStorage.setItem('savedPage', '4');
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
    <div class="dropdown-container">
      <button id="dropbtn-table" class="dropbtn-table"><span class="lang-select-risks"></span></button>
      <div class="dropdown-selector-table" id="myDropdown-table">
        <p><input type="checkbox" checked="true" value="true" id="select-high-risk-table">
          <label for="select-high-risk" class="lang-high-risk-issues"></label><br>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-medium-risk-table">
          <label for="select-medium-risk" class="lang-medium-risk-issues"></label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-low-risk-table">
          <label for="select-low-risk" class="lang-low-risk-issues"></label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-info-risk-table">
          <label for="select-info-risk" class="lang-info-risk-issues"></label>
        </p>
      </div>
    </div>
    <div class="table-container">
      <h2 class="lang-acceptable-findings"></h2>
      <table class="issues-table" id="non-issues-table">
        <thead>
          <tr>
          <th class="issue-column">
            <span class="table-header lang-name"></span>
            <span class="material-symbols-outlined" id="sort-on-issue2">swap_vert</span>
          </th>
          <th class="type-column">
            <span class="table-header lang-type"></span>
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


  // retrieve issues from tray application
  const issues = JSON.parse(sessionStorage.getItem('ScanResult'));

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

  const tableHeaders = [
    'lang-issue-table',
    'lang-acceptable-findings',
    'lang-name',
    'lang-type',
    'lang-risk',
    'lang-high-risk-issues',
    'lang-medium-risk-issues',
    'lang-low-risk-issues',
    'lang-info-risk-issues',
    'lang-select-risks',
    'lang-acceptable',
    'lang-low',
    'lang-medium',
    'lang-high',
    'lang-info',
  ];
  const localizationIds = [
    'Issues.IssueTable',
    'Issues.AcceptableFindings',
    'Issues.Name',
    'Issues.Type',
    'Issues.Risk',
    'Dashboard.HighRisk',
    'Dashboard.MediumRisk',
    'Dashboard.LowRisk',
    'Dashboard.InfoRisk',
    'Dashboard.SelectRisks',
    'Issues.Acceptable',
    'Issues.Low',
    'Issues.Medium',
    'Issues.High',
    'Issues.Info',
  ];
  for (let i = 0; i < tableHeaders.length; i++) {
    getLocalization(localizationIds[i], tableHeaders[i]);
  }
}

/**
 * Returns the risk level based on the given numeric level.
 * @param {number} level - The numeric representation of the risk level.
 * @return {string} The risk level corresponding to the numeric input:
 */
export function toRiskLevel(level) {
  switch (level) {
  case 0:
    return '<td class="lang-acceptable"></td>';
  case 1:
    return '<td class="lang-low"></td>';
  case 2:
    return '<td class="lang-medium"></td>';
  case 3:
    return '<td class="lang-high"></td>';
  case 4:
    return '<td class="lang-info"></td>';
  }
}

/** Fill the table with issues
 *
 * @param {HTMLTableSectionElement} tbody Table to be filled
 * @param {Issue} issues Issues to be filled in
 * @param {Bool} isIssue True for issue table, false for non issue table
 * @param {Bool} isListenersAdded True for the first time the eventlisteners is called
 */
export async function fillTable(tbody, issues, isIssue, isListenersAdded=true) {
  const language = await getUserSettings();
  let currentIssue;

  issues.forEach((issue) => {
    switch (language) {
    case 0:
      currentIssue = dataDe[issue.issue_id];
      break;
    case 1:
      currentIssue = dataEnGB[issue.issue_id];
      break;
    case 2:
      currentIssue = dataEnUS[issue.issue_id];
      break;
    case 3:
      currentIssue = dataEs[issue.issue_id];
      break;
    case 4:
      currentIssue = dataFr[issue.issue_id];
      break;
    case 5:
      currentIssue = dataNl[issue.issue_id];
      break;
    case 6:
      currentIssue = dataPt[issue.issue_id];
      break;
    default:
      currentIssue = dataEnGB[issue.issue_id];
    }

    if (isIssue) {
      if (currentIssue[issue.result_id]) {
        const issueData = currentIssue[issue.result_id];
        if (issueData.Severity != '0') {
          const riskLevel = toRiskLevel(issueData.Severity);
          const row = document.createElement('tr');

          row.innerHTML = `
              <td class="issue-link" data-result-id="${issue.result_id}">${issueData.Name}</td>
              <td>${currentIssue.Type}</td>
              ${riskLevel}
            `;

          row.cells[0].id = issue.issue_id;
          row.setAttribute('data-result-id', issue.result_id);
          tbody.appendChild(row);
        }
      }
    } else {
      if (currentIssue[issue.result_id]) {
        const issueData = currentIssue[issue.result_id];
        if (issueData.Severity == '0') {
          const row = document.createElement('tr');
          row.innerHTML = `
              <td class="issue-link" data-result-id="${issue.result_id}">${issueData.Name}</td>
              <td>${currentIssue.Type}</td>
            `;
          row.cells[0].id = issue.issue_id;
          tbody.appendChild(row);
        }
      }
    }
  });

  // Add links to issue information pages
  const issueLinks = document.querySelectorAll('.issue-link');
  issueLinks.forEach((link) => {
    link.addEventListener('click', () => openIssuePage(link.id, link.getAttribute('data-result-id')));
  });

  // Add buttons to sort on columns
  if (isListenersAdded) {
    if (isIssue) {
      document.getElementById('sort-on-issue').addEventListener('click', () => sortTable(tbody, 0));
      document.getElementById('sort-on-type').addEventListener('click', () => sortTable(tbody, 1));
      document.getElementById('sort-on-risk').addEventListener('click', () => sortTable(tbody, 2));
    } else {
      document.getElementById('sort-on-issue2').addEventListener('click', () => sortTable(tbody, 0));
      document.getElementById('sort-on-type2').addEventListener('click', () => sortTable(tbody, 1));
    }
    isListenersAdded = false;
  }
  // Re-apply localization to the dynamically created table rows
  const tableHeaders = [
    'lang-low',
    'lang-medium',
    'lang-high',
    'lang-info',
  ];

  const localizationIds = [
    'Issues.Low',
    'Issues.Medium',
    'Issues.High',
    'Issues.Info',
  ];
  for (let i = 0; i < tableHeaders.length; i++) {
    getLocalization(localizationIds[i], tableHeaders[i]);
  }

  // Sort the table in descending order of risk severity
  const sorting = JSON.parse(sessionStorage.getItem('IssuesSorted'));
  if (sorting) {
    const table = tbody.closest('table');
    if (table) {
      const direction = sorting.direction === 'ascending' ? 'descending' : 'ascending';
      table.setAttribute('data-sort-direction', direction);
      sortTable(tbody, parseInt(sorting.column));
    }
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
      // Change Info to lower severity
      let severityA = parseInt(a.getAttribute('data-severity'));
      if (severityA === 4) {
        severityA = 0;
      }
      let severityB = parseInt(b.getAttribute('data-severity'));
      if (severityB === 4) {
        severityB = 0;
      }
      if (direction === 'ascending') {
        return severityB - severityA;
      } else {
        return severityA - severityB;
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
  const columnValue = column.toString();
  sessionStorage.setItem('IssuesSorted', JSON.stringify({'column': columnValue, 'direction': direction}));
}

/* istanbul ignore next */
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
export function changeTable() {
  const selectedHigh = document.getElementById('select-high-risk-table').checked;
  const selectedMedium = document.getElementById('select-medium-risk-table').checked;
  const selectedLow = document.getElementById('select-low-risk-table').checked;
  const selectedInfo = document.getElementById('select-info-risk-table').checked;

  const issues = JSON.parse(sessionStorage.getItem('ScanResult'));

  const issueTable = document.getElementById('issues-table').querySelector('tbody');

  // Filter issues based on the selected risk levels
  const filteredIssues = issues.filter((issue) => {
    const issueData = dataEnGB[issue.issue_id][issue.result_id];
    if (issueData) { // Until checks that return an error are added to database
      return (
        (selectedLow && issueData.Severity === 1) ||
        (selectedMedium && issueData.Severity === 2) ||
        (selectedHigh && issueData.Severity === 3) ||
        (selectedInfo && issueData.Severity === 4)
      );
    }
    return false
  });

  // Clear existing table rows
  issueTable.innerHTML = '';

  // Refill tables with filtered issues
  fillTable(issueTable, filteredIssues, true, false);
}

/**
 * Retrieves the user settings including the preferred language.
 *
 * This function asynchronously loads user settings and returns the user's
 * preferred language as an integer. The language is represented by the
 * following integers:
 * 0 - German
 * 1 - English (UK)
 * 2 - English (US)
 * 3 - Spanish
 * 4 - French
 * 5 - Dutch
 * 6 - Portuguese
 *
 * @function getUserSettings
 * @return {Promise<number>} A promise that resolves to the user's preferred language as an integer.
 */
export async function getUserSettings() {
  try {
    const userSettings = await loadUserSettings();
    const language = userSettings.Language;
    return language;
  } catch (error) {
    logError('Error loading user settings:', error);
  }
}
