import {ScanNow as scanNowGo} from '../../wailsjs/go/main/Tray.js';
import {GetDataBaseData as getDataBaseData} from '../../wailsjs/go/main/DataBase.js';
import {openHomePage} from './home.js';
import * as runTime from '../../wailsjs/runtime/runtime.js';
import * as rc from './risk-counters.js';

/** Call ScanNow in backend and store result in sessionStorage */
try {
  scanNowGo()
    .then((result) => {
      // place result in session storage
      sessionStorage.setItem('ScanResult', JSON.stringify(result));
      // place severities in session storage
      setSeverities(result);

      runTime.WindowShow();
      runTime.WindowMaximise();
      runTime.LogPrint(sessionStorage.getItem('ScanResult'));
    })
    .catch((err) => {
      console.error(err);
    });
} catch (err) {
  console.error(err);
}

// counts the occurences of each level: 0 = safe, 1 = low, 2 = medium, 3 = high
const countOccurences = (severities, level) => severities.filter(item => item.severity === level).length;

/** Sets the severities collected from the checks and database in session storage
 *
 * @param {Check[]} input Checks to get severities from
 * @param {int[]} ids List of result ids to get corresponding severities
 */
function setSeverities(input) {
  getDataBaseData(input)
    .then((result) => {
      sessionStorage.setItem("DataBaseData", JSON.stringify(result));
      let high = countOccurences(result, 3);
      let medium = countOccurences(result, 2);
      let low = countOccurences(result, 1);
      let safe = countOccurences(result, 0);
      sessionStorage.setItem("RiskCounters", JSON.stringify(new rc.RiskCounters(high, medium, low, safe)))
      openHomePage();
    })
    .catch((err) => {
      console.error(err);
    });
}

/** Get random result ids for each check
 *
 * @param {int} length Array length
 * @return {int} List of random result ids
 */
function randomResultIDs(length) {
  const IDs = [];
  for (let i = 0; i < length; i++) {
    IDs[i] = Math.floor(Math.random() * 4);
  }
  return IDs;
}
