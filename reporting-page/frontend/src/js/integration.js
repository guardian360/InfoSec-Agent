import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {getLocalization} from './localize.js';
// import {LoadUserSettings as loadUserSettings} from '../../wailsjs/go/main/App.js';

export let currentStep = 1;
/** Load the content of the Integration page */
export function openIntegrationPage() {
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('integration-button');
  sessionStorage.setItem('savedPage', '7');

  document.getElementById('page-contents').innerHTML = `
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lighthouse API</title>
  </head>
  <body>
  <div class="api-key-container">
    <div class="api-key-form">
      <h1 class="lang-api-key-title"></h1>
      <p class="lang-api-key-description"></p>

      <div class="api-key-step" id="step1">
        <h3 class="lang-step-1-title"></h3>
        <p class="lang-step-1"></p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 1 Image">
        </div>
      </div>

      <div class="api-key-step" id="step2">
        <h3 class="lang-step-2-title"></h3>
        <p class="lang-step-2"></p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 2 Image">
        </div>
      </div>

      <div class="api-key-step" id="step3">
        <h3 class="lang-step-3-title"></h3>
        <p class="lang-step-3"></p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 3 Image">
        </div>
      </div>

      <div id="steps">
      <button class="api-key-button lang-previous" id="prevBtn"></button>
      <button class="api-key-button lang-next" id="nextBtn"></button>
      </div>

      <h2 class="lang-enter-api"></h2>
      <input type="password" class="api-key-input" id="apiKeyInput">
      <button class="api-key-button lang-connect" id="apiKeyButtonClick" disabled></button>
      <button class="api-key-button lang-disconnect" id="disconnectButton"></button>
      <div class="api-key-status" id="status"></div>
      </div>
    </div>
  </body>
  </html>`;

  document.getElementById('nextBtn').addEventListener('click', () => nextStep());
  document.getElementById('prevBtn').addEventListener('click', () => prevStep());
  document.getElementById('apiKeyButtonClick').addEventListener('click', () => connectToAPI());
  document.getElementById('disconnectButton').addEventListener('click', () => disconnectFromAPI());
  // When the integration with the lighthouse API is completed this section will be uncommented
  // and can be used to load the integration details.
  // const userSettings = loadUserSettings();
  // if (userSettings.IntegrationKey !== '') {
  //   document.getElementById('apiKeyButtonClick').style.display = 'none';
  //   document.getElementById('disconnectButton').style.display = 'inline-block';
  //   document.getElementById('apiKeyInput').value = userSettings.IntegrationKey;
  // } else {
  document.getElementById('apiKeyButtonClick').style.display = 'inline-block';
  document.getElementById('disconnectButton').style.display = 'none';
  // }

  showStep(currentStep);

  // Localize the static content of the home page
  const staticHomePageContent = [
    'lang-api-key-description',
    'lang-api-key-title',
    'lang-step-1-title',
    'lang-step-2-title',
    'lang-step-3-title',
    'lang-step-1',
    'lang-step-2',
    'lang-step-3',
    'lang-previous',
    'lang-next',
    'lang-enter-api',
    'lang-connect',
    'lang-disconnect',
    'lang-connected-api',
    'lang-disconnected-api',
    'lang-please-enter-api',
  ];
  const localizationIds = [
    'Integration.ApiKeyDescription',
    'Integration.ApiKeyTitle',
    'Integration.Step1Title',
    'Integration.Step2Title',
    'Integration.Step3Title',
    'Integration.Step1',
    'Integration.Step2',
    'Integration.Step3',
    'Integration.Previous',
    'Integration.Next',
    'Integration.EnterApi',
    'Integration.Connect',
    'Integration.Disconnect',
    'Integration.ConnectedApi',
    'Integration.DisconnectedApi',
    'Integration.PleaseEnterApi',
  ];
  for (let i = 0; i < staticHomePageContent.length; i++) {
    getLocalization(localizationIds[i], staticHomePageContent[i]);
  }
}

/**
 * This function shows the step of the API key connection process
 * @param {number} step The step to show
 */
function showStep(step) {
  const steps = document.querySelectorAll('.api-key-step');
  steps.forEach((s) => s.style.display = 'none');
  document.getElementById('step' + step).style.display = 'block';
  currentStep = step;
}

/**
 * This function navigates to the next step of the API key connection process
 */
function nextStep() {
  if (currentStep < 3) {
    showStep(currentStep + 1);
  }
}

/**
 * This function navigates to the previous step of the API key connection process
 */
function prevStep() {
  if (currentStep > 1) {
    showStep(currentStep - 1);
  }
}

/**
 * This function connects to the API using the entered API key
 */
export function connectToAPI() {
  const apiKey = document.getElementById('apiKeyInput').value;

  const status = document.getElementById('status');
  if (apiKey.trim() === '') {
    status.classList.remove('lang-connected-api');
    status.classList.remove('lang-disconnected-api');
    status.classList.add('lang-please-enter-api');
    status.style.color = 'red';
  } else {
    status.classList.remove('lang-please-enter-api');
    status.classList.remove('lang-disconnected-api');
    status.classList.add('lang-connected-api');
    status.style.color = 'green';
    // Hide API key input after connection
    document.getElementById('apiKeyButtonClick').style.display = 'none';
    document.getElementById('disconnectButton').style.display = 'inline-block';
  }
}

/**
 * This function disconnects from the API
 */
export function disconnectFromAPI() {
  // Dummy disconnect logic
  const status = document.getElementById('status');
  status.classList.remove('lang-please-enter-api');
  status.classList.remove('lang-connected-api');
  status.classList.add('lang-disconnected-api');
  status.style.color = 'red';
  // Show API key input
  document.getElementById('apiKeyButtonClick').style.display = 'inline-block';
  // Hide disconnect button
  document.getElementById('disconnectButton').style.display = 'none';
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('integration-button').addEventListener('click', () => openIntegrationPage());
  } catch (error) {
    logError('Error in integration.js: ' + error);
  }
}
