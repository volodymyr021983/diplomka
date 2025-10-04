<template>
    <form v-if="visible" @submit.prevent="CreateChannel">
            <fieldset class="fieldset" >
            <legend>Channel name:</legend>
            <input class="input" v-model="channelName" placeholder="Channel name" type="text" required><br>
            </fieldset>
    
            <fieldset class="fieldset">
            <legend>Channel type:</legend>
            <label>
            <input v-model="channelType" type="radio" value="Text" />
            Text
            </label><br>

            <label>
            <input v-model="channelType" type="radio" value="Voice"  />
            Voice
            </label>
            </fieldset>
            
            <p v-if="creationError">{{ creationError }}</p>
            <button type="submit" > Create channel2 </button>
    </form>
</template>
<script setup lang = "js">
import { ref } from 'vue'

const channelName = ref('')
const channelType = ref('')
const creationError = ref('')

const props = defineProps({
  visible: Boolean,
  serverId: String,
})


const emits = defineEmits(['close'])
const close = () => emits('close')

async function CreateChannel(){
    creationError.value = ""
    try{
        const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/create-channel/${props.serverId}`,{
            method: "POST",
            headers: {
            "Content-Type": "application/json"
            },
            body: JSON.stringify({
                channelName: channelName.value,
                channelType: channelType.value
            })
        })
        if (response.ok){
            close()
        }
        else{
            creationError.value = "Unexpected error occurs"
        }
    }catch(err){
        creationError.value = "Unexpected error occurs"
    }
}

</script>



<style scoped>
form {
  background-color: #fff;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  width: 320px;
  margin: 60px auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
  font-family: sans-serif;
}

.fieldset {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 12px;
}

legend {
  font-weight: bold;
  color: #444;
}

.input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 6px;
  margin-top: 6px;
  font-size: 14px;
}

button[type="submit"] {
  background-color: #4f46e5;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 16px;
  font-size: 15px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

button[type="submit"]:hover {
  background-color: #4338ca;
}

p {
  color: red;
  font-size: 14px;
  margin: 0;
  text-align: center;
}
</style>
