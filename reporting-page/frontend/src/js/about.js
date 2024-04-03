import { CloseNavigation, MarkSelectedNavigationItem } from "./navigation-menu.js";
import { retrieveTheme } from "./personalize.js";

/** Load the content of the About page */
function openAboutPage() {
  CloseNavigation();
  MarkSelectedNavigationItem("about-button");
  
  document.getElementById("page-contents").innerHTML = `
  <div class="dashboard-data"></div>
  `;
    document.onload = retrieveTheme();
  }
  
document.getElementById("about-button").addEventListener("click", () => openAboutPage());