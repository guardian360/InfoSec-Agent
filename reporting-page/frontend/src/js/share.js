import * as htmlToImage from 'html-to-image';
import imageCompression from 'browser-image-compression';
import {LogError as logError} from '../../wailsjs/go/main/Tray';
import {getUserSettings} from './issues';
import {GetImagePath as getImagePath, GetLighthouseState as getLighthouseState} from '../../wailsjs/go/main/App.js';

/**
 * Create image as an url from an html node
 * @param {HTMLElement} node html node to turn into an image
 * @param {Int} width width of the resulting image
 * @param {Int} height height of the resulting image
 * @return {URL} url of created image
 */
export async function getImage(node, width, height) {
  // create Blob from node
  // facebook standard for shared images is 1200x630 or 600x315 (1.91:1)
  const imageOptions = {width: width, height: height};
  const imageBlob = (await htmlToImage.toBlob(node, imageOptions));

  // pass Blob and the quality option to be compressed
  const compressionOptions = {initialQuality: 1};
  const compressedFile = await imageCompression(imageBlob, compressionOptions);

  // create URL from Blob
  return window.URL.createObjectURL(compressedFile);
}

/**
 * Set the correct image in the share-node
 * @param {string} social social media to set the image for
 * @param {HTMLElement} progress html element of the progress bar with text above it
 */
export async function setImage(social, progress) {
  let lighthousePath;
  const lighthouseState = await getLighthouseState();
  switch (lighthouseState) {
  case 0:
    lighthousePath = 'first-state.png';
    break;
  case 1:
    lighthousePath = 'second-state.png';
    break;
  case 2:
    lighthousePath = 'third-state.png';
    break;
  case 3:
    lighthousePath = 'fourth-state.png';
    break;
  case 4:
    lighthousePath = 'fifth-state.png';
    break;
  default:
    lighthousePath = 'first-state.png';
  }
  const lighthouseFullPath = await getImagePath('sharing/' + lighthousePath);
  const node = document.getElementById('share-node');
  const socialStyle = JSON.parse(sessionStorage.getItem('ShareSocial'));
  // Set the background to the current state
  node.style.width = socialStyle.width + 'px';
  node.style.height = socialStyle.height + 'px';
  node.style.backgroundImage = 'url(' + lighthouseFullPath + ')';
  node.style.backgroundSize = 'cover';
  node.style.backgroundPosition = 'center';

  node.innerHTML = `
  <div class="image-header">
    <p class="image-link">github.com/InfoSec-Agent/InfoSec-Agent</p>
  </div>
  <div class="image-footer">
    <div class="image-left" id="image-left">
      ${progress.innerHTML}
    </div>
    <div class="image-right">
      <p id="image-logo-text"></p>
      <img id="logo" alt="logo" src="./src/assets/images/InfoSec-Agent-logo.png" style="width: 75px; height: 75px;">
    </div>
  </div>
  `;
  if (social == 'instagram') {
    document.getElementById('image-left').style.marginTop = '50px';
  }
}

/**
 * Save the image created from a node in the downloads folder with selected social media format
 * @param {HTMLElement} node node to download the image for
 */
export async function saveProgress(node) {
  try {
    const social = JSON.parse(sessionStorage.getItem('ShareSocial'));
    const imageUrl = await getImage(node, social.width, social.height);

    const nowDate = new Date();
    let date = nowDate.getDate()+'-'+(nowDate.getMonth()+1)+'-'+nowDate.getFullYear();
    // change date if localization is en-US
    const language = await getUserSettings();
    if (language == 2) date = (nowDate.getMonth()+1)+'-'+nowDate.getDate()+'-'+nowDate.getFullYear();

    // download image
    const linkElement = document.createElement('a');
    linkElement.download = 'Info-Sec-Agent_'+date+'_'+social.name+'.png';
    linkElement.href = imageUrl;
    linkElement.click();
  } catch (error) {
    throw new Error(`saveProgress couldn't be completed: ${error}`);
  }
}

/** Open the selected social media page */
export async function shareProgress() {
  const social = JSON.parse(sessionStorage.getItem('ShareSocial'));

  // choose which social media page to open
  switch (social.name) {
  case 'facebook':
    window.open('https://www.facebook.com/', 'Facebook');
    break;
  case 'x':
    window.open('https://x.com/', 'X');
    break;
  case 'linkedin':
    window.open('https://www.linkedin.com/', 'Linkedin');
    break;
  case 'instagram':
    window.open('https://www.instagram.com/', 'Instagram');
    break;
  default:
    logError('Sharing failed: social media link not present');
    break;
  }
}

// Different social media's to share to, with specifications for image size
export const socialMediaSizes = {
  facebook: {
    name: 'facebook',
    height: 315,
    width: 600,
  },
  x: {
    name: 'x',
    height: 315,
    width: 600,
  },
  linkedin: {
    name: 'linkedin',
    height: 315,
    width: 600,
  },
  instagram: {
    name: 'instagram',
    height: 300,
    width: 300,
  },
};

/**
 * Select the social media and set it in sessionstorage
 * @param {string} social social media to share to
 */
export function selectSocialMedia(social) {
  const socials = document.getElementsByClassName('select-button');
  for (let i = 0; i < socials.length; i++) {
    socials[i].classList.remove('selected');
  };

  document.getElementById('select-' + social).classList.add('selected');
  setImage(social, document.getElementById('progress-segment'));
  sessionStorage.setItem('ShareSocial', JSON.stringify(socialMediaSizes[social]));
}
