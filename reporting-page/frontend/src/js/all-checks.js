import {closeNavigation, markSelectedNavigationItem} from './navigation-menu.js';
import {retrieveTheme} from './personalize.js';
import {getLocalization} from './localize.js';
import {LogError as logError} from '../../wailsjs/go/main/Tray.js';
import {openIssuePage} from './issue.js';

/** Load the content of the About page */
export function openAllChecksPage() {
  console.log('test');
  retrieveTheme();
  closeNavigation(document.body.offsetWidth);
  markSelectedNavigationItem('all-checks-button');
  sessionStorage.setItem('savedPage', '8');

  document.getElementById('page-contents').innerHTML = `
  <div class="all-checks">
    <div class="all-checks-container">
      <div class="all-checks-segment all-checks-title"> <!-- title top segment -->
        <p class="">Security areas<p> 
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Applications</p>
        </div>
        <p class="all-checks-segment-text">
          Applications are computer software made for users with a certain functionality. Most 
          applications are not malicious when installed from the right place, but as a user you 
          have to be aware of what you are downloading and using. Here is the check we run 
          regarding applications: 
        </p>
        <div class="checksList" id="securityApplications"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Devices</p>
        </div>
        <p class="all-checks-segment-text">
          Your computer is a device, but it can also be connected to other devices. These devices are mostly 
          tools like a mouse, keyboard and headset. But because these devices are connected directly to your 
          computer, they have a lot of access to it. To make sure no malicious devices, used to steal data 
          or comprimise the security of you computer, you need to check whether you know the devices 
          connected to your PC. Here are some checks we run regarding devices:
        </p>
        <div class="checksList" id="securityDevices"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Network</p>
        </div>
        <p class="all-checks-segment-text">
          Most of your data will be stored on your computer, but some data will be transferred across 
          networks. Malicious users can intercept your data if it is not transferred in a secure manner. 
          Here are some checks we run regarding your network:
        </p>
        <div class="checksList" id="securityNetwork"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Operating System</p>
        </div>
        <p class="all-checks-segment-text">
          The Operating System is the backbone of computers. It therefore plays an integral part in keeping
          your computer secure. When setting up your OS, it already enables and disables some features to 
          keep your computer secure, but not all settings are preferable. Here are some checks we run 
          regarding your OS:
        </p>
        <div class="checksList" id="securityOS"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Passwords</p>
        </div>
        <p class="all-checks-segment-text">
          Passwords are keys which prevent malicious users from having access to your computer or accounts. 
          Having a password is almost always required when creating an account somewhere, but having a 
          password does not always mean it is secure enough. Easy passwords can be guessed, but hard ones 
          may be hard to remember. Here are some checks we run regarding your passwords:
        </p>
        <div class="checksList" id="securityPasswords"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Other</p>
        </div>
        <p class="all-checks-segment-text">
          Keeping your computer secure can also be done in some other ways. Here are some checks we run 
          to check if your computer is secure:
        </p>
        <div class="checksList" id="securityOther"></div>
      </div>
      <div class="all-checks-segment all-checks-title"> <!-- title bottom segment -->
        <p class="">Privacy areas<p> 
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Permissions</p>
        </div>
        <p class="all-checks-segment-text">
          Applications installed on your computer want to have access to or permission to use certain 
          features or devices connected to your PC. It is your responsibility as a user to give 
          those permissions to certain applications or not. It is therefore recommended to search 
          on the internet if an application you use should or should not have certain permissions. 
          Here are some checks we run regarding permissions: 
        </p>
        <div class="checksList" id="privacyPermissions"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Browser</p>
        </div>
        <p class="all-checks-segment-text">
          To have access to the internet, you need to have a browser. While a browser is really useful, 
          it is not always keeping your data private. To stay safe while browsing on the internet, you 
          have to always keep in mind what you are clicking on. There are some recommendations regarding 
          what you do and have in your brower. Here are some checks we run regarding your browser:
        </p>
        <div class="checksList" id="privacyBrowserChrome"></div>
      </div>
      <div class="all-checks-segment">
        <div class="all-checks-segment-header">
          <p class="">Other</p>
        </div>
        <p class="all-checks-segment-text">
          Keeping your data private can also be done in some other ways. Here are some checks we run 
          to check if your data is private:
        </p>
        <div class="checksList" id="privacyOther"></div>
      </div>
    </div>
  </div>
  `;

  const elements = document.getElementsByClassName('checksList');
  for (let i = 0; i < elements.length; i++) {
    elements[i].innerHTML = createBulletList(elements[i].id);
  }
  const issues = JSON.parse(sessionStorage.getItem('DataBaseData'));
  const checks = document.getElementsByTagName('li');
  for (let i = 0; i < checks.length; i++) {
    const issue = issues.find((issue) => issue.id == checks[i].id);
    checks[i].addEventListener('click', () => openIssuePage(issue.jsonkey,issue.severity));
  }

  // Localize the static content of the about page
  const staticAboutPageConent = [
  ];
  const localizationIds = [
  ];
  for (let i = 0; i < staticAboutPageConent.length; i++) {
    getLocalization(localizationIds[i], staticAboutPageConent[i]);
  }
}

/* istanbul ignore next */
if (typeof document !== 'undefined') {
  try {
    document.getElementById('all-checks-button').addEventListener('click', () => openAllChecksPage());
  } catch (error) {
    logError('Error in all-checks.js: ' + error);
  }
}

const areaLists = {
  'securityApplications': [
    ['Applications that start at system boot',20],
  ],
  'securityDevices': [
    ['Devices connected via bluetooth',1],
    ['Devices connected via USB ports',2],
  ],
  'securityNetwork': [
    ['Open ports used by applications',11],
    ['System Message Block enabled/disabled',13],
  ],
  'securityOS': [
    ['User Account Control enabled/disabled',14],
    ['Windows Defender enabled/disabled',15],
    ['Windows login method',17],
    ['Windows version up to date',18],
    ['Secure boot enabled/disabled',19],
    ['Automatic log in enabled/disabled',33],
    ['Windows Firewall enabled/disabled',37],
  ],
  'securityPasswords': [
    ['Password manager installed',5],
    ['Last Windows password changed',16],
    ['Windows password length',38],
  ],
  'securityOther': [
    ['Guest account enabled/disabled',3],
    ['Remote desktop enabled/disabled',12],
    ['CIS Audit list registry settings',32],
    ['Remote Procedure Call enabled/disabled',34],
    ['Credential Guard enabled/disabled',39],
  ],
  'privacyPermissions': [
    ['Applications with location permissions',6],
    ['Applications with microphone permissions',7],
    ['Applications with webcam permissions',8],
    ['Applications with appointments permissions',9],
    ['Applications with contacts permissions',10],
  ],
  'privacyBrowserChrome': [
    ['Google Chrome has an adblocker installed',21],
    ['Phishing domains found in Google Chrome history',23],
    ['Google Chrome search engine in use',25],
    ['Google Chrome cookies stored',35],
  ],
  'privacyBrowserEdge': [
    ['Microsoft Edge has an adblocker installed',22],
    ['Phishing domains found in Microsoft Edge history',24],
    ['Microsoft Edge search engine in use',26],
    ['Microsoft Edge cookies stored',36],
  ],
  'privacyBrowserFirefox': [
    ['Mozilla Firefox has an adblocker installed',29],
    ['Phishing domains found in Mozilla Firefox history',31],
    ['Mozilla Firefox search engine in use',30],
    ['Mozilla Firefox cookies stored',27],
    ['Mozilla Firefox extensions installed',28],
  ],
  'privacyOther': [
    ['Advertisement ID/Network sharing enabled/disabled',4],
  ],
}

/**
   * create a bullet list for each entry of a security or privacy area
   * @param {string} listID id of the list to create a bullet list for
   */
function createBulletList(listId) {
  const list = areaLists[listId];
  let resultLine = `<ul>`;
  list.forEach((check) => {
    resultLine += `<li id="${check[1]}">${check[0]}</li>`;
  });
  resultLine += `</ul>`;
  return resultLine;
}