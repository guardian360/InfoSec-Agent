import { GetLocalization } from './localize.js';
import { retrieveTheme } from "./personalize";

/** Load the content of the Privacy Dashboard page */
function openPrivacyDashboardPage() {
  document.getElementById("page-contents").innerHTML = `
  <div class="dashboard-data"></div>
  `;
  document.onload = retrieveTheme();
}

document.getElementById("privacy-dashboard-button").addEventListener("click", () => openPrivacyDashboardPage());