import * as htmlToImage from 'html-to-image';
import imageCompression from 'browser-image-compression';
import {LogError as logError} from '../../wailsjs/go/main/Tray';
import {getUserSettings} from './issues';

export async function getImage(node, width, height) {
  console.log(width,height);
  // create Blob from node
  // facebook standard for shared images is 1200x630 or 600x315 (1.91:1)
  const imageOptions = { width: width, height: height }
  const imageBlob = (await htmlToImage.toBlob(node, imageOptions))

  // pass Blob and the quality option to be compressed
  const compressionOptions = { initialQuality: 1 }
  const compressedFile = await imageCompression(imageBlob, compressionOptions)

  // create URL from Blob
  return window.URL.createObjectURL(compressedFile)
}

export async function saveProgress(node) {
  try {
    const social = JSON.parse(sessionStorage.getItem('ShareSocial'));
    const imageUrl = await getImage(node,social.height,social.width);

    var nowDate = new Date(); 
    var date = nowDate.getDate()+'/'+(nowDate.getMonth()+1)+'/'+nowDate.getFullYear();
    // change date if localization is en-US
    const language = await getUserSettings();
    if (language == 2) date = (nowDate.getMonth()+1)+'/'+nowDate.getDate()+'/'+nowDate.getFullYear();

    // download image
    const linkElement = document.createElement('a')
    linkElement.download = 'Info-Sec-Agent_'+date+'.png'
    linkElement.href = imageUrl
    linkElement.click()
  } catch (error) {
    throw new Error(`Something went wrong: ${error}`)
  }
}

export async function shareProgress(node) {
  const social = JSON.parse(sessionStorage.getItem('ShareSocial'));
  const imageUrl = await getImage(node,social.height,social.width);
  

  let socialUrl = '';
  switch (social.name) {
    case 'facebook' :
      console.log(social.name);
      break;
    case 'x' :
      console.log(social.name);
      break;
    case 'linkedin' :
      console.log(social.name);
      break;
    case 'instagram' :
      console.log(social.name);
      break;
    default:
      console.log(social.name);
      break;
  }
  if (socialUrl == '') {
    logError('Sharing failed: incorrect url');
    return
  }

  // window.open('https://nl.wikipedia.org/wiki/Hoofdpagina','_blank');
}

const socialMediaSizes = {
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
}

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

  console.log(social);
  document.getElementById('select-' + social).classList.add('selected');
  sessionStorage.setItem('ShareSocial', JSON.stringify(socialMediaSizes[social]));
}