<template>
<button @click="AcceptInvite">Accept The invitation To the server: {{ server_name }}</button>
</template>


<script setup lang="js">
import { ref } from "vue";
import { useRoute } from 'vue-router'
const route = useRoute()

const server_name = ref("")
const invite_code = route.params.invite_code
async function AcceptInvite(){
    try{
        const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/invite/token/${invite_code}`)
        
        if (response.status == 302){
            const data = await response.json();
            window.location.href = `/server/${data.Server_id}/${data.Channel_id}`
        }
        if (response.ok){
            console.log("ZAHEL NA SERVAL")
        }
    
    }catch(err){
        alert(err)
    }
}
</script>