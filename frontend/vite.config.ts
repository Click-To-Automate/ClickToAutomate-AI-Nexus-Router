import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    // Wails embeds frontend/dist; backend standalone copies from there into public/
    outDir: 'dist',
    emptyOutDir: true,
  }
})
