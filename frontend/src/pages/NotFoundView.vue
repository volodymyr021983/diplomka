<template>
  <div class="not-found-container">
    <h1 class="not-found-text">404 - Not Found</h1>
    <p class="not-found-subtext">The page you are looking for does not exist.</p>
  </div>
</template>

<script setup>
</script>

<style scoped>
.not-found-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 80vh;
  text-align: center;
  background-color: #f8f8f8;
}

.not-found-text {
  font-size: 6rem;
  font-weight: bold;
  color: #ff4d4f;
  margin: 0;
}

.not-found-subtext {
  font-size: 1.5rem;
  color: #555;
  margin-top: 1rem;
}
</style>
