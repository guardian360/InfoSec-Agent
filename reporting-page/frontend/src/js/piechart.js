// import * as this.rc from "./risk-counters.js"

import { RiskCounters } from "./risk-counters.js";

export class PieChart {
  pieChart;
  rc;
  /** Create a piechart showing the risk counters
   * 
   * @param {string=} canvas id of the canvas where the piechart would be placed
   * @param {RiskCounters} riskCounters Risk counters used to retrieve data to be put in the chart
   */
  constructor(canvas,riskCounters) {
    this.rc = riskCounters;
    if (canvas !== undefined) { this.CreatePieChart(canvas); }
  }

  /** Creates a pie chart for risks 
   * 
   * @param {string} canvas html canvas where pie chart will be placed
   */
  CreatePieChart(canvas) {
    this.pieChart = new Chart(canvas, {
      type: "doughnut",
      data: this.GetData(),
      options: this.GetOptions()
    });
  }

  /** Creates the data portion for a piechart using the different levels of risks */ 
  GetData() {
    var xValues = ["No risk", "Low risk", "Medium risk", "High risk"];
    var yValues = [this.rc.lastnoRisk, this.rc.lastLowRisk, this.rc.lastMediumRisk, this.rc.lastHighRisk];
    var barColors = [this.rc.noRiskColor, this.rc.lowRiskColor, this.rc.mediumRiskColor, this.rc.highRiskColor];

    return {
      labels: xValues,
      datasets: [{
        backgroundColor: barColors,
        data: yValues
      }]
    }
  }

  /** Creates the options for a pie chart 
   * 
   * @returns {options} Options for pie chart
   */ 
  GetOptions() {
    return {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: "Security Risks Overview"
      }
    }
  }
}


