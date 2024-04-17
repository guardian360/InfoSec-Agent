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

  /**
 * Creates the data portion for a piechart using the different levels of risks
 * @return {ChartData} The data for the pie chart.
 */
  async getData() {
    // const xValues = ['No risk', 'Low risk', 'Medium risk', 'High risk'];
    const xValues = [
      await getLocalizationString('Dashboard.Safe'), 
      await getLocalizationString('Dashboard.LowRisk'), 
      await getLocalizationString('Dashboard.MediumRisk'), 
      await getLocalizationString('Dashboard.HighRisk')
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

  /** Creates the options for a pie chart
   *
   * @return {options} Options for pie chart
   */
  async getOptions() {
    return {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: await getLocalizationString("Dashboard.SecurityRisksOverview"),
      },
    };
  }
}


