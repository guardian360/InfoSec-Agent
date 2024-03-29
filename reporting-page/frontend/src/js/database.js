import { ScanNow } from "../../wailsjs/go/main/Tray.js";
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

