import 'jsdom-global/register.js';
import test from 'unit.js';
import {JSDOM} from 'jsdom';
import {jest} from '@jest/globals';
import {mockPageFunctions,mockGetLocalization,storageMock,clickEvent} from './mock.js';

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
    it('openIntegrationPage calls show step', async function() {
        // Arrange
        const integration = await import('../src/js/integration.js');
        
        // Act
        integration.openIntegrationPage();
        const steps = document.querySelectorAll('.api-key-step');

        // Assert
        steps[0].style.display = 'block';
        steps[1].style.display = 'none';
        steps[2].style.display = 'none';
    });
    it('nextStep should show the next step on the page', async function() {
        // Arrange
        const integration = await import('../src/js/integration.js');

        // Act
        // calls nextStep
        document.getElementById('nextBtn').dispatchEvent(clickEvent);
        const steps = document.querySelectorAll('.api-key-step');

        // Assert
        steps[0].style.display = 'none';
        steps[1].style.display = 'block';
        steps[2].style.display = 'none';
    });
    it('nextStep should show the next step on the page', async function() {
        // Arrange
        const integration = await import('../src/js/integration.js');

        // Act
        // calls nextStep
        document.getElementById('prevBtn').dispatchEvent(clickEvent);
        const steps = document.querySelectorAll('.api-key-step');

        // Assert
        steps[0].style.display = 'block';
        steps[1].style.display = 'none';
        steps[2].style.display = 'none';
    });
});
