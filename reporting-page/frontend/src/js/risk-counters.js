/**
 * Represents the correct count for each risk assessment.
 */
export class RiskCounters {
  highRiskColor;
  mediumRiskColor;
  lowRiskColor;
  noRiskColor;

  allHighRisks = [1, 1, 1];
  allMediumRisks = [1, 1, 1];
  allLowRisks = [1, 1, 1];
  allNoRisks = [1, 1, 1];

  lastHighRisk;
  lastMediumRisk;
  lastLowRisk;
  lastnoRisk;

  count = this.allHighRisks.length;
  /** Create the risk-Counters with the right colors
   *
   * @param {int} high Last recorded amount of high risks
   * @param {int} medium Last recorded amount of medium risks
   * @param {int} low Last recorded amount of low risks
   * @param {int} safe Last recorded amount of safe risks
   * @param {boolean} [testing=false] Specifies if the class is being used in testing, normally set to *false*
   */
  constructor(high, medium, low, safe, testing=false) {
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
    this.allNoRisks.push(safe);

    this.lastHighRisk = high;
    this.lastMediumRisk = medium;
    this.lastLowRisk = low;
    this.lastnoRisk = safe;
  }

  updateRiskcounter(high, medium, low, safe) {
    this.allHighRisks.push(high);
    this.allMediumRisks.push(medium);
    this.allLowRisks.push(low);
    this.allNoRisks.push(safe);

    this.lastHighRisk = high;
    this.lastMediumRisk = medium;
    this.lastLowRisk = low;
    this.lastnoRisk = safe;
    this.count = this.allHighRisks.length;
  }
}

// sessionStorage.setItem("RiskCounters",JSON.stringify(new RiskCounters()));

