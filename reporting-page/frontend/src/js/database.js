import {ScanNow as scanNowGo, LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {GetDataBaseData as getDataBaseData} from '../../wailsjs/go/main/DataBase.js';
import {openHomePage} from './home.js';
import {
  WindowShow as windowShow, 
  WindowMaximise as windowMaximise,
  LogPrint as logPrint} from '../../wailsjs/runtime/runtime.js';
import * as rc from './risk-counters.js';
import {updateRiskcounter} from './risk-counters.js';
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
          await setSeverities(scanResult);
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
    windowShow();
    windowMaximise();
    logPrint(sessionStorage.getItem('ScanResult'));
  } catch (err) {
    // Handle any errors that occurred during scanTest or subsequent actions
    logError('Error in scanTest: ' + err);
  }
}

// Check if scanTest has already been called before
if (sessionStorage.getItem('scanTest') === null) {
  // Call scanTest() only if it hasn't been called before
  scanTest();

  // Set the flag in sessionStorage to indicate that scanTest has been called
  sessionStorage.setItem('scanTest', 'called');
}

// counts the occurences of each level: 0 = acceptable, 1 = low, 2 = medium, 3 = high
const countOccurences = (severities, level) => severities.filter((item) => item.severity === level).length;

/** Sets the severities collected from the checks and database in session storage
 *
 * @param {Check[]} input Checks to get severities from
 * @param {int[]} ids List of result ids to get corresponding severities
 */
async function setSeverities(input) {
  try {
    const result = await getDataBaseData(input);
    sessionStorage.setItem('DataBaseData', JSON.stringify(result));
    const high = countOccurences(result, 3);
    const medium = countOccurences(result, 2);
    const low = countOccurences(result, 1);
    const acceptable = countOccurences(result, 0);
    if (sessionStorage.getItem('RiskCounters') === null || sessionStorage.getItem('RiskCounters') === undefined) {
      sessionStorage.setItem('RiskCounters', JSON.stringify(new rc.RiskCounters(high, medium, low, acceptable)));
      openHomePage();
    } else {
      let riskCounter = JSON.parse(sessionStorage.getItem('RiskCounters'));
      console.log(riskCounter);
      riskCounter = updateRiskcounter(riskCounter, high, medium, low, acceptable);
      sessionStorage.setItem('RiskCounters', JSON.stringify(riskCounter));
    }
  } catch (err) {
    logError(err);
  }
}
