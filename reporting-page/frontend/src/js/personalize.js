import cs from "../customize.json" assert { type: "json" };
import { CloseNavigation } from "./navigation-menu.js";
import { MarkSelectedNavigationItem } from "./navigation-menu.js";
/** Load the content of the Personalize page */
export function openPersonalizePage() {
  CloseNavigation();
  //MarkSelectedNavigationItem("home-button");
  document.getElementById("page-contents").innerHTML = `
  <div class="setting">
    <span class="setting-description favicon-title ">Favicon</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="input-file-icon">Change favicon</label>
      <input class="personalize-input-invisible" type="file" id="input-file-icon" accept=".ico, .png">
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">Navigation image</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="input-file-picture">Update image</label>
      <input class="personalize-input-invisible" type="file" id="input-file-picture" accept="image/jpeg, image/png, image/jpg">   
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">Navigation title</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="newTitle">Update title</label>
      <input type="text" id="newTitle">
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">Font</span>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">Background color Left nav</span>
    <div class="personalize-button-container">
      <label class="personalize-label" for="input-color-background">Change background</label>
      <input class="personalize-input-invisible" type="color" id="input-color-background">   
    </div>
  </div>
  <hr class="solid">
  <div class="setting">
    <span class="personalize-description">text color</span>
  </div>
  <hr class="solid">
  <div class="setting">
    <form action="" class="color-picker>
      <fieldset>
        <legend>Pick a theme</legend>
        <label for="normal">normal</label>
        <input type="radio" name="theme" id="normal" checked>
        <label for="dark">Dark</label>
        <input type="radio" name="theme" id="dark">
        <label for="blue">Blue</label>
        <input type="radio" name="theme" id="blue">
  </div>
  `;
  const faviconInput = document.getElementById('input-file-icon');//add eventlistener for changing Favicon
  faviconInput.addEventListener('change', handleFaviconChange);
  
  const pictureInput = document.getElementById('input-file-picture'); //add eventlistener for changing navication picture
  pictureInput.addEventListener('change', handlePictureChange);
  
  const newTitleInput = document.getElementById('newTitle'); //add eventlistener for changing navigation title
  newTitleInput.addEventListener('input', handleTitleChange);

  const inputBackgroundNav = document.getElementById('input-color-background'); //add eventlistener for changing navigation title
  inputBackgroundNav.addEventListener('change', handleLeftBackgroundNav);

  
  /*save themes*/
  const themes = document.querySelectorAll('[name="theme"]');
  themes.forEach(themeOption => {
    themeOption.addEventListener("click", () => {
      localStorage.setItem("theme", themeOption.id);
    });
  });
  
  const activeTheme = localStorage.getItem("theme");
  themes.forEach(themeOption => {
    if(themeOption.id === activeTheme){
      themeOption.checked = true;
    }
  });
  document.documentElement.className= activeTheme;

}
  
/* Changes the favicon*/
export function handleFaviconChange(icon) {
  const file = icon.target.files[0]; // Get the selected file
  if (file) {
    const reader = new FileReader();
    reader.onload = function(e) {
      const picture = e.target.result;
      const favicon = document.querySelector('link[rel="icon"]');
      if(favicon){
        favicon.href = picture;
      }
      else{
        const newFavicon = document.createElement('link');
        newFavicon.rel = 'icon';
        newFavicon.href = picture;
        document.head.appendChild(newFavicon);
      }
      localStorage.setItem("favicon", picture);
    };
    reader.readAsDataURL(file); // Read the selected file as a Data URL
  }
}

/* Changes the navigation picture*/
export function handlePictureChange(picture) {
  const file = picture.target.files[0]; // Get the selected file
  const reader = new FileReader();
  reader.onload = function(e) {
    const logo = document.getElementById('logo');
    logo.src = e.target.result; // Set the source of the logo to the selected image
    localStorage.setItem("picture", e.target.result)
    };
  reader.readAsDataURL(file); // Read the selected file as a Data URL
}

/* Changes the title of the page to value of element with id:"newTitle" */
export function handleTitleChange() {
  const newTitle = document.getElementById('newTitle').value; // Get the value from the input field
  const titleElement = document.getElementById('title'); 
  titleElement.textContent = newTitle; // Set the text content to the new title
  localStorage.setItem("title", newTitle);
}

/*Change the left background of the navigation*/
export function handleLeftBackgroundNav(){
  const colorPicker = document.getElementById('input-color-background');
  const color = colorPicker.value;
  let temp = document.getElementsByClassName('left-nav')[0];
  temp.style.backgroundColor = color;
}

export function retrieveTheme(){
  const activeTheme = localStorage.getItem("theme");
  document.documentElement.className = activeTheme;
}

//achtergrond navigation
//normale achtergrond
//kleur text
//text font
//kleur buttons