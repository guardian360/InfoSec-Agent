import { ScanNow } from "../../wailsjs/go/main/Tray.js";
import { GetAllSeverities } from "../../wailsjs/go/main/DataBase.js";
import * as runTime from "../../wailsjs/runtime/runtime.js";
import * as rc from "./risk-counters.js"
// var sqlite3 = require('sqlite3');

/** Call ScanNow in backend and store result in sessionStorage */
try {
  ScanNow()
    .then((result) => {
      // Update result with data back from App.Greet()
      sessionStorage.setItem("ScanResult",JSON.stringify(result));
      setSeverities(result,randomResultIDs(result.length))
      // runTime.WindowReloadApp();

      runTime.WindowShow();
      runTime.WindowMaximise();
      runTime.LogPrint(sessionStorage.getItem("ScanResult"));
      // console.log(JSON.parse(sessionStorage.getItem("ScanResult")));
      // window.alert(sessionStorage.getItem("ScanResult"));
    })
    .catch((err) => {
      console.error(err);
    });
  } catch (err) {
    console.error(err);
}

// function countOccurences(severities, riskLevel) {
//   const count = {};
//   for (const num of severities) {
//     count[riskLevel] = count[riskLevel] ? count[riskLevel] + 1 : 1;
//   }
//   return count[riskLevel]
// }
const countOccurences = (severities, riskLevel) => severities.filter(item => item ===riskLevel).length;

function setSeverities(input, ids) {
  GetAllSeverities(input, ids)
    .then((result) => {
      console.log(result);
      let high = countOccurences(result, 3);
      let medium = countOccurences(result, 2);
      let low = countOccurences(result, 1);
      let safe = countOccurences(result, 0);
      sessionStorage.setItem("RiskCounters",JSON.stringify(new rc.RiskCounters(high,medium,low,safe)))
    })
}

function randomResultIDs(length) {
  let IDs = [];
  for (let i = 0; i < length; i++) {
    IDs[i] = Math.floor(Math.random() * 4)
  }
  return IDs;
}