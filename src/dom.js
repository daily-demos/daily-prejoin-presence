const hiddenClassName = 'hidden';

export function setupCreateButton(onclick) {
  const createBtn = document.getElementById('createRoom');
  createBtn.onclick = () => {
    onclick();
  };
}

export function addToPresenceList(label) {
  const participants = document.getElementById('participants');
  const item = document.createElement('li');
  item.innerText = label;
  participants.appendChild(item);
}

export function showPresence() {
  getPresenceEle().classList.remove('hidden');
}

export function hidePresence() {
  getPresenceEle().classList.add('hidden');
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
