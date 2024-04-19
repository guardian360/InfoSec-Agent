import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';

/** Load the content of the Integration page */
function openIntegrationPage() {
  closeNavigation();
  markSelectedNavigationItem('integration-button');

  document.getElementById('page-contents').innerHTML = `
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lighthouse API</title>
  </head>
  <body>
  <div class="apiKeyContainer">
    <div class="apiKeyForm">
      <h1 class="apiKeyTitle">API Key Connection</h1>

      <h2>Step-by-Step Guide to Obtain API Key:</h2>
      <div class="apiKeyStep" id="step1">
        <h3>Step 1: Register for an account</h3>
        <p>Visit our website and register for a new account.</p>
        <div class="apiKeyImageContainer">
          <img class="apiKeyImage" src="https://via.placeholder.com/400" alt="Step 1 Image">
        </div>
      </div>

      <div class="apiKeyStep" id="step2">
        <h3>Step 2: Navigate to API Settings</h3>
        <p>Once logged in, navigate to your account settings and find the API section.</p>
        <div class="apiKeyImageContainer">
          <img class="apiKeyImage" src="https://via.placeholder.com/400" alt="Step 2 Image">
        </div>
      </div>

      <div class="apiKeyStep" id="step3">
        <h3>Step 3: Generate API Key</h3>
        <p>Generate a new API key and copy it.</p>
        <div class="apiKeyImageContainer">
          <img class="apiKeyImage" src="https://via.placeholder.com/400" alt="Step 3 Image">
        </div>
      </div>

      <div id="steps">
      <button class="apiKeyButton" id="prevBtn">&#8249; Previous</button>
      <button class="apiKeyButton" id="nextBtn">Next &#8250;</button>
      </div>

      <h2>Enter Your API Key:</h2>
      <input type="password" class="apiKeyInput" id="apiKeyInput">
      <button class="apiKeyButton" id="apiKeyButtonClick">Connect</button>
      <button class="apiKeyButton" id="disconnectButton">Disconnect</button>
      <div class="apiKeyStatus" id="status"></div>
      </div>
  </div>
</body>
</html>`;

document.getElementById('nextBtn').addEventListener('click', () => nextStep());
document.getElementById('prevBtn').addEventListener('click', () => prevStep());
document.getElementById('apiKeyButtonClick').addEventListener('click', () => connectToAPI());
document.getElementById('disconnectButton').addEventListener('click', () => disconnectFromAPI());
document.getElementById('disconnectButton').style.display = 'none';

let currentStep = 1;
showStep(currentStep);

function showStep(step) {
  const steps = document.querySelectorAll('.apiKeyStep');
  steps.forEach(s => s.style.display = 'none');
  document.getElementById('step' + step).style.display = 'block';
  currentStep = step;
}

function nextStep() {
  if (currentStep < 3) {
    showStep(currentStep + 1);
  }
}

function prevStep() {
  if (currentStep > 1) {
    showStep(currentStep - 1);
  }
}

function connectToAPI() {
  const apiKey = document.getElementById('apiKeyInput').value;
  // Dummy API connection logic
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

function disconnectFromAPI() {
  // Dummy disconnect logic
  const status = document.getElementById('status');
  status.innerHTML = 'Disconnected from API.';
  status.style.color = 'red';
  // Show API key input
  document.getElementById('apiKeyButtonClick').style.display = 'block';
  // Hide disconnect button
  document.getElementById('disconnectButton').style.display = 'none';
}
document.onload = retrieveTheme();
}

document.getElementById('integration-button').addEventListener('click', () => openIntegrationPage());
