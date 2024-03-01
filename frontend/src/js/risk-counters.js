import "../css/home.css";
import "../css/dashboard.css";
import "../css/color-palette.css";

// this file should contain code to put the correct count for each risk assessment.

var allHighRisks = [1,2,3,4,5,6,5];
var allMediumRisks = [1,2,3,4,5,6,3];
var allLowRisks = [1,2,3,4,5,6,2];
var allNoRisks = [1,2,3,4,5,6,4];

var graphShowAmount = document.getElementById("graph-interval").value;

var graphShowHighRisks = true;
var graphShowMediumRisks = true;
var graphShowLowRisks = true;
var graphShowNoRisks = true;

// change counters according to collected data
document.getElementById("high-risk-counter").innerHTML = allHighRisks.slice(-1)[0];
document.getElementById("medium-risk-counter").innerHTML = allMediumRisks.slice(-1)[0];
document.getElementById("low-risk-counter").innerHTML = allLowRisks.slice(-1)[0];
document.getElementById("no-risk-counter").innerHTML = allNoRisks.slice(-1)[0];