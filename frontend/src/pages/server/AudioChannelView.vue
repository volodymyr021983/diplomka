<template>

<video ref="clientVideoElement" autoplay muted controls width="320" height="240"></video>
<div id="container"></div>
<button @click="ConnectToVoice">CONNECT</button>
<button @click="GetClients">GETCLIENTS</button>
<button @click="Disconnect">DISCONNECT</button>


</template>

<script setup lang="js">
import {ref} from 'vue'
import { useRoute } from 'vue-router'
import Session from 'supertokens-web-js/recipe/session';

const route = useRoute()

const serverVideoElement = ref(null)
const clientVideoElement = ref(null)
const localClientStream = ref(null)

const serverID = route.params.server_id
const channelID = route.params.channel_id
const userId = ref('')

const PeerConnection = ref(null)
const offer = ref(null)
const answer = ref(null)
const channelTracksElements = new Map()

const wsSignalingConn = new WebSocket(`${import.meta.env.VITE_WSS_API_URL}/api/signaling/${serverID}/${channelID}`)
//very important listener where code will process the messages from signaling server as offer answer etc
wsSignalingConn.addEventListener("message", async (event) =>{
  const msg = JSON.parse(event.data)
  switch(msg.type){
   
    case "connected":
      wsSignalingConn.send(JSON.stringify({type: "connected"}))
      console.log("websocket connected")
      break;
    case "new_ice_candidate":
      if(msg.payload != null){
      const arrivedUserId = msg.userid
      console.log("ice candidate from:", arrivedUserId)
        await PeerConnection.value.addIceCandidate(msg.payload);   
    }
      break;
    case "conn_answer":
      console.log("payload:")
      console.log(msg.payload)
      const answerOf = new RTCSessionDescription(msg.payload)
      await PeerConnection.value.setRemoteDescription(answerOf)
      break;
    case "conn_offer":
      console.log("offer received payload:")

      const offerAns = new RTCSessionDescription(msg.payload)
      await PeerConnection.value.setRemoteDescription(offerAns)
       answer.value = await PeerConnection.value.createAnswer()
      await PeerConnection.value.setLocalDescription(answer.value)
      sendAnswer()
      break;
  }
})

//function responsible for RTCPeerConnection creation, adding tracks
async function CreatePeerConnection(){
  PeerConnection.value = new RTCPeerConnection({iceServers:[{
    urls:[
      'stun:stun.l.google.com:19302','stun:stun1.l.google.com:19302'
    ]
  }]})
  await GetUserMedia()

  localClientStream.value.getTracks().forEach(track =>{
    PeerConnection.value.addTrack(track, localClientStream.value)
  })
  //listener for responsible for sending ice candidates to the signalling server
  PeerConnection.value.addEventListener("icecandidate", (event)=>{
    console.log("icecandidate found!")
    wsSignalingConn.send(JSON.stringify({type: 'new_ice_candidate', userid: userId.value, payload: event.candidate}))
  })
  PeerConnection.value.addEventListener("iceconnectionstatechange", (event)=>{
    console.log(`icecandidate state change: ${PeerConnection.value.iceConnectionState}`)
  })
  PeerConnection.value.addEventListener("connectionstatechange", (event)=>{
    console.log(`connection state change: ${PeerConnection.value.connectionState}`)
    if(PeerConnection.value.connectionState == 'disconnected'){
      PeerConnection.value.close()
    }
  })
  PeerConnection.value.addEventListener("negotiationneeded", (event)=>{
    console.log(`negotiation needed`)
  })
  PeerConnection.value.addEventListener("track", (event) =>{
    console.log("track arrived")
     console.log("Number of streams in event:", event.streams.length);
     event.streams.forEach(stream =>{
      const video = document.createElement('video')
        video.setAttribute('autoplay', '');   // или video.autoplay = true;
        video.setAttribute('controls', '');
        video.width = 320;
        video.height = 240;
        video.srcObject = stream
        document.getElementById('container').appendChild(video);
        channelTracksElements.set(stream, video)
     })
  })
}
async function GetClients() {
  wsSignalingConn.send(JSON.stringify({ type: 'get_clients'}));
}
//function responsible for fetching user devises such as microphone and camera
async function GetUserMedia() {
  try {
    localClientStream.value = await navigator.mediaDevices.getUserMedia({ video: false, audio: true })

    if (clientVideoElement.value) {
      clientVideoElement.value.srcObject = localClientStream.value
    }
  } catch (err) {
    console.error("Error getting user media:", err)
    throw err 
  }
}

async function ConnectToVoice(){
    try{
    userId.value = await Session.getUserId()
    await CreatePeerConnection()
       offer.value = await PeerConnection.value.createOffer()
       await PeerConnection.value.setLocalDescription(offer.value)
       await sendJoinRequest(); 
       await sendOffer()
       

    }catch(err){
        alert(err)
    }
}
async function sendAnswer(){
  console.log("answer sended")
  wsSignalingConn.send(JSON.stringify({type: 'conn_answer', payload: answer.value}))
}
async function sendOffer(){
    console.log("sending offer");
    wsSignalingConn.send(JSON.stringify({ type: 'conn_offer', payload: offer.value}));
}
async function sendJoinRequest(){
  console.log("Sending join_room request");
    wsSignalingConn.send(JSON.stringify({ type: 'join_channel'}));
}
async function Disconnect(){
  if(PeerConnection.value != null){
  wsSignalingConn.send(JSON.stringify({type: 'disconnect_channel'}));
  }
}
async function CreateVideoElement(){
  const newVid = document.createElement("video")
}
</script>

<style scoped>
video{
border:2px solid blue;
padding:5px;
/**Or add your own style**/
}
</style>