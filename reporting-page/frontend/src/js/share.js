import * as htmlToImage from 'html-to-image';
import imageCompression from 'browser-image-compression';
import {LogError as logError} from '../../wailsjs/go/main/Tray';
import {getUserSettings} from './issues';

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
 */
function setImage(social) {
  const node = document.getElementById('share-node');
  node.innerHTML = `
  <img class="api-key-image" src="https://placehold.co/600x315" alt="Step 1 Image">
  `;
  if (social == 'instagram') {
    node.innerHTML = `
    <img class="api-key-image" src="https://placehold.co/300x300" alt="Step 1 Image">   `;
  }
}

/**
 * Save the image created from a node in the downloads folder with selected social media format
 * @param {HTMLElement} node node to download the image for
 */
export async function saveProgress(node) {
  try {
    const social = JSON.parse(sessionStorage.getItem('ShareSocial'));
    const imageUrl = await getImage(node, social.height, social.width);

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
export function shareProgress() {
  const social = JSON.parse(sessionStorage.getItem('ShareSocial'));

  // choose which
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
    height: 600,
    width: 315,
  },
  x: {
    name: 'x',
    height: 600,
    width: 315,
  },
  linkedin: {
    name: 'linkedin',
    height: 600,
    width: 315,
  },
  instagram: {
    name: 'instagram',
    height: 300,
    width: 300,
  },
};

// on startup set the social media to share to to facebook
sessionStorage.setItem('ShareSocial', JSON.stringify(socialMediaSizes['facebook']));

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
  setImage(social);
  sessionStorage.setItem('ShareSocial', JSON.stringify(socialMediaSizes[social]));
}
