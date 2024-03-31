import { retrieveTheme } from "./personalize";

/** Load the content of the Integration page */
function openIntegrationPage() {
    document.getElementById("page-contents").innerHTML = `
    <div class="dashboard-data"></div>
    `;

    document.onload = retrieveTheme();
  }
  
  document.getElementById("integration-button").addEventListener("click", () => openIntegrationPage());