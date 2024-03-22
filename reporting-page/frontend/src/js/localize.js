import { Localize } from '../../wailsjs/go/main/App';

export function GetLocalization(messageId, elementId) {
    Localize(messageId).then((result) => {
        document.getElementById(elementId).innerHTML = result;
    });
}