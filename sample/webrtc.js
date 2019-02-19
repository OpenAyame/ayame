// シグナリングサーバのURLを指定する
let wsUrl = 'ws://localhost:3000/ws';
document.getElementById("url").value = wsUrl;
let ws = null;
let roomId = randomString(9);
document.getElementById("roomId").value = roomId;
const clientId = randomString(17);
const localVideo = document.getElementById('local-video');
const remoteVideo = document.getElementById('remote-video');
let localStream = null;
let peerConnection = null;
const iceServers = [{ 'urls': 'stun:stun.l.google.com:19302' }];
const peerConnectionConfig = {
  'iceServers': iceServers
};
let isNegotiating = false;


// 接続処理
function connect() {
  roomId = document.getElementById("roomId").value;
  if (roomId.length < 2 || !roomId){
    alert("部屋のID を指定してください");
    return;
  }
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
      "room_id": roomId
    }))
    ws.onmessage = (event) => {
      console.log('ws onmessage() data:', event.data);
      const message = JSON.parse(event.data);
      switch(message.type){
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
        default: {
          console.log('Invalid message type: ');
          break;
        }
      }
    };

    if (!peerConnection) {
      console.log('make Offer');
      peerConnection = prepareNewConnection(true);
    }
    else {
      console.warn('peer already exist.');
    }
  };
  ws.onerror = (error) => {
    console.error('ws onerror() ERROR:', error);
  };
}

// 切断処理
function disconnect(){
  if (peerConnection) {
    if(peerConnection.iceConnectionState !== 'closed'){
      // peer connection を閉じる
      peerConnection.close();
      peerConnection = null;
      const message = JSON.stringify({ type: 'close'});
      console.log('sending close message');
      if(ws) {
        ws.send(message);
        ws.close();
        ws = null;
      }
      else {
        console.error('websocket connection does not exist!');
      }
      cleanupVideoElement(remoteVideo);
      return;
    }
  }
  console.log('peerConnection is closed.');
}

// ws url の変更
function onChangeWsUrl() {
  wsUrl = document.getElementById("url").value;
  console.log('ws url changes', wsUrl);
}

function onChangeRoomId() {
  roomId = document.getElementById("roomId").value;
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
          const offer = await peerConnection.createOffer({
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
        if (peerConnection) {
          disconnect();
        }
        break;
      case 'disconnected':
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

  return peer;
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
  if (peerConnection) {
    console.error('peerConnection already exist!');
  }
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
