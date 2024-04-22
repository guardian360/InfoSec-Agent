import {ScanNow as scanNowGo, LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {GetDataBaseData as getDataBaseData} from '../../wailsjs/go/main/DataBase.js';
import {openHomePage} from './home.js';
import * as runTime from '../../wailsjs/runtime/runtime.js';
import * as rc from './risk-counters.js';
import {updateRiskcounter} from './risk-counters.js';
import data from '../database.json' assert { type: 'json' };
/** Call ScanNow in backend and store result in sessionStorage */
export async function scanTest() {
  try {
    await new Promise((resolve, reject) => {
      scanNowGo()
        .then(async (scanResult) => {
          // Handle the scan result
          // For example, save it in session storage
          sessionStorage.setItem('ScanResult', JSON.stringify(scanResult));
          // Set severities in session storage
          await setAllSeverities(scanResult);
          // Resolve the promise with the scan result
          resolve(scanResult);
        })
        .catch((err) => {
          // Log any errors from scanNowGo
          logError('Error in scanNowGo: ' + err);
          // Reject the promise with the error
          reject(err);
        });
    });

    // Perform other actions after scanTest is complete
    runTime.WindowShow();
    runTime.WindowMaximise();
    runTime.LogPrint(sessionStorage.getItem('ScanResult'));
  } catch (err) {
    // Handle any errors that occurred during scanTest or subsequent actions
    logError('Error in scanTest: ' + err);
  }
}

// Check if scanTest has already been called before
if (sessionStorage.getItem('scanTest') === null) {
  // Call scanTest() only if it hasn't been called before
  scanTest().then(r => {});

  // Set the flag in sessionStorage to indicate that scanTest has been called
  sessionStorage.setItem('scanTest', 'called');
}

// counts the occurences of each level: 0 = acceptable, 1 = low, 2 = medium, 3 = high
const countOccurences = (severities, level) => severities.filter((item) => item.severity === level).length;

/** Sets the severities collected from the checks and database in session storage of all types
 *
 * @param {Check[]} input Checks to get severities from
 */
async function setAllSeverities(input) {
  const result = await getDataBaseData(input);
  sessionStorage.setItem('DataBaseData', JSON.stringify(result));
  await setSeverities(result, '');
  await setSeverities(result, 'Security');
  await setSeverities(result, 'Privacy');
}

/** Sets the severities collected from the database in session storage
 *
 * @param {DataBaseData[]} input DataBaseData retrieved from database
 * @param {string} type Type of issues to set the severities of in session storage
 */
async function setSeverities(input, type) {
  try {
    if (type !== '') {
      input = input.filter((item) => data[item.jsonkey] !== undefined);
      input = input.filter((item) => data[item.jsonkey].Type === type);
    }
    const high = countOccurences(input, 3);
    const medium = countOccurences(input, 2);
    const low = countOccurences(input, 1);
    const acceptable = countOccurences(input, 0);
    if (sessionStorage.getItem(type + 'RiskCounters') === null ||
        sessionStorage.getItem(type + 'RiskCounters') === undefined) {
      sessionStorage.setItem(type + 'RiskCounters', JSON.stringify(new rc.RiskCounters(high, medium, low, acceptable)));
      openHomePage();
    } else {
      let riskCounter = JSON.parse(sessionStorage.getItem(type + 'RiskCounters'));
      riskCounter = updateRiskcounter(riskCounter, high, medium, low, acceptable);
      sessionStorage.setItem(type + 'RiskCounters', JSON.stringify(riskCounter));
    }
  } catch (err) {
    logError(err);
  }
}
