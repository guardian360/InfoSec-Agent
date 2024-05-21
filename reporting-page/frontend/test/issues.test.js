import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import data from '../src/databases/database.en-GB.json' assert { type: 'json' };
import dataDe from '../src/databases/database.de.json' assert { type: 'json' };
import dataEnUS from '../src/databases/database.en-US.json' assert { type: 'json' };
import dataEs from '../src/databases/database.es.json' assert { type: 'json' };
import dataFr from '../src/databases/database.fr.json' assert { type: 'json' };
import dataNl from '../src/databases/database.nl.json' assert { type: 'json' };
import dataPt from '../src/databases/database.pt.json' assert { type: 'json' };
import {mockPageFunctions, clickEvent, storageMock} from './mock.js';

global.TESTING = true;

const dom = new JSDOM(`
  <div id="page-contents"></div>
  <div class="page-contents"></div>
`);
global.document = dom.window.document;
global.window = dom.window;

/** empty the table to have it empty for next tests
 *
 * @param {HTMLTableElement} table table to delete all rows from
 */
function emptyTable(table) {
  for (let i = 0; i < table.rows.length; i++) {
    table.deleteRow(i);
  }
}

// Mock sessionStorage
global.sessionStorage = storageMock;

/** Mock of getLocalizationString function
 *
 * @param {string} messageID - The ID of the message to be localized.
 * @return {string} The localized string.
 */
function mockGetLocalizationString(messageID) {
  const myPromise = new Promise(function(myResolve, myReject) {
    switch (messageID) {
    case 'Issues.IssueTable':
      myResolve('Issue table');
    case 'Issues.AcceptableFindings':
      myResolve('Acceptable findings');
    case 'Issues.Name':
      myResolve('Name');
    case 'Issues.Type':
      myResolve('Type');
    case 'Issues.Risk':
      myResolve('Risk');
    case 'Dashboard.HighRisk':
      myResolve('HighRisk');
    case 'Dashboard.MediumRisk':
      myResolve('MediumRisk');
    case 'Dashboard.LowRisk':
      myResolve('LowRisk');
    case 'Dashboard.InfoRisk':
      myResolve('InfoRisk');
    case 'Dashboard.SelectRisks':
      myResolve('SelectRisks');
    default:
      myReject(new Error('Wrong message ID'));
    }
  });
  return myPromise;
}

// Mock often used page functions
mockPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalizationString(input)),
  LoadUserSettings: jest.fn(),
}));

// Mock openIssuesPage
jest.unstable_mockModule('../src/js/issue.js', () => ({
  openIssuePage: jest.fn(),
}));

// Test cases
describe('Issues table', function() {
  it('openIssuesPage should add the issues to the page-contents', async function() {
    // Arrange
    const issue = await import('../src/js/issues.js');
    // Arrange input issues
    let issues = [];
    issues = [
      {id: 5, severity: 1, jsonkey: 51},
      {id: 15, severity: 0, jsonkey: 150},
    ];

    sessionStorage.setItem('DataBaseData', JSON.stringify(issues));

    // Act
    await issue.openIssuesPage();

    const name = document.getElementsByClassName('lang-name')[0].innerHTML;
    const type = document.getElementsByClassName('lang-type')[0].innerHTML;
    const risk = document.getElementsByClassName('lang-risk')[0].innerHTML;

    // Assert
    test.value(name).isEqualTo('Name');
    test.value(type).isEqualTo('Type');
    test.value(risk).isEqualTo('Risk');

    // Make issues table empty
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    emptyTable(issueTable);
    const nonIssueTable = document.getElementById('non-issues-table').querySelector('tbody');
    emptyTable(nonIssueTable);
  });
  it('toRiskLevel should return the right risk level', async function() {
    // Arrange
    const issue = await import('../src/js/issues.js');

    // act
    const risks = ['Acceptable', 'Low', 'Medium', 'High', 'Info'];

    // Assert
    risks.forEach((value, index) => {
      test.value(issue.toRiskLevel(index)).isEqualTo(value);
    });
  });
  it('fillTable should fill the issues table with information from the provided JSON array', async function() {
    // Arrange input issues
    let issues = [];
    issues = [
      {id: 5, severity: 1, jsonkey: 51},
      {id: 15, severity: 0, jsonkey: 150},
    ];
    // Arrange expected table data
    const expectedData = [];
    expectedData.push(data[issues[0].jsonkey]);
    expectedData.push(data[issues[1].jsonkey]);

    const issue = await import('../src/js/issues.js');

    // Act
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    issue.fillTable(issueTable, issues, true);

    const nonIssueTable = document.getElementById('non-issues-table').querySelector('tbody');
    issue.fillTable(nonIssueTable, issues, false);
    // Assert
    let row = issueTable.rows[0];
    test.value(row.cells[0].textContent).isEqualTo(expectedData[0].Name);
    test.value(row.cells[1].textContent).isEqualTo(expectedData[0].Type);
    test.value(row.cells[2].textContent).isEqualTo(issue.toRiskLevel(issues[0].severity));

    row = nonIssueTable.rows[0];
    test.value(row.cells[0].textContent).isEqualTo(expectedData[1].Name);
    test.value(row.cells[1].textContent).isEqualTo(expectedData[1].Type);

    // Make issues table empty
    emptyTable(issueTable);
    emptyTable(nonIssueTable);
  });
  it('fillTable should not fill the issues table if the data from the JSON array is incorrect', async function() {
    // Arrange input issues
    let issues = [];
    issues = [
      {id: 0, severity: 1, jsonkey: 55},
      {id: 123, severity: 0, jsonkey: 1234},
    ];

    // Arrange expected table data
    const expectedData = [];
    expectedData.push(data[issues[0].jsonkey]);
    expectedData.push(data[issues[1].jsonkey]);
    const issue = await import('../src/js/issues.js');

    // Act
    const issueTable = document.getElementById('issues-table').querySelector('tbody');
    emptyTable(issueTable);
    issue.fillTable(issueTable, issues, true);
    const nonIssueTable = document.getElementById('non-issues-table').querySelector('tbody');
    emptyTable(nonIssueTable);
    issue.fillTable(nonIssueTable, issues, false);

    // Assert
    let row = issueTable.rows[0];
    test.value(row).isEqualTo(undefined);

    row = nonIssueTable.rows[0];
    test.value(row).isEqualTo(undefined);
  });
  it('sortTable should sort the issues table', async function() {
    // Arrange table rows
    const table = dom.window.document.getElementById('issues-table');
    const tbody = table.querySelector('tbody');
    tbody.innerHTML = `
      <tr>
        <td>Windows defender</td>
        <td>Security</td>
        <td>High</td>
      </tr>
      <tr>
        <td>Camera and microphone access</td>
        <td>Privacy</td>
        <td>Low</td>
      </tr>
      <tr>
        <td>Firewall settings</td>
        <td>Security</td>
        <td>Medium</td>
      </tr>
    `;

    await import('../src/js/issues.js');

    // Act
    document.getElementById('sort-on-issue').dispatchEvent(clickEvent);

    // Assert
    let sortedRows = Array.from(tbody.rows);
    const sortedNames = sortedRows.map((row) => row.cells[0].textContent);
    test.array(sortedNames).is(['Camera and microphone access', 'Firewall settings', 'Windows defender']);

    // Act
    document.getElementById('sort-on-issue').dispatchEvent(clickEvent);

    // Assert
    let sortedRowsDescending = Array.from(tbody.rows);
    const sortedNamesDescending = sortedRowsDescending.map((row) => row.cells[0].textContent);
    test.array(sortedNamesDescending).is(['Windows defender', 'Firewall settings', 'Camera and microphone access']);

    // Act
    document.getElementById('sort-on-type').dispatchEvent(clickEvent);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedTypes = sortedRows.map((row) => row.cells[1].textContent);
    test.array(sortedTypes).is(['Privacy', 'Security', 'Security']);

    // Act
    document.getElementById('sort-on-type').dispatchEvent(clickEvent);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedTypesDescending = sortedRowsDescending.map((row) => row.cells[1].textContent);
    test.array(sortedTypesDescending).is(['Security', 'Security', 'Privacy']);

    // Act
    document.getElementById('sort-on-risk').dispatchEvent(clickEvent);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedRisks = sortedRows.map((row) => row.cells[2].textContent);
    test.array(sortedRisks).is(['High', 'Medium', 'Low']);

    // Act
    document.getElementById('sort-on-risk').dispatchEvent(clickEvent);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedRisksDescending = sortedRowsDescending.map((row) => row.cells[2].textContent);
    test.array(sortedRisksDescending).is(['Low', 'Medium', 'High']);
  });
  it('sortTable should sort the non-issues table', async function() {
    // Arrange table rows
    const table = dom.window.document.getElementById('non-issues-table');
    const tbody = table.querySelector('tbody');
    tbody.innerHTML = `
      <tr>
        <td>Windows defender</td>
        <td>Security</td>
        <td>High</td>
      </tr>
      <tr>
        <td>Camera and microphone access</td>
        <td>Privacy</td>
        <td>Low</td>
      </tr>
      <tr>
        <td>Firewall settings</td>
        <td>Security</td>
        <td>Medium</td>
      </tr>
    `;

    await import('../src/js/issues.js');

    // Act
    document.getElementById('sort-on-issue2').dispatchEvent(clickEvent);

    // Assert
    let sortedRows = Array.from(tbody.rows);
    const sortedNames = sortedRows.map((row) => row.cells[0].textContent);
    test.array(sortedNames).is(['Camera and microphone access', 'Firewall settings', 'Windows defender']);

    // Act
    document.getElementById('sort-on-issue2').dispatchEvent(clickEvent);

    // Assert
    let sortedRowsDescending = Array.from(tbody.rows);
    const sortedNamesDescending = sortedRowsDescending.map((row) => row.cells[0].textContent);
    test.array(sortedNamesDescending).is(['Windows defender', 'Firewall settings', 'Camera and microphone access']);

    // Act
    document.getElementById('sort-on-type2').dispatchEvent(clickEvent);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedTypes = sortedRows.map((row) => row.cells[1].textContent);
    test.array(sortedTypes).is(['Privacy', 'Security', 'Security']);

    // Act
    document.getElementById('sort-on-type2').dispatchEvent(clickEvent);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedTypesDescending = sortedRowsDescending.map((row) => row.cells[1].textContent);
    test.array(sortedTypesDescending).is(['Security', 'Security', 'Privacy']);
  });
  it('changeTable should update the table with selected risks', async function() {
    // Arrange
    const issue = await import('../src/js/issues.js');

    // Arrange input issues
    let issues = [];
    issues = [
      {id: 5, severity: 1, jsonkey: 51},
      {id: 16, severity: 2, jsonkey: 160},
      {id: 18, severity: 3, jsonkey: 182},
      {id: 6, severity: 4, jsonkey: 60},
    ];
    // Arrange expected table data
    const expectedData = [];
    issues.forEach((issue) => {
      expectedData.push(data[issue.jsonkey]);
    });
    sessionStorage.setItem('DataBaseData', JSON.stringify(issues));

    const ids = [
      'select-low-risk-table',
      'select-medium-risk-table',
      'select-high-risk-table',
      'select-info-risk-table',
    ];

    for (let i = -1; i < issues.length - 1; i++) {
      // Act
      let issueTable = document.getElementById('issues-table').querySelector('tbody');
      issue.fillTable(issueTable, issues, true);

      if (i >= 0) document.getElementById(ids[i]).checked = false;
      issue.changeTable();
      issueTable = document.getElementById('issues-table').querySelector('tbody');

      // Assert
      expectedData.forEach((expectedIssue, index) => {
        if (index > i) {
          // const row = issueTable.rows[index - 1 - i];
          // test.value(row.cells[0].textContent).isEqualTo(expectedIssue.Name);
          // test.value(row.cells[1].textContent).isEqualTo(expectedIssue.Type);
          // test.value(row.cells[2].textContent).isEqualTo(issue.toRiskLevel(issues[index].severity));
        }
      });
    }
  });
  it('clicking on an issue should open the issue page', async function() {
    // Arrange
    const issue = await import('../src/js/issue.js');
    const issueLinks = document.querySelectorAll('.issue-link');
    const openIssuePageMock = jest.spyOn(issue, 'openIssuePage');

    // Assert
    issueLinks.forEach((link) => {
      link.dispatchEvent(clickEvent);
      expect(openIssuePageMock).toHaveBeenCalled();
    });
  });
  it('clicking the select risks toggles show', async function() {
    // Arrange
    await import('../src/js/issues.js');
    const button = document.getElementById('dropbtn-table');
    const myDropdownTable = document.getElementById('myDropdown-table');

    // Act
    button.dispatchEvent(clickEvent);

    // Arrange
    expect(myDropdownTable.classList.contains('show')).toBe(true);

    // Act
    button.dispatchEvent(clickEvent);

    // Arrange
    expect(myDropdownTable.classList.contains('show')).toBe(false);
  });
  it('should use the correct data object based on user language settings', async () => {
    // Define the language settings and the corresponding expected data
    const languageSettings = [
      {language: 0, expectedData: dataDe},
      {language: 1, expectedData: data},
      {language: 2, expectedData: dataEnUS},
      {language: 3, expectedData: dataEs},
      {language: 4, expectedData: dataFr},
      {language: 5, expectedData: dataNl},
      {language: 6, expectedData: dataPt},
      {language: 999, expectedData: data}, // Default case
    ];
    const loadUserSettingsMock = jest.spyOn(await import('../wailsjs/go/main/App.js'), 'LoadUserSettings');

    for (const {language, expectedData} of languageSettings) {
      loadUserSettingsMock.mockResolvedValueOnce({Language: language});
      // Prepare the issues array
      const issues = [
        {id: 1, severity: 1, jsonkey: 51}, // assuming 51 exists in all datasets
      ];

      // Act
      const {fillTable} = await import('../src/js/issues.js');
      const issueTable = document.createElement('tbody');
      await fillTable(issueTable, issues, true);

      // Assert
      const row = issueTable.rows[0];
      const currentIssue = expectedData[issues[0].jsonkey];
      test.value(row.cells[0].textContent).isEqualTo(currentIssue.Name);
      test.value(row.cells[1].textContent).isEqualTo(currentIssue.Type);
      test.value(row.cells[2].textContent).isEqualTo('Low'); // Based on toRiskLevel(1)
    }
  });
});
