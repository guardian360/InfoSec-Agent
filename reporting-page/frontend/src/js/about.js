import { CloseNavigation } from "./navigation-menu";
import { MarkSelectedNavigationItem } from "./navigation-menu";

/** Load the content of the About page */
function openAboutPage() {
  CloseNavigation();
  MarkSelectedNavigationItem("about-button");
  
  document.getElementById("page-contents").innerHTML = `
  <div class="dashboard-data"></div>
  `;
}
  
document.getElementById("about-button").addEventListener("click", () => openAboutPage());