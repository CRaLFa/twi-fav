import react from '@vitejs/plugin-react'
import path from 'node:path'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@rorygudka/react-infinite-scroller': path.join(
        import.meta.dirname,
        'node_modules/@rorygudka/react-infinite-scroller/index.js',
      ),
      'react-tweet-theme': path.join(
        import.meta.dirname,
        'node_modules/react-tweet/dist/twitter-theme',
      ),
    },
  },
  base: '/twi-fav/',
})
