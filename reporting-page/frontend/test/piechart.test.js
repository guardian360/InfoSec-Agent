import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {PieChart} from '../src/js/piechart.js';
import {RiskCounters} from '../src/js/risk-counters.js';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
  <div class="data-column piechart">
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
  case 'Dashboard.Safe':
    return 'Acceptable';
  case 'Dashboard.LowRisk':
    return 'Low';
  case 'Dashboard.MediumRisk':
    return 'Medium';
  case 'Dashboard.HighRisk':
    return 'High';
  case 'Dashboard.TotalRisksOverview':
    return 'Total Risks Overview';
  case 'Dashboard.SecurityRisksOverview':
    return 'Security Risks Overview';
  case 'Dashboard.PrivacyRisksOverview':
    return 'Privacy Risks Overview';
  }
}

// test cases
describe('Risk level distribution piechart', function() {
  // arrange
  const rc = new RiskCounters(true);
  let p = new PieChart(undefined, rc);
  it('getData should fill the piechart with the correct data', function() {
    // arrange
    const expectedXValues = ['Acceptable', 'Low', 'Medium', 'High'];
    const expectedYValues = [4, 3, 2, 1];
    const expectedColors = [
      'rgb(255, 255, 0)',
      'rgb(255, 0, 0)',
      'rgb(0, 0, 255)',
      'rgb(0, 255, 255)',
    ];
    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',

      lastHighRisk: 1,
      lastMediumRisk: 2,
      lastLowRisk: 3,
      lastnoRisk: 4,
    };

    p = new PieChart(undefined, mockRiskCounters);

    // act

    /** asynchronous function to call p.getData() */
    async function getData() {
      return await p.getData(mockGetLocalizationString);
    }

    // assert
    getData().then((result) =>{
      test.array(result.labels).is(expectedXValues);
      test.array(result.datasets[0].backgroundColor).is(expectedColors);
      test.array(result.datasets[0].data).is(expectedYValues);
    })
      .catch((error) =>
        console.log(error));
  });
  it('getOptions should return the correct piechart options', function() {
    // arrange
    const titles = ['Total', 'Security', 'Privacy'];

    titles.forEach((title) => {
      // arrange
      const expectedOptions = {
        maintainAspectRatio: false,
        title: {
          display: true,
          text: (title + ' Risks Overview'),
        },
      };

      // act

      /** asynchronous function to call p.getOptions() */
      async function getOptions() {
        return await p.getOptions(title, mockGetLocalizationString);
      }

      // assert
      getOptions().then((result) => {
        test.object(result).is(expectedOptions);
      })
        .catch((error) =>
          console.log(error));
    });
  });
});
