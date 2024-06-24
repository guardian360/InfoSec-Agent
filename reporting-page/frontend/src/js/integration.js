import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';

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
      <h1 class="api-key-title">Connect to the Lighthouse API</h1>

      <div class="api-key-step" id="step1">
        <h3>Step 1: Register for an account</h3>
        <p>Visit our website and register for a new account.</p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 1 Image">
        </div>
      </div>

      <div class="api-key-step" id="step2">
        <h3>Step 2: Navigate to API Settings</h3>
        <p>Once logged in, navigate to your account settings and find the API section.</p>
        <div class="api-key-image-container">
          <img class="api-key-image" src="https://via.placeholder.com/400" alt="Step 2 Image">
        </div>
      </div>

      <div class="api-key-step" id="step3">
        <h3>Step 3: Generate API Key</h3>
        <p>Generate a new API key and copy it.</p>
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
      <button class="api-key-button" id="apiKeyButtonClick">Connect</button>
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
  document.getElementById('disconnectButton').style.display = 'none';

  showStep(currentStep);
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
    status.innerHTML = 'Please enter your API key.';
    status.style.color = 'red';
  } else {
    status.innerHTML = 'Connected to API.';
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
