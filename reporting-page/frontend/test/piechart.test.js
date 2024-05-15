import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {mockChart} from './mock.js';
import {RiskCounters} from '../src/js/risk-counters.js';
import {jest} from '@jest/globals';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
  <div class="data-column pie-chart">
    <canvas id="pieChart"></canvas>
  </div>
`, {
  url: 'http://localhost',
});
global.document = dom.window.document;
global.window = dom.window;

/** Mock of getLocalizationString function
 *
 * @param {string} messageID - The ID of the message to be localized.
 * @return {string} The localized string.
 */
function mockGetLocalizationString(messageID) {
  switch (messageID) {
  case 'Dashboard.Acceptable':
    return 'Acceptable';
  case 'Dashboard.LowRisk':
    return 'Low';
  case 'Dashboard.MediumRisk':
    return 'Medium';
  case 'Dashboard.HighRisk':
    return 'High';
  case 'Dashboard.InfoRisk':
    return 'Info';
  case 'Dashboard.TotalRisksOverview':
    return 'Total Risks Overview';
  case 'Dashboard.SecurityRisksOverview':
    return 'Security Risks Overview';
  case 'Dashboard.PrivacyRisksOverview':
    return 'Privacy Risks Overview';
  }
}

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalizationString(input)),
}));

// Mock chart constructor
mockChart();


// test cases
describe('Risk level distribution piechart', function() {
  it('getData should fill the piechart with the correct data', async function() {
    // arrange
    const chart = await import('../src/js/piechart.js');

    // arrange
    const expectedXValues = ['Acceptable', 'Low', 'Medium', 'High', 'Info'];
    const expectedYValues = [4, 3, 2, 1, 0];
    const expectedColors = [
      'rgb(255, 255, 0)',
      'rgb(255, 0, 0)',
      'rgb(0, 0, 255)',
      'rgb(0, 255, 255)',
      'rgb(255,255,255)',
    ];
    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',
      infoColor: 'rgb(255,255,255)',

      lastInfoRisk: 0,
      lastHighRisk: 1,
      lastMediumRisk: 2,
      lastLowRisk: 3,
      lastNoRisk: 4,
    };

    const p = new chart.PieChart(undefined, mockRiskCounters);

    // act
    const result = await p.getData();

    // assert
    test.array(result.labels).is(expectedXValues);
    test.array(result.datasets[0].backgroundColor).is(expectedColors);
    test.array(result.datasets[0].data).is(expectedYValues);
  });
  it('getOptions should return the correct piechart options', async function() {
    // arrange
    const chart = await import('../src/js/piechart.js');
    const rc = new RiskCounters();
    const p = new chart.PieChart(undefined, rc);

    const titles = ['Total', 'Security', 'Privacy'];

    titles.forEach(async (title) => {
      // arrange
      const expectedOptions = {
        maintainAspectRatio: false,
        title: {
          display: true,
          text: (title + ' Risks Overview'),
        },
      };

      // act
      const result = await p.getOptions(title);

      // assert
      test.object(result).is(expectedOptions);
    });
  });
  it('Creating a piechart should call getOptions and getData', async function() {
    // arrange
    const piechart = await import('../src/js/piechart.js');
    const chart = await import('chart.js/auto');
    const rc = new RiskCounters();
    const getDataMock = jest.spyOn(piechart.PieChart.prototype, 'getData');
    const getOptionsMock = jest.spyOn(piechart.PieChart.prototype, 'getOptions');

    // act
    const p = new piechart.PieChart('pieChart', rc);
    await p.createPieChart();

    // assert
    expect(getDataMock).toHaveBeenCalled();
    expect(getOptionsMock).toHaveBeenCalled();
    expect(chart.Chart).toHaveBeenCalled();
    getDataMock.mockRestore();
    getOptionsMock.mockRestore();
  });
});
