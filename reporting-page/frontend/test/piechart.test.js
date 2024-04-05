import 'jsdom-global/register.js';
import test from 'unit.js';
import { JSDOM } from "jsdom";
import { PieChart } from '../src/js/piechart.js';
import { RiskCounters } from '../src/js/risk-counters.js';

// Mock page
const dom = new JSDOM(`
  <div class="data-column piechart">
    <canvas id="pieChart"></canvas>
  </div>
`, {
url: 'http://localhost'
});
global.document = dom.window.document
global.window = dom.window

// test cases
describe("Risk level distribution piechart", function() {
  // arrange
  let rc = new RiskCounters(true);
  let p = new PieChart(undefined,rc);
  it("getData should fill the piechart with the correct data", function() {
    // arrange
    const expectedXValues = ["No risk", "Low risk", "Medium risk", "High risk"];
    const expectedYValues = [4,3,2,1];
    const expectedColors = ["rgb(255, 255, 0)","rgb(255, 0, 0)","rgb(0, 0, 255)","rgb(0, 255, 255)"];

    const mockRiskCounters = {
      highRiskColor : "rgb(0, 255, 255)",
      mediumRiskColor : "rgb(0, 0, 255)",
      lowRiskColor : "rgb(255, 0, 0)",
      noRiskColor : "rgb(255, 255, 0)",

      lastHighRisk : 1,
      lastMediumRisk : 2,
      lastLowRisk : 3,
      lastnoRisk : 4,
    };
    p = new PieChart(undefined,mockRiskCounters);

    // act 
    const resultData = p.GetData();

    // assert
    test.array(resultData.labels).is(expectedXValues);
    test.array(resultData.datasets[0].backgroundColor).is(expectedColors);
    test.array(resultData.datasets[0].data).is(expectedYValues);
  })
  it("getOptions should return the correct piechart options", function() {
    // arrange
    const expectedOptions = {
      maintainAspectRatio: false,
      title: {
        display: true,
        text: "Security Risks Overview"
      }
    };

    // act
    const resultOptions = p.GetOptions();

    // assert
    test.object(resultOptions).is(expectedOptions);
  })
});