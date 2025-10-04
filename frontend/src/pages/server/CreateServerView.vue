<template>
  <div class="create-server">
    <h1>Create a new server</h1>

    <form @submit.prevent="createServer">
      <div>
        <label for="name">Server Name</label>
        <input
          id="name"
          v-model="name"
          type="text"
          placeholder="My Server"
          required
        />
      </div>
      <button type="submit">Create</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const name = ref("");


async function createServer() {
  try {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_API_URL}/api/server/create`,{
      method: "POST",
      body: JSON.stringify({
        Servername: name.value
      }),  
    })

    if (!response.ok) {
      throw new Error("Failed to create server");
    }

    const data = await response.json();
    console.log(data.body.OK)
    // редирект на страницу сервера
  } catch (err) {
    alert(err.message);
  }
}
</script>

<style scoped>
.create-server {
  max-width: 400px;
  margin: 2rem auto;
  padding: 1rem;
  border: 1px solid #ccc;
  border-radius: 8px;
}

form > div {
  margin-bottom: 1rem;
}

label {
  display: block;
  font-weight: bold;
  margin-bottom: 0.3rem;
}

input,
textarea {
  width: 100%;
  padding: 0.4rem;
}
</style>
