import { ScanNow } from "../../wailsjs/go/main/Tray.js";
import { GetSeverities } from "../../wailsjs/go/main/DataBase.js";
import * as runTime from "../../wailsjs/runtime/runtime.js";
// var sqlite3 = require('sqlite3');

/** Call ScanNow in backend and store result in sessionStorage */
try {
  ScanNow()
    .then((result) => {
      // Update result with data back from App.Greet()
      sessionStorage.setItem("ScanResult",JSON.stringify(result));
      runTime.WindowShow();
      runTime.WindowMaximise();
      runTime.LogPrint(sessionStorage.getItem("ScanResult"));
      // console.log(JSON.parse(sessionStorage.getItem("ScanResult")));
      window.alert(sessionStorage.getItem("ScanResult"));
    })
    .catch((err) => {
      console.error(err);
    });
  } catch (err) {
    console.error(err);
}

function getSeverities(input, ids) {
  GetSeverities(input, ids)
    .then((result) => {
      
    })
}

function randomResultIDs(length) {
  let IDs = [];
  for (let i = 0; i < length; i++) {
    IDs[i] = Math.floor(Math.random() * 4)
  }
  return IDs;
}