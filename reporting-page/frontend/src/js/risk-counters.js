/**
 * Represents the correct count for each risk assessment.
 */
export class RiskCounters {
  highRiskColor;
  mediumRiskColor;
  lowRiskColor;
  noRiskColor;

  allHighRisks = [];
  allMediumRisks = [];
  allLowRisks = [];
  allNoRisks = [];

  lastHighRisk;
  lastMediumRisk;
  lastLowRisk;
  lastNoRisk;

  count = 1;
  /** Create the risk-Counters with the right colors
   *
   * @param {int} high Last recorded amount of high risks
   * @param {int} medium Last recorded amount of medium risks
   * @param {int} low Last recorded amount of low risks
   * @param {int} acceptable Last recorded amount of acceptable risks
   * @param {boolean} [testing=false] Specifies if the class is being used in testing, normally set to *false*
   */
  constructor(high, medium, low, acceptable, testing=false) {
    if (testing) {
      this.highRiskColor = 'rgb(0, 255, 255)';
      this.mediumRiskColor = 'rgb(0, 0, 255)';
      this.lowRiskColor = 'rgb(255, 0, 0)';
      this.noRiskColor = 'rgb(255, 255, 0)';
    } else {
      this.highRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--high-risk-color');
      this.mediumRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--medium-risk-color');
      this.lowRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--low-risk-color');
      this.noRiskColor = getComputedStyle(document.documentElement).getPropertyValue('--no-risk-color');
    }
    this.allHighRisks.push(high);
    this.allMediumRisks.push(medium);
    this.allLowRisks.push(low);
    this.allNoRisks.push(acceptable);

    this.lastHighRisk = high;
    this.lastMediumRisk = medium;
    this.lastLowRisk = low;
    this.lastNoRisk = acceptable;
  }
}

/**
 * Updates the RiskCounters instance with new risk assessments and recalculates the maximum count.
 *
 * @param {RiskCounters} rc - The RiskCounters instance to be updated.
 * @param {number} high - The last recorded amount of high risks.
 * @param {number} medium - The last recorded amount of medium risks.
 * @param {number} low - The last recorded amount of low risks.
 * @param {number} acceptable - The last recorded amount of acceptable risks.
 * @return {RiskCounters} The updated RiskCounters instance.
 */
export function updateRiskCounter(rc, high, medium, low, acceptable) {
  rc.allHighRisks.push(high);
  rc.allMediumRisks.push(medium);
  rc.allLowRisks.push(low);
  rc.allNoRisks.push(acceptable);

  rc.lastHighRisk = high;
  rc.lastMediumRisk = medium;
  rc.lastLowRisk = low;
  rc.lastNoRisk = acceptable;
  rc.count = calculateMaxCount(rc);
  return rc;
}

/**
 * Calculates the maximum length among all arrays of a RiskCounters instance.
 *
 * @param {RiskCounters} rc - The RiskCounters instance containing arrays to calculate maximum length from.
 * @return {number} The maximum length among all arrays.
 */
function calculateMaxCount(rc) {
  return Math.max(
    rc.allHighRisks.length,
    rc.allMediumRisks.length,
    rc.allLowRisks.length,
    rc.allNoRisks.length,
  );
}
