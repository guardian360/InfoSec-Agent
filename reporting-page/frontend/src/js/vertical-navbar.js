import logo from '../assets/images/logoTeamA-transformed.png';
import { Localize } from '../../wailsjs/go/main/App';
import { GetLocalization } from './localize.js';

document.getElementById('logo').src = logo;

let navbarItems = ["settings", "home", "security_dashboard", "issues", "integration", "about"]
let localizationIds = ["Navigation.Settings", "Navigation.Home", "Navigation.SecurityDashboard", "Navigation.Issues", "Navigation.Integration", "Navigation.About"]
for (let i = 0; i < navbarItems.length; i++) {
    GetLocalization(localizationIds[i], navbarItems[i])
}
