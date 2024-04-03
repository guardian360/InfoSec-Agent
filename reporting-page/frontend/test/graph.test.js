import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {Graph} from '../src/js/graph.js';
import {RiskCounters} from '../src/js/risk-counters.js';

// Mock page
const dom = new JSDOM(`
  <div class="graph-row">
  <div class="graph-column issues-graph-buttons">
    <H2>In this graph you are able to see the distribution of different issues we have found over the past 5 times we ran a check.</H2>
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

// test cases
describe('risk-graph', function() {
  // arrange
  const rc = new RiskCounters(true);
  let g = new Graph(undefined, rc);
  it('Should have togglable risks', function() {
    // act
    g.ToggleRisks('high', false);
    g.ToggleRisks('medium', false);
    g.ToggleRisks('low', false);
    g.ToggleRisks('no', false);

    // assert
    test.value(g.graphShowHighRisks).isEqualTo(false);
    test.value(g.graphShowMediumRisks).isEqualTo(false);
    test.value(g.graphShowLowRisks).isEqualTo(false);
    test.value(g.graphShowNoRisks).isEqualTo(false);

    // act
    g.ToggleRisks('high', false);
    g.ToggleRisks('medium', false);
    g.ToggleRisks('low', false);
    g.ToggleRisks('no', false);

    // assert
    test.value(g.graphShowHighRisks).isEqualTo(true);
    test.value(g.graphShowMediumRisks).isEqualTo(true);
    test.value(g.graphShowLowRisks).isEqualTo(true);
    test.value(g.graphShowNoRisks).isEqualTo(true);
  });
  it('Should have a togglable dropdown button', function() {
    // act
    g.GraphDropdown();

    // assert
    test.value(document.getElementById('myDropdown').classList.contains('show')).isEqualTo(true);

    // act
    g.GraphDropdown();

    // assert
    test.value(document.getElementById('myDropdown').classList.contains('show')).isEqualTo(false);
  });
  it('Should contain the right data', function() {
    // arrange
    const expectedData = {
      'labels': [1, 2, 3, 4, 5],
      'datasets': [{
        'label': 'Safe issues',
        'data': [3, 4, 5, 6, 4],
        'backgroundColor': 'rgb(255, 255, 0)',
      }, {
        'label': 'Low risk issues',
        'data': [3, 4, 5, 6, 3],
        'backgroundColor': 'rgb(255, 0, 0)',
      }, {
        'label': 'Medium risk issues',
        'data': [3, 4, 5, 6, 2],
        'backgroundColor': 'rgb(0, 0, 255)',
      }, {
        'label': 'High risk issues',
        'data': [3, 4, 5, 6, 1],
        'backgroundColor': 'rgb(0, 255, 255)',
      }],
    };

    const mockRiskCounters = {
      highRiskColor: 'rgb(0, 255, 255)',
      mediumRiskColor: 'rgb(0, 0, 255)',
      lowRiskColor: 'rgb(255, 0, 0)',
      noRiskColor: 'rgb(255, 255, 0)',

      allHighRisks: [1, 2, 3, 4, 5, 6, 1],
      allMediumRisks: [1, 2, 3, 4, 5, 6, 2],
      allLowRisks: [1, 2, 3, 4, 5, 6, 3],
      allNoRisks: [1, 2, 3, 4, 5, 6, 4],
    };

    g = new Graph(undefined, mockRiskCounters);

    // act
    const resultData = g.GetData();

    // assert
    test.array(resultData.labels).is(expectedData.labels);
    test.array(resultData.datasets).is(expectedData.datasets);
  });
  it('Should use the right options', function() {
    // arrange
    const expectedOptions = {
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

    // act
    const resultOptions = g.GetOptions();

    // assert
    test.object(resultOptions).is(expectedOptions);
  });
});
