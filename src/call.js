import {
  addToPresenceList,
  hidePresence,
  setupCreateButton,
  showCall,
  showEntry,
  showPresence,
} from './dom';

window.addEventListener('DOMContentLoaded', () => {
  setupCreateButton(createRoom);
  showEntry();
});

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
    let msg = `${errMsg}. Status code: ${res.status}`;
    if (body) {
      msg += `; ${body}`;
    }
    console.error(`failed to create Daily room: ${msg}`);
  }
  joinRoom(body.url, body.name);
}

function joinRoom(roomURL, roomName) {
  showCall();

  // Set up the call frame and join the call
  const callContainer = document.getElementById('callFrame');
  const callFrame = window.DailyIframe.createFrame(callContainer, {
    showLeaveButton: true,
    iframeStyle: {
      position: 'fixed',
      width: 'calc(100% - 1rem)',
      height: 'calc(100% - 5rem)',
    },
  });
  callFrame.on('joined-meeting', () => {
    // Presence view is no longer needed once the participant is in the room.
    hidePresence();
  });
  callFrame.join({ url: roomURL });

  // Fetch existing participants and show them next to the call frame
  fetchParticipants(roomName);
}

async function fetchParticipants(roomName) {
  const url = `/.netlify/functions/presence?`;
  const errMsg = 'failed to fetch room presence';

  // Create Daily room
  let res;
  try {
    res = await fetch(
      url +
        new URLSearchParams({
          roomName,
        })
    );
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
    let msg = `${errMsg}. Status code: ${res.status}`;
    if (body) {
      msg += `; ${body}`;
    }
    console.error(`failed to fetch presence for Daily room: ${msg}`);
  }
  for (let i = 0; i < body.length; i += i) {
    const p = body[i];

    const label = p.name ? p.name : p.id;
    addToPresenceList(label);
  }
  showPresence();
}
