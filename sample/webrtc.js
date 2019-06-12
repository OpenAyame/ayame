// シグナリングサーバのURLを指定する
let wsUrl = 'ws://localhost:3000/ws';
const roomStorageKey = "OPEN-AYAME-SAMPLE-ROOM-IDS";
const roomInput = document.getElementById("roomId");
const recentRoomDiv = document.getElementById("recent-rooms");
const clientId = randomString(17);
const localVideo = document.getElementById('local-video');
const remoteVideo = document.getElementById('remote-video');
const connectButton = document.getElementById('connect-button');
const disconnectButton = document.getElementById('disconnect-button');
document.getElementById("url").value = wsUrl;
let ws = null;
let roomId = randomString(9);
roomInput.value = roomId;
let roomIds = [];
let localStream = null;
let peerConnection = null;
const iceServers = [{ 'urls': 'stun:stun.l.google.com:19302' }];
const peerConnectionConfig = {
  'iceServers': iceServers
};
let isNegotiating = false;
disconnectButton.disabled = true;

// 接続処理
function connect() {
  roomId = document.getElementById("roomId").value;
  if (roomId.length < 2 || !roomId){
    alert("部屋のID を指定してください");
    return;
  }
  let newRoomIds = [];
  if(roomIds.length > 0 && roomId === roomIds[0]) newRoomIds = [...roomIds];
  else {
    newRoomIds = [roomId, ...roomIds];
  }
  localStorage.setItem(roomStorageKey, JSON.stringify(newRoomIds));
  recentRoomDiv.style.display = 'none';
  isNegotiating = false;
  // 新規に websocket を作成
  if(!ws){
    ws = new WebSocket(wsUrl);
  }
  // ws のコールバックを定義する
  ws.onopen = (event) => {
    console.log('ws open()');
    ws.send(JSON.stringify({
      "type": "register",
      "room_id": roomId,
      "client_id": clientId
    }))
    ws.onmessage = (event) => {
      console.log('ws onmessage() data:', event.data);
      const message = JSON.parse(event.data);
      console.log(message.type)
      switch(message.type){
        case 'ping': {
          console.log('Received Ping, Send Pong.');
          ws.send(JSON.stringify({
            "type": "pong"
          }))
          break;
        }
        case 'offer': {
          console.log('Received offer ...');
          setOffer(message);
          break;
        }
        case 'answer': {
          console.log('Received answer ...');
          setAnswer(message);
          break;
        }
        case 'candidate': {
          console.log('Received ICE candidate ...');
          const candidate = new RTCIceCandidate(message.ice);
          console.log(candidate);
          addIceCandidate(candidate);
          break;
        }
        case 'close': {
          console.log('peer is closed ...');
          disconnect();
          break;
        }
        case 'reject': {
          console.log('connection is rejected...');
          disconnect();
          break;
        }
        case 'accept': {
          connectButton.disabled = true;
          disconnectButton.disabled = false;
          if (!peerConnection) {
            console.log('make Offer');
            peerConnection = prepareNewConnection(true);
          }
          else {
            console.warn('peer already exist.');
          }
          break;
        }
        default: {
          console.log('Invalid message type: ');
          break;
        }
      }
    };

  };
  ws.onerror = (error) => {
    console.error('ws onerror() ERROR:', error);
  };
  ws.onclose = (event) => {
    disconnect();
  };
}

// 切断処理
function disconnect(){
  connectButton.disabled = false;
  disconnectButton.disabled = true;
  if (peerConnection) {
    if(peerConnection.iceConnectionState !== 'closed'){
      // peer connection を閉じる
      peerConnection.close();
      cleanupVideoElement(remoteVideo);
    }
  }
  if(ws && ws.readyState < 2){
    ws.close();
  }
  ws = null;
  isNegotiating = false;
  peerConnection = null;
  recentRoomDiv.style.display = 'block';
  loadLocalRoomIds();
  console.log('peerConnection is closed.');
}

// ws url の変更
function onChangeWsUrl() {
  wsUrl = document.getElementById("url").value;
  console.log('ws url changes', wsUrl);
}

function onChangeRoomId() {
  roomId = roomInput;
  console.log('room id changes', roomId);
}

// ICE candaidate受信時にセットする
function addIceCandidate(candidate) {
  if (peerConnection) {
    peerConnection.addIceCandidate(candidate);
  }
  else {
    console.error('PeerConnection does not exist!');
    return;
  }
}

function sendIceCandidate(candidate) {
  if (ws) {
  console.log('---sending ICE candidate ---', candidate);
  const message = JSON.stringify({ type: 'candidate', ice: candidate});
  console.log('sending candidate=' + message);
    ws.send(message);
  }
  else {
    console.error('websocket connection does not exist!');
  }
}

async function startVideo() {
  try{
    localStream = await navigator.mediaDevices.getUserMedia({video: true, audio: true});
    playVideo(localVideo,localStream);
  } catch(error){
    console.error('mediaDevice.getUserMedia() error:', error);
  }
}

async function playVideo(element, stream) {
  element.srcObject = stream;
  try {
    await element.play();
  } catch(error) {
    console.log('error auto play:' + error);
  }
}

function prepareNewConnection(isOffer) {
  const peer = new RTCPeerConnection(peerConnectionConfig);
  if ('ontrack' in peer) {
    let tracks = [];
    peer.ontrack = event => {
      tracks.push(event.track);
      console.log('-- peer.ontrack()', event);
      let mediaStream = new MediaStream(tracks);
      playVideo(remoteVideo, mediaStream);
    };
  }else{
    peer.onaddstream = event => {
      console.log('-- peer.onaddstream()');
      event.stream.onaddtrack = e => {
        event.track.onended = event => remoteVideo.srcObject = remoteVideo.srcObject;
      }
      playVideo(remoteVideo, event.stream);
    };
  }

  peer.onicecandidate = event => {
    if (event.candidate) {
      console.log('-- peer.onicecandidate()', event.candidate);
      sendIceCandidate(event.candidate);
    } else {
      console.log('empty ice event');
    }
  };

    peer.onnegotiationneeded = async () => {
      if (isNegotiating) {
        console.log("SKIP nested negotiations");
        return;
      }
      try {
        isNegotiating = true;
        if(isOffer){
          const offer = await peer.createOffer({
            'offerToReceiveAudio': true,
            'offerToReceiveVideo': true
          })
          console.log('createOffer() succsess in promise');
          await peer.setLocalDescription(offer);
          console.log('setLocalDescription() succsess in promise');
          sendSdp(peer.localDescription);
          isNegotiating = false;
      }
      } catch(error){
      console.error('setLocalDescription(offer) ERROR: ', error);
    }
  }

  // ICEのステータスが変更になったときの処理
  peer.oniceconnectionstatechange = () => {
    console.log('ICE connection Status has changed to ' + peer.iceConnectionState);
    switch (peer.iceConnectionState) {
      case 'connected':
        isNegotiating = false;
        break;
      case 'closed':
      case 'failed':
      case 'disconnected':
        cleanupVideoElement(remoteVideo);
        disconnect();
        break;
    }
  };
  peer.onsignalingstatechange = (e) => {
    console.log('signaling state changes:', peer.signalingState);
  }
  // ローカルのMediaStreamを利用できるようにする
  if (localStream) {
    console.log('Adding local stream...');
    const videoTrack = localStream.getVideoTracks()[0];
    const audioTrack = localStream.getAudioTracks()[0];
    if(videoTrack){
      peer.addTrack(videoTrack, localStream);
    }
    if(audioTrack){
      peer.addTrack(audioTrack, localStream);
    }
  } else {
    console.warn('no local stream, but continue.');
  }

  if (isUnifiedPlan(peer)) {
    console.log('peer is unified plan');
    peer.addTransceiver('video', {direction: 'recvonly'});
    peer.addTransceiver('audio', {direction: 'recvonly'});
  }
  return peer;
}


function isUnifiedPlan(peer) {
  const config = peer.getConfiguration();
  return ('addTransceiver' in peer) && (!('sdpSemantics' in config) || config.sdpSemantics === "unified-plan");
}

// sdp を ws で送る
function sendSdp(sessionDescription) {
  if(ws){
  console.log('---sending sdp ---');
  const message = JSON.stringify(sessionDescription);
  console.log('sending SDP=' + message);
    ws.send(message);
  }
  else {
    console.error('websocket connection does not exist!');
  }
}


// Answer SDP を生成する
async function makeAnswer() {
  console.log('sending Answer. Creating remote session description...' );
  if (!peerConnection) {
    console.error('peerConnection NOT exist!');
    return;
  }
  try{
    let answer = await peerConnection.createAnswer();
    console.log('createAnswer() succsess in promise');
    await peerConnection.setLocalDescription(answer);
    console.log('setLocalDescription() succsess in promise');
    sendSdp(peerConnection.localDescription);
  } catch(error){
    console.error(error);
  }
}

// Offer SDP を生成する
async function setOffer(sessionDescription) {
  peerConnection = prepareNewConnection(false);
  try{
    await peerConnection.setRemoteDescription(sessionDescription);
    console.log('setRemoteDescription(answer) success in promise');
    makeAnswer();
  } catch(error){
    console.error('setRemoteDescription(offer) ERROR: ', error);
  }
}

// Answer SDP を生成する
async function setAnswer(sessionDescription) {
  if (!peerConnection) {
    console.error('peerConnection DOES NOT exist!');
    return;
  }
  try{
    await peerConnection.setRemoteDescription(sessionDescription);
    console.log('setRemoteDescription(answer) succsess in promise');
  } catch(error){
    console.error('setRemoteDescription(answer) ERROR: ', error);
  }
}

// video element を初期化する
function cleanupVideoElement(element) {
  let stream = element.srcObject;
  if(stream){
    let tracks = stream.getTracks();
    tracks.forEach(function(track) {
      track.stop();
    });
    element.srcObject = null;
  }
}

// Return a random numerical string.
function randomString(strLength) {
  var result = [];
  strLength = strLength || 5;
  var charSet = '0123456789';
  while (strLength--) {
    result.push(charSet.charAt(Math.floor(Math.random() * charSet.length)));
  }
  return result.join('');
}

startVideo();


// ここからルームID の取得の処理

function setRoomId(r) {
  roomId = r;
  roomInput.value = roomId;
}

function loadLocalRoomIds() {
  const roomUl = document.getElementById('recent-item-list');
  const itemJSON = localStorage.getItem(roomStorageKey);
  if (itemJSON) {
    roomIds = JSON.parse(itemJSON).slice(0, 7);
  }
  const fragment = document.createDocumentFragment();
  roomUl.innerHTML = '';
  roomIds.forEach(r => {
    const roomLi = document.createElement('li');
    roomLi.id = 'recent-items-' + roomId;
    roomLi.innerHTML = `<a onclick="setRoomId(${r})">${r}</a>`;
    fragment.appendChild(roomLi);
  });

  roomUl.appendChild(fragment);
}

loadLocalRoomIds();
