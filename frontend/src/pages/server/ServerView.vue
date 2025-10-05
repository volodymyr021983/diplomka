<template>
  <div class="container">
    <button class="primary-btn" @click="DeleteServer">DeleteServer</button>
    <button class="primary-btn" @click="showPopUp = !showPopUp">Create Channel</button>
    <CreateChannelComponent 
      :visible="showPopUp" 
      :serverId="serverId" 
      @close="showPopUp = false"
    />

    <div class="message-section">
      <input v-model="inputText" class="input-field" placeholder="Type a message..." />
      <button class="primary-btn" @click="Sending">Send Message</button>
    </div>
   
    <pre ref="chat" id="chat"></pre>
   
    <ul class="channel-list">
      <li v-for="(channel, index) in channels" :key="index">
        <button @click="Redirect(channel.ChannelId)" class="server-btn">
          {{ channel.Channelname }}
        </button>
        <button @click="DeleteChannel(channel.ChannelId)">Delete Channel</button>
      </li>
    </ul>
    
    <InviteLinkComponent
    :serverId="serverId"
    />
    </div>

</template>

<script setup lang="js">
import { useRoute, useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'
import { doesSessionExist } from 'supertokens-web-js/recipe/session'
import CreateChannelComponent from '@/components/server/CreateChannelComponent.vue'
import InviteLinkComponent from '@/components/server/InviteLinkComponent.vue'
const showPopUp = ref(false)

const inputText = ref('')
const route = useRoute()
const router = useRouter()
const isSessionExists = ref(false)
const channels = ref([])
const chat = ref(null);
let ws = null
const serverId = route.params.server_id
const channelId = route.params.channel_id

async function Sending() {
  ws.send(inputText.value)
}

async function CloseConnIfExists() {
  if (ws && ws.OPEN) {
    ws.close()
    ws.send("CONNECTION CLOSED")
  }
}
async function DeleteServer(){
  try{
  const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/delete-server/${serverId}`,{
      method: "GET",
      headers: { "Content-Type": "application/json" },
    })
    if(response.ok){
      console.log("deleted!")
    }
  }catch(err){
    alert(err)
  }
}
async function DeleteChannel(ChannelId){
  try{
    console.log(ChannelId)
    const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/delete-channel/${serverId}`,{
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        Channel_id: ChannelId
      })
    })
  }catch(err){
    alert(err)
  }
}
async function Redirect(ChannelId) {
  router.push(`/server/${serverId}/${ChannelId}`)
  chat.value.innerText = ''
  CloseConnIfExists(ws)
  ConnectToChannel(ChannelId)
}

async function ConnectToChannel(ChannelId) {
  const wsUri = `ws://localhost:8080/api/server/connect/${serverId}/${ChannelId}`
  ws = new WebSocket(wsUri);
  ws.addEventListener("message", (e) => {
    chat.value.innerText += e.data + "\n";
    if(e.data == "close") ws.close()
  })
}

async function GetServerChannels() {
  try {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/getChannels/${serverId}`, {
      headers: { "Content-Type": "application/json" }
    })
    if(response.ok){
      const data = await response.json();
      channels.value = data
    }
  } catch(err) {
    alert(err)
  }
}

onMounted(async () => {
  await doesSessionExistWrapper()
  if(isSessionExists.value){
    GetServerChannels()
    ConnectToChannel(channelId)
  }
})

async function doesSessionExistWrapper() {
  isSessionExists.value = await doesSessionExist()
}
</script>

<style scoped>
.container {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}

.primary-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 16px;
  margin: 5px 0;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.primary-btn:hover {
  background-color: #45a049;
}

.secondary-btn {
  background-color: #008CBA;
  color: white;
  border: none;
  padding: 8px 14px;
  margin-top: 10px;
  border-radius: 5px;
  cursor: pointer;
}

.secondary-btn:hover {
  background-color: #007bb5;
}

.input-field {
  padding: 8px;
  width: 70%;
  margin-right: 10px;
  border-radius: 5px;
  border: 1px solid #ccc;
}

.message-section {
  display: flex;
  align-items: center;
  margin: 15px 0;
}

.channel-list {
  list-style-type: none;
  padding: 0;
}

.server-btn {
  background-color: #f0f0f0;
  border: 1px solid #ccc;
  padding: 6px 12px;
  border-radius: 5px;
  margin: 5px 0;
  width: 100%;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.3s;
}

.server-btn:hover {
  background-color: #ddd;
}

.invite-link {
  margin-top: 10px;
  word-break: break-all;
}

#chat {
      text-align: left;
      background: #f1f1f1;
      width: 500px;
      min-height: 300px;
      padding: 20px;
    }
</style>
