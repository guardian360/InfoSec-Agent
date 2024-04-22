import {getLocalizationString} from './localize.js';

/**
 * Represents a PieChart object for displaying risk counters.
 */
export class PieChart {
  pieChart;
  rc;
  /** Create a pie-chart showing the risk counters
   *
   * @param {string=} canvas id of the canvas where the pie-chart would be placed
   * @param {RiskCounters} riskCounters Risk counters used to retrieve data to be put in the chart
   * @param {string} type Title for the type of risks shown in the chart
   */
  constructor(canvas, riskCounters, type) {
    this.rc = riskCounters;
    if (canvas !== undefined) {
      this.createPieChart(canvas, type).then((r) => {});
    }
  }

  /** Creates a pie chart for risks
   *
   * @param {string} canvas html canvas where pie chart will be placed
   * @param {string} type Title for the type of risks shown in the chart
   */
  async createPieChart(canvas, type) {
    this.pieChart = new Chart(canvas, {
      type: 'doughnut',
      data: await this.getData(),
      options: await this.getOptions(type),
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
    const yValues = [this.rc.lastNoRisk, this.rc.lastLowRisk, this.rc.lastMediumRisk, this.rc.lastHighRisk];
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
   * @param {string} type Title for the type of risks shown in the chart
   * @param {*} getString Function to retrieve localized text
   * @return {ChartData} The options for the pie chart
   */
  async getOptions(type, getString = getLocalizationString) {
    return {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: await getString('Dashboard.' + type + 'RisksOverview'),
      },
    };
  }
}


