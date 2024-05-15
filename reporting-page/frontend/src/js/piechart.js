import {getLocalizationString} from './localize.js';
import {Chart} from 'chart.js/auto';

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
    this.createPieChart(canvas, type).then(() => {});
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
   * @return {ChartData} The data for the pie chart
   */
  async getData() {
    const xValues = [
      await getLocalizationString('Dashboard.Acceptable'),
      await getLocalizationString('Dashboard.LowRisk'),
      await getLocalizationString('Dashboard.MediumRisk'),
      await getLocalizationString('Dashboard.HighRisk'),
      await getLocalizationString('Dashboard.InfoRisk'),
    ];
    const yValues = [this.rc.lastNoRisk, this.rc.lastLowRisk,
      this.rc.lastMediumRisk, this.rc.lastHighRisk, this.rc.lastInfoRisk];
    const barColors = [this.rc.noRiskColor, this.rc.lowRiskColor,
      this.rc.mediumRiskColor, this.rc.highRiskColor, this.rc.infoColor];

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
   * @return {ChartData} The options for the pie chart
   */
  async getOptions(type) {
    return {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: await getLocalizationString('Dashboard.' + type + 'RisksOverview'),
      },
    };
  }
}


