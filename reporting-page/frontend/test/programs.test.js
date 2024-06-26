import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, storageMock} from './mock.js';

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
  while (table.rows.length > 0) {
    table.deleteRow(0);
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
    case 'Programs.ProgramTable':
      myResolve('Programs table');
    case 'Programs.Name':
      myResolve('Name');
    case 'Programs.Version':
      myResolve('Version');
    case 'Programs.Search':
      myResolve('Search...');
    case 'Programs.ReloadText':
      myResolve('ReloadText');
    case 'Programs.ReloadButton':
      myResolve('Reload page');
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

// Test cases
describe('Programs table', function() {
  it('openProgramsPage should add the programs to the page-contents', async function() {
    // Arrange
    const programs = await import('../src/js/programs.js');
    // Arrange input issues
    const list =
      { // Programs, level 0
        issue_id: 43,
        result_id: 0,
        result: [
          'Program 1 | 1.0.0',
          'Program 2 | 2.0.0',
        ],
      };

    sessionStorage.setItem('ProgramList', JSON.stringify(list));

    // Act
    await programs.openProgramsPage();

    const name = document.getElementsByClassName('lang-name')[0].innerHTML;
    const version = document.getElementsByClassName('lang-version')[0].innerHTML;

    // Assert
    test.value(name).isEqualTo('Name');
    test.value(version).isEqualTo('Version');

    // Make programs table empty
    const programsTable = document.getElementById('program-table').querySelector('tbody');
    emptyTable(programsTable);
  });
  it('fillProgramTable should fill the programs table with information from the provided JSON array', async function() {
    // Arrange input issues
    const issues = [
      { // Privacy, level 0
        issue_id: 21,
        result_id: 0,
        result: [],
      },
      { // Programs, level 0
        issue_id: 43,
        result_id: 0,
        result: [
          'Program 1 | 1.0.0',
          'Program 2 | 2.0.0',
        ],
      },
    ]; ;
    const foundObject = await issues.find((obj) => obj.issue_id === 43);
    const programsTable = document.getElementById('program-table').querySelector('tbody');
    const programs = await import('../src/js/programs.js');

    // Act
    programs.fillProgamTable(programsTable, foundObject.result);

    // Assert
    const row = programsTable.rows[0];
    test.value(row.cells[0].textContent).isEqualTo(foundObject.result[0].split(' | ')[0]);
    test.value(row.cells[1].textContent).isEqualTo(foundObject.result[0].split(' | ')[1]);

    // Make programs table empty
    emptyTable(programsTable);
  });
  it('sortProgramTable should sort the programs table', async function() {
    // Arrange table rows
    const table = dom.window.document.getElementById('program-table');
    const tbody = table.querySelector('tbody');
    tbody.innerHTML = `
      <tr data-severity="3">
        <td>Program B</td>
        <td>3</td>
      </tr>
      <tr data-severity="1">
        <td>Program C</td>
        <td>2</td>
      </tr>
      <tr data-severity="2">
        <td>Program A</td>
        <td>1</td>
      </tr>
    `;

    const programs = await import('../src/js/programs.js');

    // Act
    programs.sortProgramTable(tbody, 'ascending');

    // Assert
    const sortedRows = Array.from(tbody.rows);
    const sortedNames = sortedRows.map((row) => row.cells[0].textContent);
    test.array(sortedNames).is(['Program A', 'Program B', 'Program C']);

    // Act
    programs.sortProgramTable(tbody, 'descending');

    // Assert
    const sortedRowsDescending = Array.from(tbody.rows);
    const sortedNamesDescending = sortedRowsDescending.map((row) => row.cells[0].textContent);
    test.array(sortedNamesDescending).is(['Program C', 'Program B', 'Program A']);
  });
  it('searchTable should filter the table based on the search query', async function() {
    // Arrange table rows
    const table = dom.window.document.getElementById('program-table');
    const tbody = table.querySelector('tbody');
    tbody.innerHTML = `
      <tr>
        <td>Program A</td>
        <td>1</td>
      </tr>
      <tr>
        <td>Program B</td>
        <td>2</td>
      </tr>
      <tr>
        <td>Program C</td>
        <td>3</td>
      </tr>
    `;

    const programs = await import('../src/js/programs.js');

    // Act
    programs.searchTable(tbody, 'Program B');

    // Assert
    test.value(tbody.rows[0].style.display).isEqualTo('none');
    test.value(tbody.rows[1].style.display).isEqualTo('');

    // Act
    programs.searchTable(tbody, 'Program');

    // Assert
    test.value(tbody.rows[0].style.display).isEqualTo('');
    test.value(tbody.rows[1].style.display).isEqualTo('');
  });
});
