import { CloseNavigation, MarkSelectedNavigationItem } from "./navigation-menu.js";
import { retrieveTheme } from "./personalize.js";

/** Load the content of the Integration page */
function openIntegrationPage() {
    CloseNavigation();
    MarkSelectedNavigationItem("integration-button");
    
    document.getElementById("page-contents").innerHTML = `
    <div class="dashboard-data"></div>
    `;

    document.onload = retrieveTheme();
  }
  
  document.getElementById("integration-button").addEventListener("click", () => openIntegrationPage());