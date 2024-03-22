/** Load the content of the Personalize page */
export function openPersonalizePage() {
  document.getElementById("page-contents").innerHTML = `
  <h1>Customize Page</h1>
  <div class="setting">
    <span class="setting-description">Favicon</span>
    <div class="favicon-button-container">
      <label for="input-file-icon"></label>
      <input type="file" id="input-file-icon" accept=".ico, .png">
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Picture(Top left)</span>
    <div class="picture-button-container">
      <label for="input-file-picture"></label>
      <input type="file" id="input-file-picture" accept="image/jpeg, image/png, image/jpg">   
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Name(Top left)</span>
      <label for="newTitle">Enter new title:</label>
      <input type="text" id="newTitle">
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Font</span>
    <label class="switch">
      <input type="checkbox">
      <span class="slider round"></span>
    </label>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">Background color</span>
    <label class="switch">
      <input type="checkbox">
      <span class="slider round"></span>
    </label>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="setting-description">text color</span>
    <label class="switch">
      <input type="checkbox">
      <span class="slider round"></span>
    </label>
  </div>
  `;
  const faviconInput = document.getElementById('faviconInput');//add eventlistener for changing Favicon
  faviconInput.addEventListener('change', handleFaviconSelect);
  
  const pictureInput = document.getElementById('input-file-picture'); //add eventlistener for changing navication picture
  pictureInput.addEventListener('change', changePicture);
  
  const newTitleInput = document.getElementById('newTitle'); //add eventlistener for changing navigation title
  newTitleInput.addEventListener('input', handleTitleChange);
  }

export function changePicture(picture) {
  const file = picture.target.files[0]; // Get the selected file
  const reader = new FileReader();
  reader.onload = function(e) {
    const logo = document.getElementById('logo'); // Get the logo element
    logo.src = e.target.result; // Set the source of the logo to the selected image
    localStorage.setItem("picture", e.target.result)
    };
  reader.readAsDataURL(file); // Read the selected file as a Data URL
}

/** Select af file to set as the icon
 * 
 * @param {Event} event File event from which file is taken
 */
function handleFaviconSelect(event) {
  const file = event.target.files[0]; // Get the selected file
  if (file) {
    const reader = new FileReader(); // Create a new FileReader object
    reader.onload = function(e) {
      const favicon = document.createElement('link'); // Create a new link element for favicon
      favicon.rel = 'icon'; // Set rel attribute to 'icon' for favicon
      favicon.type = 'image/png'; // Set type attribute to 'image/png' for favicon
      favicon.href = e.target.result; // Set the href attribute to the selected image
      const head = document.querySelector('head'); // Get the <head> element
      head.appendChild(favicon); // Append the favicon link to the head
    };
    reader.readAsDataURL(file); // Read the selected file as a Data URL
  }
}

/** Changes the title of the page to value of element with id:"newTitle" */
function handleTitleChange() {
  const newTitle = document.getElementById('newTitle').value; // Get the value from the input field
  const titleElement = document.getElementById('title'); 
  titleElement.textContent = newTitle; // Set the text content to the new title
}