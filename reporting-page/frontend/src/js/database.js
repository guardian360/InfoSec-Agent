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

// function UseDatabase() {
//     let db = new sqlite3.Database('database.db', sqlite3.OPEN_READWRITE, (err) => {
//         if (err) {
//           console.error(err.message);
//         }
//         console.log('Connected to the chinook database.');
//       });
    
//     // insert one row 
//     db.run(`INSERT INTO issues(Issue ID, Result ID, Severity, JSON Key) VALUES(?,?,?,?)`, ['1,1,1,"PasswordManager"'], function(err) {
//         if (err) {
//           return console.log(err.message);
//         }
//         // get the last insert id
//         console.log(`A row has been inserted with rowid ${this.lastID}`);
//       });
    
//     db.close();
// }
