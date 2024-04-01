import { CloseNavigationHamburger } from "./navigation-menu";

/** Load the content of the Integration page */
function openIntegrationPage() {
    CloseNavigationHamburger();
    document.getElementById("page-contents").innerHTML = `
    <div class="dashboard-data"></div>
    `;
  }
  
  document.getElementById("integration-button").addEventListener("click", () => openIntegrationPage());