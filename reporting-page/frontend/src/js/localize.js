import {Localize} from '../../wailsjs/go/main/App';
/**
 * Retrieves localized message using the provided message ID and updates HTML elements with the specified class.
 * @param {string} messageId - The ID of the message to be localized.
 * @param {string} elementClass - The class name of HTML elements to be updated with the localized message.
 */
export function getLocalization(messageId, elementClass) {
  Localize(messageId).then((result) => {
    const elements = document.getElementsByClassName(elementClass);
    for (let i = 0; i < elements.length; i++) {
      elements[i].innerHTML = result;
    }
  });
}
