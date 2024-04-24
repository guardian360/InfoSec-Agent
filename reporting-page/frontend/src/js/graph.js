import {getLocalizationString} from './localize.js';

/**
 * Represents a graph for displaying risk counters.
 */
export class Graph {
  graphShowHighRisks = true;
  graphShowMediumRisks = true;
  graphShowLowRisks = true;
  graphShowNoRisks = true;
  graphShowInfoRisks = true;

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
    if (canvas !== undefined) this.createGraphChart(canvas).then(() => {});
  }

  /** Creates a graph in the form of a bar chart for risks
   *
   * @param {string} canvas html canvas where bar chart will be placed
   */
  async createGraphChart(canvas) {
    this.barChart = new Chart(canvas, {
      type: 'bar',
      data: await this.getData(), // The data for our dataset
      options: await this.getOptions(), // Configuration options go here
    });
  }

  /** Updates the graph, should be called after a change in graph properties */
  async changeGraph() {
    this.graphShowAmount = document.getElementById('graph-interval').value;
    this.barChart.data = await this.getData();
    this.barChart.update();
  }

  /** Toggles a risks to show in the graph
   *
   * @param {string} category Category corresponding to risk
   * @param {boolean} [change=true] Changes graph after call, normally set to *true*
   */
  toggleRisks(category, change = true) {
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
    case 'info':
      this.graphShowInfoRisks = !this.graphShowInfoRisks;
      break;
    default:
      break;
    }
    if (change) this.changeGraph().then(() => {});
  }

  /** toggles 'show' class on element with id:"myDropDown" */
  graphDropdown() {
    document.getElementById('myDropdown').classList.toggle('show');
  }

  /** Creates data for a bar chart
   *
   * @param {*} getString Function to retrieve localized text
   * @return {ChartData} The data for the bar chart
   */
  async getData(getString = getLocalizationString) {
    /**
     * Labels created for the x-axis
     * @type {!Array<string>}
     */
    const labels = [];
    for (let i = 1; i <= Math.min(this.rc.allNoRisks.length, this.graphShowAmount); i++) {
      labels.push(i.toString());
    }

    const noRiskData = {
      label: await getString('Dashboard.Safe'),
      data: this.rc.allNoRisks.slice(Math.max(this.rc.allNoRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.noRiskColor,
    };

    const lowRiskData = {
      label: await getString('Dashboard.LowRisk'),
      data: this.rc.allLowRisks.slice(Math.max(this.rc.allLowRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.lowRiskColor,
    };

    const mediumRiskData = {
      label: await getString('Dashboard.MediumRisk'),
      data: this.rc.allMediumRisks.slice(Math.max(this.rc.allMediumRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.mediumRiskColor,
    };

    const highRiskData = {
      label: await getString('Dashboard.HighRisk'),
      data: this.rc.allHighRisks.slice(Math.max(this.rc.allHighRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.highRiskColor,
    };

    const infoRiskData = {
      label: await getString('Dashboard.InfoRisk'),
      data: this.rc.allInfoRisks.slice(Math.max(this.rc.allInfoRisks.length - this.graphShowAmount, 0)),
      backgroundColor: this.rc.infoColor,
    };

    const datasets = [];

    if (this.graphShowNoRisks) datasets.push(noRiskData);
    if (this.graphShowLowRisks) datasets.push(lowRiskData);
    if (this.graphShowMediumRisks) datasets.push(mediumRiskData);
    if (this.graphShowHighRisks) datasets.push(highRiskData);
    if (this.graphShowInfoRisks) datasets.push(infoRiskData);

    return {
      labels: labels,
      datasets: datasets,
    };
  }

  /** Create the options for a bar chart
   *
   * @return {Options} Options for graph chart
   */
  getOptions() {
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


