<template>
<button class="secondary-btn" @click="GetInviteLink">Get Invite Link</button>
    <h1 class="invite-link">Invite link: {{ link }}</h1>
    <p v-if="linkError">{{ linkError }}</p>
</template>


<script setup lang="js">
import { ref } from 'vue'

const link = ref('')
const linkError = ref('')

const props = defineProps({
    serverId: String,
})
async function GetInviteLink() {
    linkError.value = ""
  try {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/invite/${props.serverId}`, {
      method: "GET",
      headers: { "Content-Type": "application/json" }
    })
    if(response.ok){
      const data = await response.json();
      link.value = `${import.meta.env.VITE_FRONTEND_API_URL}/server/invite/${data.Invite_code}`
    }
  } catch(err) {
    alert(err)
    linkError.value = "Unexpected error occurs"
  }
}

</script>
