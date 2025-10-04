
<template>
<button class="sign-up-btn" @click="SignUp">Redirect To SignUp</button>
<button class="sign-up-btn" @click="Redirect">Redirect to Test</button>
<button class="sign-up-btn" @click="SignIn">Redirect to SignIn</button>
<button class="sign-up-btn" @click="CreateServer">Redirect to CreateServer</button>
<button type="submit" class="sign-up-btn" @click="GetUserServers">get Servers</button>
<ul v-if="isSessionExists && isServersExists">
    <li v-for="(server, index) in servers" :key="index">
        <button @click="RedirectToServer(server.ServerId)" class="server-btn">{{ server.Servername }}</button>
        
    </li>
</ul>
</template>

<script setup lang="js">
import Session from "supertokens-web-js/recipe/session";
import {ref, onMounted} from 'vue'
import { doesSessionExist, EmailVerificationClaim } from 'supertokens-web-js/recipe/session'

const isSessionExists = ref(false)
const isServersExists = ref(false)
const servers = ref([])

async function doesSessionExistWrapper() {
      if (await doesSessionExist()) {
            isSessionExists.value = true
        } else {
            isSessionExists.value = false
      }
}


async function GetUserServers(){
    const handlerURL = `${import.meta.env.VITE_BACKEND_API_URL}/api/server/getServers`;
    try{
        const response = await fetch(handlerURL,{
            method: "GET",
            headers:{
                "Content-Type": "application/json"
            },
        })
    if (!response.ok) {
      isServersExists.value = false
        return 
    }
    isServersExists.value = true
    const data = await response.json();
    servers.value = data

    }catch (err) {
    alert(err.message);
  }
}


async function SignUp(){
    window.location.href = "/auth/signup"
}
async function Redirect(){
    window.location.href = "/test"
}
async function SignIn(){
    window.location.href = "/auth/signin"
}
async function CreateServer(){
    window.location.href = "/server/create"
}
async function RedirectToServer(ServerId){
    const channel_id = ref('')
    try{
        const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/get-server-channel/${ServerId}`,{
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        })

        if (!response.ok){
            console.log("unexpected error")
            return
        }
        const data = await response.json();
        channel_id.value = data.channel_id
    }catch(err){
        alert(err)
    }
    window.location.href = `/server/${ServerId}/${channel_id.value}`
}

onMounted(async () => {
   await doesSessionExistWrapper()

   if(isSessionExists.value){
    let validateClaims = await Session.validateClaims()
    if (validateClaims.length === 0){
    GetUserServers()
    }else {
          for (const err of validateClaims) {
              if (err.id === EmailVerificationClaim.id) {
                alert("email not verified")
              }
          }
        }
   }
})
</script>

<style scoped>
.blue-block {
      width: 100px;
      height: 100px;
      background-color: blue;
}
.sign-up-btn{
    width: 100px;
    height: 100px;
    background-color: blueviolet;
}
.server-btn {
  background-color: #4f46e5; /* blue */
  color: white;
  font-size: 16px;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.server-btn:hover {
  background-color: #4338ca; /* darker blue */
}
</style>