var highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--highRiskColor');
var mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--mediumRiskColor');
var lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--lowRiskColor');
var noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--noRiskColor');

var xValues = ["No risk", "Low risk", "Medium risk", "High risk"];
var yValues = [allNoRisks.slice(-1)[0], allLowRisks.slice(-1)[0], allMediumRisks.slice(-1)[0], allHighRisks.slice(-1)[0]];
var barColors = [
    noRiskColor,
  lowRiskColor,
  mediumRiskColor,
  highRiskColor
];

new Chart("pieChart", {
    type: "doughnut",
    data: {
      labels: xValues,
      datasets: [{
        backgroundColor: barColors,
        data: yValues
      }]
    },
    options: {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: "Security Risks Overview"
      }
    }
  });