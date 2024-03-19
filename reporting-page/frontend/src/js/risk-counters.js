// this file should contain code to put the correct count for each risk assessment.

export let highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--highRiskColor');
export let mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--mediumRiskColor');
export let lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--lowRiskColor');
export let noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--noRiskColor');

export let allHighRisks = [1,2,3,4,5,6,2];
export let allMediumRisks = [1,2,3,4,5,6,0];
export let allLowRisks = [1,2,3,4,5,6,2];
export let allNoRisks = [1,2,3,4,5,6,4];

export let lastHighRisk = allHighRisks.slice(-1)[0];
export let lastMediumRisk = allMediumRisks.slice(-1)[0];
export let lastLowRisk = allLowRisks.slice(-1)[0];
export let lastnoRisk = allNoRisks.slice(-1)[0];