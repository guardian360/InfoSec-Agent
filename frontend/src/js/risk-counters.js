
import "../css/home.css";
import "../css/dashboard.css";
import "../css/color-palette.css";
var highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--highRiskColor');
var mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--mediumRiskColor');
var lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--lowRiskColor');
var noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--noRiskColor');
// this file should contain code to put the correct count for each risk assessment.

var allHighRisks = [1,2,3,4,5,6,2];
var allMediumRisks = [1,2,3,4,5,6,0];
var allLowRisks = [1,2,3,4,5,6,2];
var allNoRisks = [1,2,3,4,5,6,4];

var graphShowAmount = document.getElementById("graph-interval").value;

var graphShowHighRisks = true;
var graphShowMediumRisks = true;
var graphShowLowRisks = true;
var graphShowNoRisks = true;

var lastHighRisk = allHighRisks.slice(-1)[0];
var lastMediumRisk = allMediumRisks.slice(-1)[0];
var lastLowRisk = allLowRisks.slice(-1)[0];
var lastnoRisk = allNoRisks.slice(-1)[0];

// change counters according to collected data
document.getElementById("high-risk-counter").innerHTML = lastHighRisk;
document.getElementById("medium-risk-counter").innerHTML = lastMediumRisk;
document.getElementById("low-risk-counter").innerHTML = lastLowRisk;
document.getElementById("no-risk-counter").innerHTML = lastnoRisk;

if (lastHighRisk > 1) {
    document.getElementById("security-status").innerHTML = "Critical";
    document.getElementById("security-status").style.backgroundColor = highRiskColor;
    document.getElementById("security-status").style.color = "rgb(255, 255, 255)";
} else if (lastMediumRisk > 1) {
    document.getElementById("security-status").innerHTML = "Medium concern";
    document.getElementById("security-status").style.backgroundColor = mediumRiskColor; 
    document.getElementById("security-status").style.color = "rgb(255, 255, 255)";
} else if (lastLowRisk > 1) {
    document.getElementById("security-status").innerHTML = "Light concern";
    document.getElementById("security-status").style.backgroundColor = lowRiskColor;
    document.getElementById("security-status").style.color = "rgb(0, 0, 0)";  
} else {
    document.getElementById("security-status").innerHTML = "Safe";
    document.getElementById("security-status").style.backgroundColor = noRiskColor;
    document.getElementById("security-status").style.color = "rgb(0, 0, 0)";  
}

document.getElementById("graph-interval").max = allNoRisks.length;