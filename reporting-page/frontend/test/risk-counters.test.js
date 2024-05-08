import test from 'unit.js';
import {RiskCounters,updateRiskCounter} from '../src/js/risk-counters.js';

global.TESTING = true;

describe('risk-counters class', function() {
    it('Calling the constructor of risk-counters should fill in the right properties', function() {
        // Arrange
        const expectedHigh = 1;
        const expectedMedium = 2;
        const expectedLow = 3;
        const expectedInfo = 4;
        const expectedNoRisk = 5;
        const rc = new RiskCounters(expectedHigh, expectedMedium, expectedLow, expectedInfo, expectedNoRisk);

        // Act
        const high = rc.allHighRisks[0];
        const medium = rc.allMediumRisks[0];
        const low = rc.allLowRisks[0];
        const info = rc.allInfoRisks[0];
        const noRisk = rc.allNoRisks[0];

        // Arrange
        test.value(high).isEqualTo(expectedHigh);
        test.value(medium).isEqualTo(expectedMedium);
        test.value(low).isEqualTo(expectedLow);
        test.value(info).isEqualTo(expectedInfo);
        test.value(noRisk).isEqualTo(expectedNoRisk);
    })
    it('updateRiskCounter should add new counts to existing counts', function() {
        // Arrange
        const expectedHigh = 1;
        const expectedMedium = 2;
        const expectedLow = 3;
        const expectedInfo = 4;
        const expectedNoRisk = 5;
        const rc = new RiskCounters(expectedHigh, expectedMedium, expectedLow, expectedInfo, expectedNoRisk);
        const newExpectedHigh = 6;
        const newExpectedMedium = 7;
        const newExpectedLow = 8;
        const newExpectedInfo = 9;
        const newExpectedNoRisk = 10;

        // Act
        updateRiskCounter(rc, newExpectedHigh, newExpectedMedium, newExpectedLow, newExpectedInfo, newExpectedNoRisk);

        // Old values
        const high = rc.allHighRisks[0];
        const medium = rc.allMediumRisks[0];
        const low = rc.allLowRisks[0];
        const info = rc.allInfoRisks[0];
        const noRisk = rc.allNoRisks[0];

        // New values
        const newHigh = rc.allHighRisks.pop();
        const newMedium = rc.allMediumRisks.pop();
        const newLow = rc.allLowRisks.pop();
        const newInfo = rc.allInfoRisks.pop();
        const newNoRisk = rc.allNoRisks.pop();

        // Arrange
        // Old values are still there
        test.value(high).isEqualTo(expectedHigh);
        test.value(medium).isEqualTo(expectedMedium);
        test.value(low).isEqualTo(expectedLow);
        test.value(info).isEqualTo(expectedInfo);
        test.value(noRisk).isEqualTo(expectedNoRisk);
        // New values are added
        test.value(newHigh).isEqualTo(newExpectedHigh);
        test.value(newMedium).isEqualTo(newExpectedMedium);
        test.value(newLow).isEqualTo(newExpectedLow);
        test.value(newInfo).isEqualTo(newExpectedInfo);
        test.value(newNoRisk).isEqualTo(newExpectedNoRisk);        
    })
})
