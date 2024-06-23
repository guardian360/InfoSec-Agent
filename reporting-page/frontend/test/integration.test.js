import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions, mockGetLocalization, storageMock, clickEvent} from './mock.js';

global.TESTING = true;

// Mock issue page
const dom = new JSDOM(`
<!DOCTYPE html>
<html>
<body>
    <div id="page-contents"></div>
</body>
</html>
`);
global.document = dom.window.document;
global.window = dom.window;

// Mock often used page functions
mockPageFunctions();

// Mock Localize function
jest.unstable_mockModule('../wailsjs/go/main/App.js', () => ({
  Localize: jest.fn().mockImplementation((input) => mockGetLocalization(input)),
}));

// Mock sessionStorage
global.sessionStorage = storageMock;

describe('Integration page', function() {
  it('openIntegrationPage calls show step and hides connect button', async function() {
    // Arrange
    const integration = await import('../src/js/integration.js');

    // Act
    integration.openIntegrationPage();
    const steps = document.querySelectorAll('.api-key-step');
    const button = document.getElementById('disconnectButton');

    // Assert
    test.value(button.style.display).isEqualTo('none');
    test.value(steps[0].style.display).isEqualTo('block');
    test.value(steps[1].style.display).isEqualTo('none');
    test.value(steps[2].style.display).isEqualTo('none');
  });
  it('prevStep should not do anything when on the first step', async function() {
    // Arrange
    await import('../src/js/integration.js');

    // Act
    // calls nextStep
    document.getElementById('prevBtn').dispatchEvent(clickEvent);
    const steps = document.querySelectorAll('.api-key-step');

    // Assert
    test.value(steps[0].style.display).isEqualTo('block');
    test.value(steps[1].style.display).isEqualTo('none');
    test.value(steps[2].style.display).isEqualTo('none');
  });
  it('nextStep should show the next step on the page', async function() {
    // Arrange
    await import('../src/js/integration.js');

    // Act
    // calls nextStep
    document.getElementById('nextBtn').dispatchEvent(clickEvent);
    const steps = document.querySelectorAll('.api-key-step');

    // Assert
    test.value(steps[0].style.display).isEqualTo('none');
    test.value(steps[1].style.display).isEqualTo('block');
    test.value(steps[2].style.display).isEqualTo('none');
  });
  it('prevStep should show the prev step on the page', async function() {
    // Arrange
    await import('../src/js/integration.js');

    // Act
    // calls nextStep
    document.getElementById('prevBtn').dispatchEvent(clickEvent);
    const steps = document.querySelectorAll('.api-key-step');

    // Assert
    test.value(steps[0].style.display).isEqualTo('block');
    test.value(steps[1].style.display).isEqualTo('none');
    test.value(steps[2].style.display).isEqualTo('none');
  });
  it('nextStep does not do anything when the final step is shown', async function() {
    // Arrange
    await import('../src/js/integration.js');

    // Act
    // calls nextStep 1 more time after being on the last step
    const steps = document.querySelectorAll('.api-key-step');
    steps.forEach(() => {
      document.getElementById('nextBtn').dispatchEvent(clickEvent);
    });

    // Assert
    test.value(steps[0].style.display).isEqualTo('none');
    test.value(steps[1].style.display).isEqualTo('none');
    test.value(steps[2].style.display).isEqualTo('block');
  });
  it('connectToAPI asks for input or shows that you are connected', async function() {
    // Arrange
    const integration = await import('../src/js/integration.js');

    // Act
    // no API key given
    integration.openIntegrationPage();
    integration.connectToAPI();
    let status = document.getElementById('status');
    let keyButton = document.getElementById('apiKeyButtonClick');

    // Assert
    test.value(status.innerHTML).isEqualTo('Please enter your API key.');

    // Act
    // API key given
    document.getElementById('apiKeyInput').value = 'abcd';
    keyButton.dispatchEvent(clickEvent);
    status = document.getElementById('status');
    keyButton = document.getElementById('apiKeyButtonClick');
    const disconnectButton = document.getElementById('disconnectButton');

    // Assert
    test.value(status.innerHTML).isEqualTo('Connected to API.');
    test.value(keyButton.style.display).isEqualTo('none');
    test.value(disconnectButton.style.display).isEqualTo('inline-block');
  });
  it('disconnectFromAPI show that you are disconnected after first being connected', async function() {
    // Arrange
    const integration = await import('../src/js/integration.js');

    // Act
    integration.openIntegrationPage();
    const disconnectButton = document.getElementById('disconnectButton');
    const keyButton = document.getElementById('apiKeyButtonClick');
    disconnectButton.dispatchEvent(clickEvent);
    const status = document.getElementById('status');

    // Assert
    test.value(status.innerHTML).isEqualTo('Disconnected from API.');
    test.value(disconnectButton.style.display).isEqualTo('none');
    test.value(keyButton.style.display).isEqualTo('inline-block');
  });
});
