import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api/admin': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        rewrite: (path) => path
      },
      '/api/product': {
        target: 'http://localhost:8001',
        changeOrigin: true,
        rewrite: (path) => path
      },
      '/api/order': {
        target: 'http://localhost:8002',
        changeOrigin: true,
        rewrite: (path) => path
      },
      '/api/payment': {
        target: 'http://localhost:8003',
        changeOrigin: true,
        rewrite: (path) => path
      },
      '/api/film': {
        target: 'http://localhost:8007',
        changeOrigin: true,
        rewrite: (path) => path
      },
      '/api/upload': {
        target: 'http://localhost:8001',
        changeOrigin: true,
        rewrite: (path) => path
      },
      '/uploads': {
        target: 'http://localhost:8001',
        changeOrigin: true
      }
    }
  }
})
