import {closeNavigation} from './navigation-menu.js';
// import {loadPersonalizeNavigation} from './navigation-menu.js';
// import logoPhoto from '../assets/images/logoTeamA-transformed.png';

/** Load the content of the Personalize page */
export function openPersonalizePage() {
  closeNavigation();

  document.getElementById('page-contents').innerHTML = `
  <div class="setting">
    <span class="setting-description favicon-title ">Favicon</span>
    <div class="personalize-button-container">
      <button class="setting-button icon-button" type="button">Change icon </button>    
      <input class="personalize-input-invisible" type="file" id="input-file-icon" accept=".ico, .png">
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">Navigation image</span>
    <div class="personalize-button-container">
      <button class="setting-button logo-button" type="button">Change logo </button>    
      <input class="personalize-input-invisible"
      type="file" 
      id="input-file-picture" 
      accept="image/jpeg, image/png, image/jpg">
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">Navigation title</span>
    <div class="personalize-button-container">
      <button class="setting-button title-button" type="button">Change title </button>
      <div id="custom-modal" class="modal">
        <div class="modal-content">
          <input type="text" id="new-title-input">
          <button id="saveTitleButton">Save</button>
        </div>
      </div>
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <form action="" class="color-picker>
      <fieldset>
        <legend>Pick a theme</legend>
        <label for="normal">Light</label>
        <input type="radio" name="theme" id="normal" checked>
        <label for="dark">Dark</label>
        <input type="radio" name="theme" id="dark">
    <div class="personalize-button-container">
      <button class="setting-button reset-button" type="button">Reset settings </button>    
    </div>
  </div>
  `;
  // add eventlistener for changing Favicon
  const changeIconButton = document.getElementsByClassName('icon-button')[0];
  const inputFileIcon = document.getElementById('input-file-icon');

  changeIconButton.addEventListener('click', function() {
    inputFileIcon.click();
  });

  inputFileIcon.addEventListener('change', handleFaviconChange);

  // add eventlistener for changing navication picture
  const changeLogoButton = document.getElementsByClassName('logo-button')[0];
  const inputLogo = document.getElementById('input-file-picture');

  changeLogoButton.addEventListener('click', function() {
    inputLogo.click();
  });

  inputLogo.addEventListener('change', handlePictureChange);

  // add eventlistener for changing navigation title
  const changeTitleButton = document.getElementsByClassName('title-button')[0];
  const customModal = document.getElementById('custom-modal');
  const newTitleInput = document.getElementById('new-title-input');
  const saveTitleButton = document.getElementById('saveTitleButton');

  changeTitleButton.addEventListener('click', function() {
    customModal.style.display = 'block';
    newTitleInput.focus();
  });

  saveTitleButton.addEventListener('click', function() {
    const newTitle = newTitleInput.value.trim();
    if (newTitle !== '') {
      handleTitleChange(newTitle);
      customModal.style.display = 'none';
    }
  });


  // add eventlistener for changing Favicon
  const changeResetButton = document.getElementsByClassName('reset-button')[0];

  changeResetButton.addEventListener('click', function() {
    resetSettings();
  });

  /*
  // add eventlistener for changing navigation title
  const inputBackgroundNav = document.getElementById('input-color-background');
  inputBackgroundNav.addEventListener('change', handleLeftBackgroundNav);
  */

  /* save themes*/
  const themes = document.querySelectorAll('[name="theme"]');
  themes.forEach((themeOption) => {
    themeOption.addEventListener('click', () => {
      localStorage.setItem('theme', themeOption.id);
      // loadPersonalizeNavigation();
    });
  });

  const activeTheme = localStorage.getItem('theme');
  themes.forEach((themeOption) => {
    if (themeOption.id === activeTheme) {
      themeOption.checked = true;
    }
  });
  document.documentElement.className= activeTheme;
}


/**
 * Handles the change event when selecting a new favicon file.
 * Updates the favicon of the document with the selected image.
 * Saves the selected image URL in the localStorage.
 * @param {Event} icon - The event object representing the change of the favicon input.
 */
export function handleFaviconChange(icon) {
  const file = icon.target.files[0]; // Get the selected file
  if (file) {
    const reader = new FileReader();
    reader.onload = function(e) {
      const picture = e.target.result;
      const favicon = document.querySelector('link[rel="icon"]');
      if (favicon) {
        favicon.href = picture;
      } else {
        const newFavicon = document.createElement('link');
        newFavicon.rel = 'icon';
        newFavicon.href = picture;
        document.head.appendChild(newFavicon);
      }
      localStorage.setItem('favicon', picture);
    };
    reader.readAsDataURL(file); // Read the selected file as a Data URL
  }
}

/**
 * Handles the change event when selecting a new picture file.
 * Updates the source of the specified image element with the selected image.
 * Saves the selected image URL in the localStorage.
 * @param {Event} picture - The event object representing the change of the picture input.
 */
export function handlePictureChange(picture) {
  const file = picture.target.files[0]; // Get the selected file
  const reader = new FileReader();
  reader.onload = function(e) {
    const logo = document.getElementById('logo');
    logo.src = e.target.result; // Set the source of the logo to the selected image
    localStorage.setItem('picture', e.target.result);
  };
  reader.readAsDataURL(file); // Read the selected file as a Data URL
}

/**
 * Handles the change event when updating the title.
 * Updates the text content of the specified title element with the new title value.
 * Saves the new title value in the localStorage.
 * @param {string} newTitle - The new title to set.
 */
export function handleTitleChange(newTitle) {
  const titleElement = document.getElementById('title');
  titleElement.textContent = newTitle; // Set the text content to the new title
  localStorage.setItem('title', newTitle);
}

/**
 * Handles the change event when updating the background color of the left navigation.
 * Retrieves the selected color from the color picker input.
 * Updates the background color of the left navigation with the selected color.
 */
export function handleLeftBackgroundNav() {
  const colorPicker = document.getElementById('input-color-background');
  const color = colorPicker.value;
  const temp = document.getElementsByClassName('left-nav')[0];
  temp.style.backgroundColor = color;
}
/**
 * Retrieves the active theme from localStorage and applies it to the document's root element.
 * The active theme class name is retrieved from the 'theme' key in localStorage.
 */
export function retrieveTheme() {
  const activeTheme = localStorage.getItem('theme');
  document.documentElement.className = activeTheme;
}
/**
 * Resets the settings by clearing localStorage and restoring default values.
 */
export function resetSettings() {
  localStorage.clear();
  logoPhoto = "frontend/src/assets/images/logoTeamA-transformed.png";
  const logo = document.getElementById('logo');
  logo.src = logoPhoto;

  const title = document.getElementById('title');
  title.textContent = 'Little Brother';

  const favicon = document.getElementById('favicon');
  favicon.href = logoPhoto;
}
