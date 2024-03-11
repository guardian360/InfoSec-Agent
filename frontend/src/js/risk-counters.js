// this file should contain code to put the correct count for each risk assessment.

export var highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--highRiskColor');
export var mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--mediumRiskColor');
export var lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--lowRiskColor');
export var noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--noRiskColor');

export var allHighRisks = [1,2,3,4,5,6,2];
export var allMediumRisks = [1,2,3,4,5,6,0];
export var allLowRisks = [1,2,3,4,5,6,2];
export var allNoRisks = [1,2,3,4,5,6,4];

export var lastHighRisk = allHighRisks.slice(-1)[0];
export var lastMediumRisk = allMediumRisks.slice(-1)[0];
export var lastLowRisk = allLowRisks.slice(-1)[0];
export var lastnoRisk = allNoRisks.slice(-1)[0];