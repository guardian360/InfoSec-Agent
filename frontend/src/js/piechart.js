import * as rc from "./risk-counters"

// Create the data portion for a piechart using the different levels of risks
export function GetData() {
  var xValues = ["No risk", "Low risk", "Medium risk", "High risk"];
  var yValues = [rc.allNoRisks.slice(-1)[0], rc.allLowRisks.slice(-1)[0], rc.allMediumRisks.slice(-1)[0], rc.allHighRisks.slice(-1)[0]];
  var barColors = [
      rc.noRiskColor,
    rc.lowRiskColor,
    rc.mediumRiskColor,
    rc.highRiskColor
  ];

  return {
    labels: xValues,
    datasets: [{
      backgroundColor: barColors,
      data: yValues
    }]
  }
}

// Create the options for a pie chart
export function GetOptions() {
  return {
    maintainAspectRatio: false,
    title: {
      display: true,
      text: "Security Risks Overview"
    }
  }
}
