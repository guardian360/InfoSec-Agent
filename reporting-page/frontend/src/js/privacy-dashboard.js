import { GetLocalization } from './localize.js';
import { CloseNavigation } from "./navigation-menu.js";
import { MarkSelectedNavigationItem } from "./navigation-menu.js";

/** Load the content of the Privacy Dashboard page */
function openPrivacyDashboardPage() {
  CloseNavigation();
  MarkSelectedNavigationItem("privacy-dashboard-button");
  
  document.getElementById("page-contents").innerHTML = `
  <div class="dashboard-data"></div>
  `;
}

document.getElementById("privacy-dashboard-button").addEventListener("click", () => openPrivacyDashboardPage());