import {ScanNow as scanNowGo, LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {openHomePage} from './home.js';
import {
  WindowShow as windowShow,
  WindowMaximise as windowMaximise,
  LogPrint as logPrint} from '../../wailsjs/runtime/runtime.js';
import * as rc from './risk-counters.js';
import {updateRiskCounter} from './risk-counters.js';
import data from '../databases/database.en-GB.json' assert { type: 'json' };

let isFirstScan = true;
/**
 * Initiates a scan and handles the result.
 *
 * @param {boolean} dialogPresent - Indicates whether a dialog is present during the scan.
 */
export async function scanTest(dialogPresent) {
  if (sessionStorage.getItem('isScanning') === null ||
      sessionStorage.getItem('isScanning') === undefined ||
      sessionStorage.getItem('isScanning') === 'false') {
    sessionStorage.setItem('isScanning', 'true');
    try {
      await new Promise((resolve, reject) => {
        scanNowGo(dialogPresent)
          .then(async (scanResult) => {
          // Handle the scan result
          // For example, save it in session storage
            sessionStorage.setItem('ScanResult', JSON.stringify(scanResult));
            // Set severities in session storage
            await setAllSeverities(scanResult);
            // set the detected windows version
            const windowsVersion = scanResult.find((i) => i.issue_id === 18);
            sessionStorage.setItem('WindowsVersion', windowsVersion.result[0]);
            // Resolve the promise with the scan result
            resolve(scanResult);
            sessionStorage.setItem('isScanning', 'false');
          })
          .catch((err) => {
          // Log any errors from scanNowGo
            logError('Error in scanNowGo: ' + err);
            // Reject the promise with the error
            reject(err);
          });
      });
      // Perform other actions after scanTest is complete
      windowShow();
      logPrint(sessionStorage.getItem('ScanResult'));
    } catch (err) {
      // Handle any errors that occurred during scanTest or subsequent actions
      logError('Error in scanTest: ' + err);
    }
  }
}

// Check if scanTest has already been called before
if (sessionStorage.getItem('scanTest') === null || sessionStorage.getItem('scanTest') === undefined) {
  // Call scanTest() only if it hasn't been called before
  scanTest(false).then((r) => {});

  // Set the flag in sessionStorage to indicate that scanTest has been called
  sessionStorage.setItem('scanTest', 'called');
}

// counts the occurrences of each level: 0 = acceptable, 1 = low, 2 = medium, 3 = high
const countOccurrences = (severities, level) => severities.filter((item) => data[item.issue_id][item.result_id].Severity === level).length;

/** Sets the severities collected from the checks and database in session storage of all types
 *
 * @param {Check[]} input Checks to get severities from
 */
async function setAllSeverities(input) {
  sessionStorage.setItem('IssuesSorting', JSON.stringify(
    {
      'column': 2,
      'direction': 'descending',
    },
  ));
  sessionStorage.setItem('IssuesFilter', JSON.stringify(
    {'high': 1, 'medium': 1, 'low': 1, 'acceptable': 1, 'info': 1},
  ));
  await setSeverities(input, '');
  await setSeverities(input, 'Security');
  await setSeverities(input, 'Privacy');
  if (isFirstScan) {
    openHomePage();
    windowMaximise();
    isFirstScan = false;
  }
}

/** Sets the severities collected from the scan in session storage
 *
 * @param {Checks[]} input Scan data retrieved from a scan
 * @param {string} type Type of issues to set the severities of in session storage
 */
async function setSeverities(input, type) {
  try {
    input = input.filter((item) => item.result_id !== -1);
    if (type !== '') {
      input = input.filter((item) => data[item.issue_id] !== undefined);
      input = input.filter((item) => data[item.issue_id].Type === type);
    }
    const info = countOccurrences(input, 4);
    const high = countOccurrences(input, 3);
    const medium = countOccurrences(input, 2);
    const low = countOccurrences(input, 1);
    const acceptable = countOccurrences(input, 0);
    if (sessionStorage.getItem(type + 'RiskCounters') === null ||
        sessionStorage.getItem(type + 'RiskCounters') === undefined) {
      sessionStorage.setItem(type + 'RiskCounters',
        JSON.stringify(new rc.RiskCounters(high, medium, low, info, acceptable)));
    } else {
      let riskCounter = JSON.parse(sessionStorage.getItem(type + 'RiskCounters'));
      riskCounter = updateRiskCounter(riskCounter, high, medium, low, info, acceptable);
      sessionStorage.setItem(type + 'RiskCounters', JSON.stringify(riskCounter));
    }
  } catch (err) {
    /* istanbul ignore next */
    logError(err);
  }
}
