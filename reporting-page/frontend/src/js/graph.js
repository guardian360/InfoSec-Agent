import * as rc from "./risk-counters"

/** Creates the data portion for a graph using the different levels of risks 
 * 
 * @param {int} graphShowAmount Number of columns to show
 * @param {bool} graphShowHighRisks Show high risks or not
 * @param {bool} graphShowMediumRisks Show medium risks or not
 * @param {bool} graphShowLowRisks Show low risks or not
 * @param {bool} graphShowNoRisks Show safe risks or not
 * 
 * @returns {data} Data for graph chart
 */ 
export function GetData(graphShowAmount, graphShowHighRisks, graphShowMediumRisks, graphShowLowRisks, graphShowNoRisks) {
  /**
   * Labels created for the x-axis
   * @type {!Array<string>}
   */
  let labels = [];
  for (var i = 1; i <= Math.min(rc.allNoRisks.length, graphShowAmount); i++) {
    labels.push(i);
  }
  
  let noRiskData = {
    label: 'Safe issues',
    data: rc.allNoRisks.slice(Math.max(rc.allNoRisks.length - graphShowAmount, 0)),
    backgroundColor: rc.noRiskColor,
  };
  
  let lowRiskData = {
    label: 'Low risk issues',
    data: rc.allLowRisks.slice(Math.max(rc.allLowRisks.length - graphShowAmount, 0)),
    backgroundColor: rc.lowRiskColor,
  };
  
  let mediumRiskData = {
    label: 'Medium risk issues',
    data: rc.allMediumRisks.slice(Math.max(rc.allMediumRisks.length - graphShowAmount, 0)),
    backgroundColor: rc.mediumRiskColor,
  };
  
  let highRiskData = {
    label: 'High risk issues',
    data: rc.allHighRisks.slice(Math.max(rc.allHighRisks.length - graphShowAmount, 0)),
    backgroundColor: rc.highRiskColor,
  };
  
  let datasets = [];
  
  if (graphShowNoRisks) datasets.push(noRiskData);
  if (graphShowLowRisks) datasets.push(lowRiskData);
  if (graphShowMediumRisks) datasets.push(mediumRiskData);
  if (graphShowHighRisks) datasets.push(highRiskData); 

  return {
    labels: labels,
    datasets: datasets
  };
}

/** Create the options for a bar chart 
 * 
 * @returns {options} Options for graph chart
 */ 
export function GetOptions() {
  return {
    scales: {
      xAxes: [{
        stacked: true
      }],
      yAxes: [{
        stacked: true
      }]
    },
    legend: {
      display: false,
    },
    maintainAspectRatio: false,
    categoryPercentage: 1,
  };
}
 
