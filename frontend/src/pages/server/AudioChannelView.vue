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

const wsSignalingConn = new WebSocket(`${import.meta.env.VITE_WSS_API_URL}/api/signaling/${serverID}/${channelID}`)

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
async function GetClients() {
  wsSignalingConn.send(JSON.stringify({ type: 'get_clients'}));
}
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
    await GetUserMedia()
        sendJoinRequest();
    }catch(err){
        alert("Oopssomthing wrong")
    }
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