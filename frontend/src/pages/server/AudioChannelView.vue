<template>

<video ref="clientVideoElement" autoplay controls width="320" height="240"></video>

<button @click="ConnectToVoice"></button>
</template>

<script setup lang="js">
import {ref} from 'vue'
const clientVideoElement = ref(null)
const localClientStream = ref(null)


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
    }catch(err){
        alert("Oopssomthing wrong")
    }
}
</script>

<style scoped>
video{
border:2px solid blue;
padding:5px;
/**Or add your own style**/
}
</style>