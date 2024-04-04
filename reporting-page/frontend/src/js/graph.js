/**
 * Represents a graph for displaying risk counters.
 */
export class Graph {
  graphShowHighRisks = true;
  graphShowMediumRisks = true;
  graphShowLowRisks = true;
  graphShowNoRisks = true;

  graphShowAmount = document.getElementById('graph-interval').value;

  barChart;
  rc;
  /** Create a bar chart showing the risk counters
   *
   * @param {string=} canvas Id of the canvas where the bar chart would be placed
   * @param {RiskCounters} riskCounters Risk counters used to retrieve data to be put in the chart
   */
  constructor(canvas, riskCounters) {
    this.rc = riskCounters;
    if (canvas !== undefined) this.CreateGraphChart(canvas);
  }

  /** Creates a graph in the form of a bar chart for risks
   *
   * @param {string} canvas html canvas where bar chart will be placed
   */
  CreateGraphChart(canvas) {
    this.barChart = new Chart(canvas, {
      type: 'bar',
      data: this.GetData(), // The data for our dataset
      options: this.GetOptions(), // Configuration options go here
    });
  }

  /** Updates the graph, should be called after a change in graph properties */
  ChangeGraph() {
    this.graphShowAmount = document.getElementById('graph-interval').value;
    this.barChart.data = this.GetData();
    console.log(this.graphShowAmount);
    this.barChart.update();
  }

  /** Toggles a risks to show in the graph
   *
   * @param {string} category Category corresponding to risk
   * @param {boolean} [change=true] Changes graph after call, normally set to *true*
   */
  ToggleRisks(category, change = true) {
    switch (category) {
    case 'high':
      this.graphShowHighRisks = !this.graphShowHighRisks;
      break;
    case 'medium':
      this.graphShowMediumRisks = !this.graphShowMediumRisks;
      break;
    case 'low':
      this.graphShowLowRisks = !this.graphShowLowRisks;
      break;
    case 'no':
      this.graphShowNoRisks = !this.graphShowNoRisks;
      break;
    default:
      break;
    }
    if (change) this.ChangeGraph();
  }

  /** toggles 'show' class on element with id:"myDropDown" */
  GraphDropdown() {
    document.getElementById('myDropdown').classList.toggle('show');
  }


  /** Creates the data portion for a graph using the different levels of risks
   *
   * @return {data} Data for graph chart
   */
  GetData() {
    /**
     * Labels created for the x-axis
     * @type {!Array<string>}
     */
    const labels = [];
    for (let i = 1; i <= Math.min(this.rc.allNoRisks.length, this.graphShowAmount); i++) {
      labels.push(i);
    }

    const noRiskData = {
      label: 'Safe issues',
      data: this.rc.allNoRisks.slice(Math.max(this.rc.allNoRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.noRiskColor,
    };

    const lowRiskData = {
      label: 'Low risk issues',
      data: this.rc.allLowRisks.slice(Math.max(this.rc.allLowRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.lowRiskColor,
    };

    const mediumRiskData = {
      label: 'Medium risk issues',
      data: this.rc.allMediumRisks.slice(Math.max(this.rc.allMediumRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.mediumRiskColor,
    };

    const highRiskData = {
      label: 'High risk issues',
      data: this.rc.allHighRisks.slice(Math.max(this.rc.allHighRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.highRiskColor,
    };

    const datasets = [];

    if (this.graphShowNoRisks) datasets.push(noRiskData);
    if (this.graphShowLowRisks) datasets.push(lowRiskData);
    if (this.graphShowMediumRisks) datasets.push(mediumRiskData);
    if (this.graphShowHighRisks) datasets.push(highRiskData);

    return {
      labels: labels,
      datasets: datasets,
    };
  }

  /** Create the options for a bar chart
   *
   * @return {Options} Options for graph chart
   */
  GetOptions() {
    return {
      scales: {
        xAxes: [{
          stacked: true,
        }],
        yAxes: [{
          stacked: true,
        }],
      },
      legend: {
        display: false,
      },
      maintainAspectRatio: false,
      categoryPercentage: 1,
    };
  }
}


