import {
  addToPresenceList,
  hidePresence,
  setupCreateButton,
  showCall,
  showEntry,
  showPresence,
} from './dom.js';

window.addEventListener('DOMContentLoaded', () => {
  setupCreateButton(createRoom);

  // If room URL and name params were provided, join that room
  const params = new URLSearchParams(window.location.search);
  const roomURL = params.get('roomURL');
  const roomName = params.get('roomName');
  if (roomURL && roomName) {
    joinRoom(roomURL, roomName);
    return;
  }

  // If room URL and name params were not provided, show
  // entry element (with call creation button)
  showEntry();
});

/**
 * Creates and joins a Daily room
 * @returns {Promise<void>}
 */
async function createRoom() {
  const url = `/.netlify/functions/createRoom`;
  const errMsg = 'failed to create room';

  // Create Daily room
  let res;
  try {
    res = await fetch(url);
  } catch (e) {
    console.error('failed to make Daily room creation request', e);
  }

  // Retrieve body, if it exists in the response
  let body;
  try {
    body = await res.json();
  } catch (e) {
    console.warn('No body present in room creation response');
  }

  // Check status code
  if (res.status !== 200) {
    const msg = `${errMsg}. Status code: ${res.status}`;
    console.error(`failed to create Daily room: ${msg}`, body);
    return;
  }
  joinRoom(body.url, body.name);
}

/**
 * Joins a Daily room and retrieves its existing participants
 * @param roomURL
 * @param roomName
 */
function joinRoom(roomURL, roomName) {
  showCall();

  // Set up the call frame and join the call
  const callContainer = document.getElementById('callFrame');
  const callFrame = window.DailyIframe.createFrame(callContainer, {
    showLeaveButton: true,
  });
  callFrame.on('joined-meeting', () => {
    // Presence view is no longer needed once the participant is in the room.
    hidePresence();
  });
  callFrame.join({ url: roomURL });

  updateInviteURL(roomURL, roomName);

  // Fetch existing participants and show them next to the call frame
  fetchParticipants(roomName);
}

function updateInviteURL(roomURL, roomName) {
  const inviteURLEle = document.getElementById('inviteURL');

  const url = new URL(window.location.href);
  const params = new URLSearchParams();
  params.set('roomURL', roomURL);
  params.set('roomName', roomName);
  url.search = params;
  const inviteURL = url.toString();
  inviteURLEle.href = inviteURL;
  inviteURLEle.innerText = inviteURL;
}

/**
 * Fetches a list of participants who are
 * already in the given room name.
 * @param roomName
 * @returns {Promise<void>}
 */
async function fetchParticipants(roomName) {
  const url = `/.netlify/functions/presence?`;
  const errMsg = 'failed to fetch room presence';

  let res;
  try {
    const reqURL =
      url +
      new URLSearchParams({
        roomName,
      });
    res = await fetch(reqURL);
  } catch (e) {
    console.error('failed to make Daily room presence request', e);
  }

  // Retrieve body, if it exists in the response
  let body;
  try {
    body = await res.json();
  } catch (e) {
    console.warn('No body present in room presence response');
  }

  // Check status code
  if (res.status !== 200) {
    const msg = `${errMsg}. Status code: ${res.status}`;
    console.error(`failed to fetch presence for Daily room: ${msg}`, body);
  }

  const count = body.length;
  if (count === 0) {
    addToPresenceList('Nobody here yet!');
    return;
  }
  for (let i = 0; i < body.length; i += 1) {
    const p = body[i];
    const label = p.userName ? p.userName : p.id;
    addToPresenceList(label);
  }
  showPresence();
}
