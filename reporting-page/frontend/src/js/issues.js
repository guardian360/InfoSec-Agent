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
      <div class="table-header-container">
        <h2 class="lang-issue-table"></h2>
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
          <p><input type="checkbox" checked="true" value="true" id="select-acceptable-risk-table">
            <label for="select-acceptable-risk" class="lang-acceptable-risk-issues"></label>
          </p>
          <p><input type="checkbox" checked="true" value="true" id="select-info-risk-table">
            <label for="select-info-risk" class="lang-info-risk-issues"></label>
          </p>
        </div>
      </div>
    </div>
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

  // Fill the table with issues
  var issues = getIssues();
  if (issues) {
    sessionStorage.setItem('IssuesList', JSON.stringify(issues));
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    fillTable(issueTable, issues);

    const sortingMethod = JSON.parse(sessionStorage.getItem('IssuesSorting'));
    if (sortingMethod) {
      refillTable(issueTable, sortingMethod);
    } else {
      console.log("IssuesSorting not found, setting default sorting method");
      const defaultSorting = {'column': 2, 'direction': 'descending'};
      sessionStorage.setItem('IssuesSorting', JSON.stringify(defaultSorting));
      refillTable(issueTable, defaultSorting);
    }
  } else {
    logError('Error in issues.js: Issues not found');
  }

  // Add event listeners for the table filter menu
  const myDropdownTable = document.getElementById('myDropdown-table');
  document.getElementById('dropbtn-table').addEventListener('click', () => myDropdownTable.classList.toggle('show'));
  document.getElementById('select-high-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-medium-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-low-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-acceptable-risk-table').addEventListener('change', changeTable);
  document.getElementById('select-info-risk-table').addEventListener('change', changeTable);

  // Add buttons to sort on columns
  document.getElementById('sort-on-issue').addEventListener('click', () => sortTable(0));
  document.getElementById('sort-on-type').addEventListener('click', () => sortTable(1));
  document.getElementById('sort-on-risk').addEventListener('click', () => sortTable(2));

  // Add localization to the page
  const tableHeaders = [
    'lang-issue-table',
    'lang-acceptable-findings',
    'lang-name',
    'lang-type',
    'lang-risk',
    'lang-high-risk-issues',
    'lang-medium-risk-issues',
    'lang-low-risk-issues',
    'lang-acceptable-risk-issues',
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
    'Dashboard.Acceptable',
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

/** Get the issues from the database
 * 
 * @returns {Issue[]} List of issues
 * */
export function getIssues() {
  // Get checks results from session storage
  const issues = JSON.parse(sessionStorage.getItem('DataBaseData'));
  
  // Get issue information in the user's preferred language
  var issueList = [];
  const language = getUserSettings();
  let currentIssue;
  console.log("Issues in getIssues:", issues);
  issues.forEach((issue) => {
    switch (language) {
    case 0:
      currentIssue = dataDe[issue.jsonkey];
      break;
    case 1:
      currentIssue = dataEnGB[issue.jsonkey];
      break;
    case 2:
      currentIssue = dataEnUS[issue.jsonkey];
      break;
    case 3:
      currentIssue = dataEs[issue.jsonkey];
      break;
    case 4:
      currentIssue = dataFr[issue.jsonkey];
      break;
    case 5:
      currentIssue = dataNl[issue.jsonkey];
      break;
    case 6:
      currentIssue = dataPt[issue.jsonkey];
      break;
    default:
      currentIssue = dataEnGB[issue.jsonkey];
    }

    // Add issue to list
    if (currentIssue) {
      const name = currentIssue.Name;
      const type = currentIssue.Type;
      const severity = issue.severity;
      const jsonkey = issue.jsonkey;

      issueList.push({"name": name, "type": type, "severity": severity, "jsonkey": jsonkey});
    }
  });
  return issueList;
}

/** Fill the table with issues
 *
 * @param {HTMLTableSectionElement} tbody Table to be filled
 * @param {Issue} issues Issues to be filled in
 * @param {Bool} isIssue True for issue table, false for non issue table
 * @param {Bool} isListenersAdded True for the first time the eventlisteners is called
 */
export function fillTable(tbody, issues) {
  // Add a table row for each issue
  issues.forEach((issue) => {
    const riskLevel = toRiskLevel(issue.severity);
    const row = document.createElement('tr');

    row.innerHTML = `
      <td class="issue-link" data-severity="${issue.severity}">${issue.name}</td>
      <td>${issue.type}</td>
      ${riskLevel}
    `;

    row.cells[0].id = issue.jsonkey;
    row.setAttribute('data-severity', issue.severity);
    tbody.appendChild(row);
  });

  // Add links to issue information pages
  const issueLinks = document.querySelectorAll('.issue-link');
  issueLinks.forEach((link) => {
    link.parentElement.addEventListener('click', () => openIssuePage(link.id, link.getAttribute('data-severity')));
  });

  // Re-apply localization to the dynamically created table rows
  const tableHeaders = [
    'lang-acceptable',
    'lang-low',
    'lang-medium',
    'lang-high',
    'lang-info',
  ];
  const localizationIds = [
    'Issues.Acceptable',
    'Issues.Low',
    'Issues.Medium',
    'Issues.High',
    'Issues.Info',
  ];
  for (let i = 0; i < tableHeaders.length; i++) {
    getLocalization(localizationIds[i], tableHeaders[i]);
  }

  // Sort the table
  const sortingMethod = JSON.parse(sessionStorage.getItem('IssuesSorting'));
  refillTable(tbody, sortingMethod);
}

/** Updates the sorting method and sorts the table
 *
 * @param {number} column Column to sort the table on
 */
export function sortTable(column) {
  // Update sorting method
  var sortingMethod = JSON.parse(sessionStorage.getItem('IssuesSorting'));
  var direction = sortingMethod.direction;
  direction = direction === 'ascending' ? 'descending' : 'ascending';
  sortingMethod = {'column': column, 'direction': direction};
  sessionStorage.setItem('IssuesSorting', JSON.stringify(sortingMethod));

  // Refill the table with the new sorting method
  const issueTable = document.getElementById('issues-table').querySelector('tbody');
  refillTable(issueTable, sortingMethod);
}

/** Sorts the rows of the issues table
 *
 * @param {HTMLTableSectionElement} tbody Table to be sorted
 * @param {Object} sortingMethod Sorting method
 */
export function refillTable(tbody, sortingMethod) {
  // Get sorting method
  const column = sortingMethod.column;
  const direction = sortingMethod.direction;

  // Sort the table rows
  const table = tbody.closest('table');
  const rows = Array.from(tbody.rows);
  rows.sort((a, b) => {
    const nameA = a.cells[0].textContent.toLowerCase();
    const nameB = b.cells[0].textContent.toLowerCase();
    const typeA = a.cells[1].textContent.toLowerCase();
    const typeB = b.cells[1].textContent.toLowerCase();
    var severityA = parseInt(a.getAttribute('data-severity'));
    if (severityA === 0) {severityA = -1;}
    if (severityA === 4) {severityA = 0;}
    console.log("Severity A: ", severityA);
    var severityB = parseInt(b.getAttribute('data-severity'));
    if (severityB === 0) {severityB = -1;}
    if (severityB === 4) {severityB = 0;}
    console.log("Severity B: ", severityB);
    // Sort on issue name
    if (column === 0) {
      if (direction === 'ascending') {
        return nameA.localeCompare(nameB) || severityB - severityA || typeA.localeCompare(typeB);
      } else {
        return nameB.localeCompare(nameA) || severityB - severityA || typeA.localeCompare(typeB);
      }
    } 
    // Sort on issue type
    else if (column === 1) {
      if (direction === 'ascending') {
        return typeA.localeCompare(typeB) || severityB - severityA || nameA.localeCompare(nameB);
      } else {
        return typeB.localeCompare(typeA) || severityB - severityA || nameA.localeCompare(nameB);
      }
    } 
    // Sort on risk level
    else if (column === 2) {
      if (direction === 'ascending') {
        return severityB - severityA || typeB.localeCompare(typeA) || nameA.localeCompare(nameB);
      } else {
        return severityA - severityB || typeB.localeCompare(typeA) || nameA.localeCompare(nameB);
      }
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

/**
 * Returns the risk level based on the given numeric level.
 * 
 * @param {number} level - The numeric representation of the risk level.
 * @return {string} The risk level corresponding to the numeric input:
 */
export function toRiskLevel(level) {
  switch (level) {
  case 0:
    return '<td><span class="table-risk-level lang-acceptable"></span></td>';
  case 1:
    return '<td><span class="table-risk-level lang-low"></span></td>';
  case 2:
    return '<td><span class="table-risk-level lang-medium"></span></td>';
  case 3:
    return '<td><span class="table-risk-level lang-high"></span></td>';
  case 4:
    return '<td><span class="table-risk-level lang-info"></span></td>';
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
  const selectedAcceptable = document.getElementById('select-acceptable-risk-table').checked;
  const selectedInfo = document.getElementById('select-info-risk-table').checked;

  var issues = JSON.parse(sessionStorage.getItem('IssuesList'));

  const issueTable = document.getElementById('issues-table').querySelector('tbody');

  // Filter issues based on the selected risk levels
  const filteredIssues = issues.filter((issue) => {
    return (
      (selectedAcceptable && issue.severity === 0) ||
      (selectedLow && issue.severity === 1) ||
      (selectedMedium && issue.severity === 2) ||
      (selectedHigh && issue.severity === 3) ||
      (selectedInfo && issue.severity === 4)
    );
  });

  // Clear existing table rows
  issueTable.innerHTML = '';

  // Refill tables with filtered issues
  fillTable(issueTable, filteredIssues);
}

// /** Fill the table with issues
//  *
//  * @param {HTMLTableSectionElement} tbody Table to be filled
//  * @param {Issue} issues Issues to be filled in
//  * @param {Bool} isIssue True for issue table, false for non issue table
//  * @param {Bool} isListenersAdded True for the first time the eventlisteners is called
//  */
// export async function fillTable(tbody, issues, isListenersAdded=true) {
//   // Add a table row for each issue
//   const language = await getUserSettings();
//   let currentIssue;
//   issues.forEach((issue) => {
//     switch (language) {
//     case 0:
//       currentIssue = dataDe[issue.jsonkey];
//       break;
//     case 1:
//       currentIssue = dataEnGB[issue.jsonkey];
//       break;
//     case 2:
//       currentIssue = dataEnUS[issue.jsonkey];
//       break;
//     case 3:
//       currentIssue = dataEs[issue.jsonkey];
//       break;
//     case 4:
//       currentIssue = dataFr[issue.jsonkey];
//       break;
//     case 5:
//       currentIssue = dataNl[issue.jsonkey];
//       break;
//     case 6:
//       currentIssue = dataPt[issue.jsonkey];
//       break;
//     default:
//       currentIssue = dataEnGB[issue.jsonkey];
//     }

//     if (currentIssue) {
//       const riskLevel = toRiskLevel(issue.severity);
//       const row = document.createElement('tr');

//       row.innerHTML = `
//         <td class="issue-link" data-severity="${issue.severity}">${currentIssue.Name}</td>
//         <td>${currentIssue.Type}</td>
//         ${riskLevel}
//       `;

//       row.cells[0].id = issue.jsonkey;
//       row.setAttribute('data-severity', issue.severity);
//       tbody.appendChild(row);
//     }
//   });

//   // Add links to issue information pages
//   const issueLinks = document.querySelectorAll('.issue-link');
//   issueLinks.forEach((link) => {
//     link.parentElement.addEventListener('click', () => openIssuePage(link.id, link.getAttribute('data-severity')));
//   });

//   // Add buttons to sort on columns
//   if (isListenersAdded) {
//     document.getElementById('sort-on-issue').addEventListener('click', () => sortTable(tbody, 0));
//     document.getElementById('sort-on-type').addEventListener('click', () => sortTable(tbody, 1));
//     document.getElementById('sort-on-risk').addEventListener('click', () => sortTable(tbody, 2));
//     isListenersAdded = false;
//   }

//   // Re-apply localization to the dynamically created table rows
//   const tableHeaders = [
//     'lang-acceptable',
//     'lang-low',
//     'lang-medium',
//     'lang-high',
//     'lang-info',
//   ];
//   const localizationIds = [
//     'Issues.Acceptable',
//     'Issues.Low',
//     'Issues.Medium',
//     'Issues.High',
//     'Issues.Info',
//   ];
//   for (let i = 0; i < tableHeaders.length; i++) {
//     getLocalization(localizationIds[i], tableHeaders[i]);
//   }

//   // Sort the table in descending order of risk severity
//   const sorting = JSON.parse(sessionStorage.getItem('IssuesSorting'));
//   if (sorting) {
//     const table = tbody.closest('table');
//     if (table) {
//       const direction = sorting.direction === 'ascending' ? 'descending' : 'ascending';
//       table.setAttribute('data-sort-direction', direction);
//       sortTable(tbody, parseInt(sorting.column));
//     }
//   }
// }

// /** Sorts the table
//  *
//  * @param {HTMLTableSectionElement} tbody Table to be sorted
//  * @param {number} column Column to sort the table on
//  */
// export function sortTable(tbody, column) {
//   const table = tbody.closest('table');
//   let direction = table.getAttribute('data-sort-direction');
//   direction = direction === 'ascending' ? 'descending' : 'ascending';
//   const rows = Array.from(tbody.rows);
//   rows.sort((a, b) => {
//     if (column !== 2) {
//       // Alphabetical sorting for other columns
//       const textA = a.cells[column].textContent.toLowerCase();
//       const textB = b.cells[column].textContent.toLowerCase();
//       if (direction === 'ascending') {
//         return textA.localeCompare(textB);
//       } else {
//         return textB.localeCompare(textA);
//       }
//     } else {
//       // Change Info to lower severity
//       let severityA = parseInt(a.getAttribute('data-severity'));
//       if (severityA === 0) {
//         severityA = -1;
//       }
//       if (severityA === 4) {
//         severityA = 0;
//       }
//       let severityB = parseInt(b.getAttribute('data-severity'));
//       if (severityB === 0) {
//         severityB = -1;
//       }
//       if (severityB === 4) {
//         severityB = 0;
//       }
//       if (direction === 'ascending') {
//         return severityB - severityA;
//       } else {
//         return severityA - severityB;
//       }
//     }
//   });
//   while (tbody.rows.length > 0) {
//     tbody.deleteRow(0);
//   }
//   rows.forEach((row) => {
//     tbody.appendChild(row);
//   });
//   table.setAttribute('data-sort-direction', direction);
//   const columnValue = column.toString();
//   sessionStorage.setItem('IssuesSorting', JSON.stringify({'column': columnValue, 'direction': direction}));
// }

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

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('issues-button').addEventListener('click', () => openIssuesPage());
  } catch (error) {
    logError('Error in issues.js: ' + error);
  }
}
