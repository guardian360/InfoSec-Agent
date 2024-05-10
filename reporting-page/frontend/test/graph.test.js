import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {RiskCounters} from '../src/js/risk-counters.js';
import {jest} from '@jest/globals';
import {mockChart} from './mock.js';

global.TESTING = true;

// Mock page
const dom = new JSDOM(`
  <div class="graph-row">
  <div class="graph-column issues-graph-buttons">
    <H2>In this graph you are able to see the distribution of different issues 
    we have found over the past 5 times we ran a check.</H2>
    <div class="dropdown">
      <button class="dropbtn" id="dropbtn">Select Risks</button>
      <div class="dropdown-selector" id="myDropdown">
        <p><input type="checkbox" checked="true" value="true" id="select-high-risk">
          <label for="select-high-risk"> High risks</label><br>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-medium-risk">
          <label for="select-medium-risk"> Medium risks</label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-low-risk">
          <label for="select-low-risk"> Low risks</label>
        </p>
        <p><input type="checkbox" checked="true" value="true" id="select-no-risk">
          <label for="select-no-risk"> Safe</label>
        </p>
      </div>
    </div>
    <a class="interval-button"><p>Change interval</p><input type="number" value="5" id="graph-interval" min=1></a>
  </div>
  <div class="graph-column issues-graph">
    <canvas id="interval-graph"></canvas>
  </div>
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
  case 'Dashboard.InfoRisk':
    return 'Info';
  case 'Dashboard.SecurityRisksOverview':
    return 'Security Risks Overview';
  }
}

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalizationString(input)),
}));

// Mock Chart constructor
mockChart();

// test cases
describe('Risk graph', function() {
  it('toggleRisks should change which risk levels are shown in the risk graph', async function() {
    // arrange
    const graph = await import('../src/js/graph.js');
    const rc = new RiskCounters();
    const g = new graph.Graph('interval-graph', rc);
    await g.createGraphChart();

    // act
    g.toggleRisks('high');
    g.toggleRisks('medium');
    g.toggleRisks('low');
    g.toggleRisks('no');

    // assert
    test.value(g.graphShowHighRisks).isEqualTo(false);
    test.value(g.graphShowMediumRisks).isEqualTo(false);
    test.value(g.graphShowLowRisks).isEqualTo(false);
    test.value(g.graphShowNoRisks).isEqualTo(false);

    // act
    g.toggleRisks('high');
    g.toggleRisks('medium');
    g.toggleRisks('low');
    g.toggleRisks('no');

    // assert
    test.value(g.graphShowHighRisks).isEqualTo(true);
    test.value(g.graphShowMediumRisks).isEqualTo(true);
    test.value(g.graphShowLowRisks).isEqualTo(true);
    test.value(g.graphShowNoRisks).isEqualTo(true);

    // act
    g.toggleRisks();

    // assert
    test.value(g.graphShowHighRisks).isEqualTo(true);
    test.value(g.graphShowMediumRisks).isEqualTo(true);
    test.value(g.graphShowLowRisks).isEqualTo(true);
    test.value(g.graphShowNoRisks).isEqualTo(true);
  });
  it('graphDropdown should show and hide a togglable dropdown button', async function() {
    // arrange
    const graph = await import('../src/js/graph.js');
    const rc = new RiskCounters();
    const g = new graph.Graph(undefined, rc);

    // act
    g.graphDropdown();

    // assert
    test.value(document.getElementById('myDropdown').classList.contains('show')).isEqualTo(true);

    // act
    g.graphDropdown();

    // assert
    test.value(document.getElementById('myDropdown').classList.contains('show')).isEqualTo(false);
  });
  it('getData should fill the graph with the correct data', async function() {
    // arrange
    const graph = await import('../src/js/graph.js');

    const expectedData = {
      'labels': ['1', '2', '3', '4', '5'],
      'datasets': [{
        'label': 'Acceptable',
        'data': [3, 4, 5, 6, 4],
        'backgroundColor': 'rgb(255, 255, 0)',
      }, {
        'label': 'Low',
        'data': [3, 4, 5, 6, 3],
        'backgroundColor': 'rgb(255, 0, 0)',
      }, {
        'label': 'Medium',
        'data': [3, 4, 5, 6, 2],
        'backgroundColor': 'rgb(0, 0, 255)',
      }, {
        'label': 'High',
        'data': [3, 4, 5, 6, 1],
        'backgroundColor': 'rgb(0, 255, 255)',
      }, {
        'label': 'Info',
        'data': [3, 4, 5, 6, 5],
        'backgroundColor': 'rgb(255, 255, 255)',
      }],
    };

    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',
      infoColor: 'rgb(255, 255, 255)',

      allHighRisks: [1, 2, 3, 4, 5, 6, 1],
      allMediumRisks: [1, 2, 3, 4, 5, 6, 2],
      allLowRisks: [1, 2, 3, 4, 5, 6, 3],
      allNoRisks: [1, 2, 3, 4, 5, 6, 4],
      allInfoRisks: [1, 2, 3, 4, 5, 6, 5],
    };

    const g = new graph.Graph(undefined, mockRiskCounters);

    // act
    const result = await g.getData();

    // assert
    test.array(result.labels).is(expectedData.labels);
    test.array(result.datasets).is(expectedData.datasets);
  });
  it('getOptions should return the correct graph options', async function() {
    // arrange
    const graph = await import('../src/js/graph.js');
    const rc = new RiskCounters();
    const g = new graph.Graph(undefined, rc);

    const expectedOptions = {
      scales: {
        x: {
          stacked: true,
        },
        y: {
          stacked: true,
        },
      },
      plugins: {
        legend: {
          display: false,
        },
      },
      maintainAspectRatio: false,
      categoryPercentage: 1,
    };

    // act
    const resultOptions = g.getOptions();

    // assert
    test.object(resultOptions).is(expectedOptions);
  });
  it('Creating a graph should call getOptions and getData', async function() {
    // arrange
    const graph = await import('../src/js/graph.js');
    const chart = await import('chart.js/auto');
    const rc = new RiskCounters();
    const getDataMock = jest.spyOn(graph.Graph.prototype, 'getData');
    const getOptionsMock = jest.spyOn(graph.Graph.prototype, 'getOptions');

    // act
    const g = new graph.Graph('interval-graph', rc);
    await g.createGraphChart();

    // assert
    expect(getDataMock).toHaveBeenCalled();
    expect(getOptionsMock).toHaveBeenCalled();
    expect(chart.Chart).toHaveBeenCalled();
    getDataMock.mockRestore();
    getOptionsMock.mockRestore();
  });
});
