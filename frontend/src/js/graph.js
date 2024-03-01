import "../css/home.css";
import "../css/dashboard.css";
import "../css/color-palette.css";

// creates labels on x-axis, for now just numbers.
function getData() {
    var labels = [];
    for (var i = 1; i <= Math.min(allNoRisks.length, graphShowAmount); i++) {
        labels.push(i);
    }
    
    var noRiskData = {
        label: 'Safe issues',
        data: allNoRisks.slice(Math.max(allNoRisks.length - graphShowAmount, 0)),
        backgroundColor: noRiskColor,
    }
    
    var lowRiskData = {
        label: 'Low risk issues',
        data: allLowRisks.slice(Math.max(allLowRisks.length - graphShowAmount, 0)),
        backgroundColor: lowRiskColor,
    }
    
    var mediumRiskData = {
        label: 'Medium risk issues',
        data: allMediumRisks.slice(Math.max(allMediumRisks.length - graphShowAmount, 0)),
        backgroundColor: mediumRiskColor,
    }
    
    var highRiskData = {
        label: 'High risk issues',
        data: allHighRisks.slice(Math.max(allHighRisks.length - graphShowAmount, 0)),
        backgroundColor: highRiskColor,
    }
    
    var datasets = [];
    console.log(graphShowLowRisks);
    if (graphShowNoRisks) {
        datasets.push(noRiskData)
    }
    if (graphShowLowRisks) {
        datasets.push(lowRiskData)
    }
    if (graphShowMediumRisks) {
        datasets.push(mediumRiskData)
    }
    if (graphShowHighRisks) {
        datasets.push(highRiskData)
    }

    var data = {
        labels: labels,
        datasets: datasets
      };
    
    return data;
}

var barChart = new Chart("graph", {

    type: 'bar',
    // The data for our dataset
    data: getData(),

    // Configuration options go here
    options: {
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
    }
});

function ChangeGraph() {
    graphShowAmount = document.getElementById('graph-interval').value;
    barChart.data = getData();
    barChart.update();
}

function ToggleHighRisks() {
    if (graphShowHighRisks) {
        graphShowHighRisks = false;
    } else {
        graphShowHighRisks = true;
    }
    ChangeGraph();
}

function ToggleMediumRisks() {
    if (graphShowMediumRisks) {
        graphShowMediumRisks = false;
    } else {
        graphShowMediumRisks = true;
    }
    ChangeGraph();
}

function ToggleLowRisks() {
    if (graphShowLowRisks) {
        graphShowLowRisks = false;
    } else {
        graphShowLowRisks = true;
    }
    ChangeGraph();
}

function ToggleNoRisks() {
    if (graphShowNoRisks) {
        graphShowNoRisks = false;
    } else {
        graphShowNoRisks = true;
    }
    ChangeGraph();
}

function GraphDropdown() {
    document.getElementById("myDropdown").classList.toggle("show");
  }
  
// Close the dropdown if the user clicks outside of it
// window.onclick = function(e) {
//   if (!e.target.matches('.dropbtn')) {
//       if (!e.target.matches('.dropdown-selector')) {
//           var myDropdown = document.getElementById("myDropdown");
//           console.log("true");
//           if (myDropdown.classList.contains('show')) {
//               myDropdown.classList.remove('show');
//       }
//     }
//   }
// }