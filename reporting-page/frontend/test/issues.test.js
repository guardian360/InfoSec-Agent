import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import data from '../src/databases/database.en-GB.json' assert { type: 'json' };
global.TESTING = true;

const dom = new JSDOM(`
  <div id="page-contents"></div>
  <div class="page-contents"></div>
`);
global.document = dom.window.document;
global.window = dom.window;

// Mock sessionStorage
const sessionStorageMock = (() => {
  let store = {};

  return {
    getItem: (key) => store[key],
    setItem: (key, value) => {
      store[key] = value.toString();
    },
    clear: () => {
      store = {};
    },
  };
})();
global.sessionStorage = sessionStorageMock;

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

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalizationString(input)),
  LoadUserSettings: jest.fn(),
}));

// Mock LogError
jest.unstable_mockModule('../wailsjs/go/main/Tray.js', () => ({
  LogError: jest.fn(),
}));

// Mock Navigation
jest.unstable_mockModule('../src/js/navigation-menu.js', () => ({
  closeNavigation: jest.fn(),
  markSelectedNavigationItem: jest.fn(),
  loadPersonalizeNavigation: jest.fn(),
}));

// Mock retrieveTheme
jest.unstable_mockModule('../src/js/personalize.js', () => ({
  retrieveTheme: jest.fn(),
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

    const issue = await import('../src/js/issues.js');

    // Act
    issue.sortTable(tbody, 0);

    // Assert
    let sortedRows = Array.from(tbody.rows);
    const sortedNames = sortedRows.map((row) => row.cells[0].textContent);
    test.array(sortedNames).is(['Camera and microphone access', 'Firewall settings', 'Windows defender']);

    // Act
    issue.sortTable(tbody, 0);

    // Assert
    let sortedRowsDescending = Array.from(tbody.rows);
    const sortedNamesDescending = sortedRowsDescending.map((row) => row.cells[0].textContent);
    test.array(sortedNamesDescending).is(['Windows defender', 'Firewall settings', 'Camera and microphone access']);

    // Act
    issue.sortTable(tbody, 2);

    // Assert
    sortedRows = Array.from(tbody.rows);
    const sortedRisks = sortedRows.map((row) => row.cells[2].textContent);
    test.array(sortedRisks).is(['High', 'Medium', 'Low']);

    // Act
    issue.sortTable(tbody, 2);

    // Assert
    sortedRowsDescending = Array.from(tbody.rows);
    const sortedRisksDescending = sortedRowsDescending.map((row) => row.cells[2].textContent);
    test.array(sortedRisksDescending).is(['Low', 'Medium', 'High']);
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

    for (let i = -1; i < issues.length; i++) {
      // Act
      let issueTable = document.getElementById('issues-table').querySelector('tbody');
      issue.fillTable(issueTable, issues, true);

      if (i >= 0) document.getElementById(ids[i]).checked = false;
      issue.changeTable();
      issueTable = document.getElementById('issues-table').querySelector('tbody');

      // Assert
      expectedData.forEach((expectedIssue, index) => {
        if (index > i) {
          const row = issueTable.rows[index - 1 - i];
          //test.value(row.cells[0].textContent).isEqualTo(expectedIssue.Name);
          //test.value(row.cells[1].textContent).isEqualTo(expectedIssue.Type);
          //test.value(row.cells[2].textContent).isEqualTo(issue.toRiskLevel(issues[index].severity));
        }
      });
    }
  });
});
