import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {getLocalization} from './localize.js';
//import {LoadUserSettings as loadUserSettings} from '../../wailsjs/go/main/App.js';

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
      <h1 class="lang-api-key-title">Connect to the Lighthouse API</h1>
      <p class="lang-api-key-description"><b>Coming soon:</b> this functionality is not available yet, you will be able to connect to the Lighthouse API in the near future.</p>

      <div class="api-key-step" id="step1">
        <h3>Step 1: Register for an account</h3>
        <p class="lang-step-1">Login to your account on https://lighthouse.guardian360.nl and go to the settings</p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 1 Image">
        </div>
      </div>

      <div class="api-key-step lang-step-2" id="step2">
        <h3>Step 2: Navigate to API Settings</h3>
        <p class="lang-step-2">Go to the API section and create a new token</p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 2 Image">
        </div>
      </div>

      <div class="api-key-step lang-step-3" id="step3">
        <h3>Step 3: Generate API Key</h3>
        <p class="lang-step-3">Copy the token and paste it in the field below, you can then click on the 'connect' button</p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 3 Image">
        </div>
      </div>

      <div id="steps">
      <button class="api-key-button" id="prevBtn">&#8249; Previous</button>
      <button class="api-key-button" id="nextBtn">Next &#8250;</button>
      </div>

      <h2>Enter Your API Key:</h2>
      <input type="password" class="api-key-input" id="apiKeyInput">
      <button class="api-key-button" id="apiKeyButtonClick" disabled>Connect</button>
      <button class="api-key-button" id="disconnectButton">Disconnect</button>
      <div class="api-key-status" id="status"></div>
      </div>
    </div>
  </body>
  </html>`;

  document.getElementById('nextBtn').addEventListener('click', () => nextStep());
  document.getElementById('prevBtn').addEventListener('click', () => prevStep());
  document.getElementById('apiKeyButtonClick').addEventListener('click', () => connectToAPI());
  document.getElementById('disconnectButton').addEventListener('click', () => disconnectFromAPI());
  // When the integration with the lighthouse API is completed this section will be uncommented and can be used to load the integration details.
  // const userSettings = loadUserSettings();
  // if (userSettings.IntegrationKey !== '') {
  //   document.getElementById('apiKeyButtonClick').style.display = 'none';
  //   document.getElementById('disconnectButton').style.display = 'inline-block';
  //   document.getElementById('apiKeyInput').value = userSettings.IntegrationKey;
  // } else {
  document.getElementById('apiKeyButtonClick').style.display = 'inline-block';
  document.getElementById('disconnectButton').style.display = 'none';
  //}

  showStep(currentStep);

  // Localize the static content of the home page
  const staticHomePageContent = [
    'lang-api-key-description',
    'lang-api-key-title',
    'lang-step-1',
    'lang-step-2',
    'lang-step-3',
  ];
  const localizationIds = [
    'Integration.apiKeyDescription',
    'Integration.apiKeyTitle',
    'Integration.step1',
    'Integration.step2',
    'Integration.step3',
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
function connectToAPI() {
  const apiKey = document.getElementById('apiKeyInput').value;
  setTimeout(() => {
    const status = document.getElementById('status');
    if (apiKey.trim() === '') {
      status.innerHTML = 'Please enter your API key.';
      status.style.color = 'red';
    } else {
      status.innerHTML = 'Connected to API.';
      status.style.color = 'green';
      // Hide API key input after connection
      document.getElementById('apiKeyButtonClick').style.display = 'none';
      document.getElementById('disconnectButton').style.display = 'inline-block';
    }
  }, 1000);
}

/**
 * This function disconnects from the API
 */
function disconnectFromAPI() {
  // Dummy disconnect logic
  const status = document.getElementById('status');
  status.innerHTML = 'Disconnected from API.';
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
