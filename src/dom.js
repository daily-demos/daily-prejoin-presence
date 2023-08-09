const hiddenClassName = 'hidden';

/**
 * Sets up the onclick handler for the call creation button.
 * @param onclick
 */
export function setupCreateButton(onclick) {
  const createBtn = document.getElementById('createRoom');
  createBtn.onclick = () => {
    onclick();
  };
}

/**
 * Adds a participant label to the presence list.
 * @param label
 */
export function addToPresenceList(label) {
  const participants = getParticipantsEle();
  const item = document.createElement('li');
  item.innerText = label;
  participants.appendChild(item);
}

export function clearPresence() {
  const participants = getParticipantsEle();
  participants.innerHTML = '';
}

export function showPresence() {
  getPresenceEle().classList.remove(hiddenClassName);
}

export function showEntry() {
  hideCallEle();
  showEntryEle();
}

export function showCall() {
  hideEntryEle();
  showCallEle();
}

function hideEntryEle() {
  getEntryEle().classList.add(hiddenClassName);
}

function showEntryEle() {
  getEntryEle().classList.remove(hiddenClassName);
}

function hideCallEle() {
  getCallEle().classList.add(hiddenClassName);
}

function showCallEle() {
  getCallEle().classList.remove(hiddenClassName);
}

function getEntryEle() {
  return document.getElementById('entry');
}

function getCallEle() {
  return document.getElementById('call');
}

function getPresenceEle() {
  return document.getElementById('presence');
}

function getParticipantsEle() {
  return document.getElementById('participants');
}
