// this file should contain code to put the correct count for each risk assessment.

export class RiskCounters {
  highRiskColor;
  mediumRiskColor; 
  lowRiskColor;
  noRiskColor;

  allHighRisks = [1,2,3,4,5,6,2];
  allMediumRisks = [1,2,3,4,5,6,0];
  allLowRisks = [1,2,3,4,5,6,2];
  allNoRisks = [1,2,3,4,5,6,4];

  lastHighRisk = this.allHighRisks.slice(-1)[0];
  lastMediumRisk = this.allMediumRisks.slice(-1)[0];
  lastLowRisk = this.allLowRisks.slice(-1)[0];
  lastnoRisk = this.allNoRisks.slice(-1)[0];

  count = this.allHighRisks.length;

  /** Create the risk-Counters with the right colors
   * 
   * @param {boolean} [testing=false] Specifies if the class is being used in testing, normally set to *false*
   */
  constructor (testing=false) {
    if (testing) {
      this.highRiskColor = "rgb(0, 255, 255)";
      this.mediumRiskColor = "rgb(0, 0, 255)";
      this.lowRiskColor = "rgb(255, 0, 0)";
      this.noRiskColor = "rgb(255, 255, 0)";
    } else {
      this.highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--highRiskColor');
      this.mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--mediumRiskColor');
      this.lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--lowRiskColor');
      this.noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--noRiskColor');
    }
  } 
}
sessionStorage.setItem("RiskCounters",JSON.stringify(new RiskCounters()));


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
  
 