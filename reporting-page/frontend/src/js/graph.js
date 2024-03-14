import * as rc from "./risk-counters"

// Create the data portion for a graph using the different levels of risks
export function GetData(graphShowAmount, graphShowHighRisks, graphShowMediumRisks, graphShowLowRisks, graphShowNoRisks) {
    var labels = [];
    for (var i = 1; i <= Math.min(rc.allNoRisks.length, graphShowAmount); i++) {
        labels.push(i);
    }
    
    var noRiskData = {
        label: 'Safe issues',
        data: rc.allNoRisks.slice(Math.max(rc.allNoRisks.length - graphShowAmount, 0)),
        backgroundColor: rc.noRiskColor,
    };
    
    var lowRiskData = {
        label: 'Low risk issues',
        data: rc.allLowRisks.slice(Math.max(rc.allLowRisks.length - graphShowAmount, 0)),
        backgroundColor: rc.lowRiskColor,
    };
    
    var mediumRiskData = {
        label: 'Medium risk issues',
        data: rc.allMediumRisks.slice(Math.max(rc.allMediumRisks.length - graphShowAmount, 0)),
        backgroundColor: rc.mediumRiskColor,
    };
    
    var highRiskData = {
        label: 'High risk issues',
        data: rc.allHighRisks.slice(Math.max(rc.allHighRisks.length - graphShowAmount, 0)),
        backgroundColor: rc.highRiskColor,
    };
    
    var datasets = [];
   
    if (graphShowNoRisks) datasets.push(noRiskData);
    if (graphShowLowRisks) datasets.push(lowRiskData);
    if (graphShowMediumRisks) datasets.push(mediumRiskData);
    if (graphShowHighRisks) datasets.push(highRiskData);

    return {
        labels: labels,
        datasets: datasets
    };
}

// Create the options for a bar chart
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
 
// Close the dropdown if the user clicks outside of it
/*window.onclick = function(e) {
  if (!e.target.matches('.dropbtn')) {
      if (!e.target.matches('.dropdown-selector')) {
          var myDropdown = document.getElementById("myDropdown");
          console.log("true");
          if (myDropdown.classList.contains('show')) {
              myDropdown.classList.remove('show');
      }
    }
  }
}*/