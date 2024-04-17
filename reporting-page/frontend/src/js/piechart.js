import {getLocalizationString} from './localize.js';

/**
 * Represents a PieChart object for displaying risk counters.
 */
export class PieChart {
  pieChart;
  rc;
  /** Create a piechart showing the risk counters
   *
   * @param {string=} canvas id of the canvas where the piechart would be placed
   * @param {RiskCounters} riskCounters Risk counters used to retrieve data to be put in the chart
   */
  constructor(canvas, riskCounters) {
    this.rc = riskCounters;
    if (canvas !== undefined) {
      this.createPieChart(canvas);
    }
  }

  /** Creates a pie chart for risks
   *
   * @param {string} canvas html canvas where pie chart will be placed
   */
  async createPieChart(canvas) {
    this.pieChart = new Chart(canvas, {
      type: 'doughnut',
      data: await this.getData(),
      options: await this.getOptions(),
      overrides: {
        plugins: {
          legend: {
            display: true,
            position: 'left',
          },
        },
      },
    });
  }

  /** Creates data for a pie chart using different levels of risks
   *
   * @param {*} getString Function to retrieve localized text
   * @return {ChartData} The data for the pie chart
   */
  async getData(getString = getLocalizationString) {
    const xValues = [
      await getString('Dashboard.Safe'),
      await getString('Dashboard.LowRisk'),
      await getString('Dashboard.MediumRisk'),
      await getString('Dashboard.HighRisk'),
    ];
    const yValues = [this.rc.lastnoRisk, this.rc.lastLowRisk, this.rc.lastMediumRisk, this.rc.lastHighRisk];
    const barColors = [this.rc.noRiskColor, this.rc.lowRiskColor, this.rc.mediumRiskColor, this.rc.highRiskColor];

    return {
      labels: xValues,
      datasets: [{
        backgroundColor: barColors,
        data: yValues,
      }],
    };
  }

  /** Creates options for a pie chart
   *
   * @param {*} getString Function to retrieve localized text
   * @return {ChartData} The options for the pie chart
   */
  async getOptions(getString = getLocalizationString) {
    return {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: await getString('Dashboard.SecurityRisksOverview'),
      },
    };
  }
}


