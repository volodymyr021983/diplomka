<template>

<video ref="clientVideoElement" autoplay controls width="320" height="240"></video>

<button @click="ConnectToVoice">CONNECT</button>
<button @click="GetClients">GETCLIENTS</button>

</template>

<script setup lang="js">
import {ref} from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()

const clientVideoElement = ref(null)
const localClientStream = ref(null)

const serverID = route.params.server_id
const channelID = route.params.channel_id

const PeerConnection = ref(null)
const offer = ref(null)

const wsSignalingConn = new WebSocket(`${import.meta.env.VITE_WSS_API_URL}/api/signaling/${serverID}/${channelID}`)
//very important listener where code will process the messages from signaling server as offer answer etc
wsSignalingConn.addEventListener("message", (event) =>{
  console.log(event.data)
  const msg = JSON.parse(event.data)
  switch(msg.type){
   
    case "connected":
      wsSignalingConn.send(JSON.stringify({type: "connected"}))
      console.log("websocket connected")
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
    wsSignalingConn.send(JSON.stringify({type: 'new_ice_candidate', payload: event.candidate}))
  })
  PeerConnection.value.addEventListener("iceconnectionstatechange", (event)=>{
    console.log(`icecandidate state change: ${event.iceConnectionState}`)
  })
  PeerConnection.value.addEventListener("connectionstatechange", (event)=>{
    console.log(`connection state change: ${event.connectionState}`)
  })
}
async function GetClients() {
  wsSignalingConn.send(JSON.stringify({ type: 'get_clients'}));
}
//function responsible for fetching user devises such as microphone and camera
async function GetUserMedia() {
  try {
    localClientStream.value = await navigator.mediaDevices.getUserMedia({ video: true })

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
    await CreatePeerConnection()
       offer.value = await PeerConnection.value.createOffer()
       await PeerConnection.value.setLocalDescription(offer.value)
       sendOffer()
        sendJoinRequest();
    }catch(err){
        alert(err)
    }
}
async function sendOffer(){
    console.log("sending offer");
    wsSignalingConn.send(JSON.stringify({ type: 'conn_offer', payload: offer.value}));
}
async function sendJoinRequest(){
  console.log("Sending join_room request");
    wsSignalingConn.send(JSON.stringify({ type: 'join_channel'}));
}
</script>

<style scoped>
video{
border:2px solid blue;
padding:5px;
/**Or add your own style**/
}
</style>