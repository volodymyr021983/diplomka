import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import fs from 'fs' // Import the Node.js file system module
//import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  
    server: {
    host: '0.0.0.0',
    port: 5173,
    https: {
      key: fs.readFileSync('./certificate/localhost+3-key.pem'),
      cert: fs.readFileSync('./certificate/localhost+3.pem'),
    },
    },
  
    plugins: [
    vue(),
    //vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
})
