import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(),tailwindcss(),],
  server: {
    proxy: {
      '/manticore':{
        target: 'http://127.0.0.1:9308',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/manticore/, '')
      },
      '/purr': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/purr/, '')
      }
    }
  }
})
