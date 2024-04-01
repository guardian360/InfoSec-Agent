import { CloseNavigation } from "./navigation-menu";
import { MarkSelectedNavigationItem } from "./navigation-menu";
import { retrieveTheme } from "./personalize";

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