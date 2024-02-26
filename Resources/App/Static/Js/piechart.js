const highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--highRiskColor');
const mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--mediumRiskColor');
const lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--lowRiskColor');
const noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--noRiskColor');

var xValues = ["No risk", "Low risk", "Medium risk", "High risk"];
var yValues = [4,4,4,4];
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