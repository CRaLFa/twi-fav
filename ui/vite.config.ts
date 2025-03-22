import react from '@vitejs/plugin-react'
import { join } from 'node:path'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@rorygudka/react-infinite-scroller': join(
        import.meta.dirname,
        'node_modules/@rorygudka/react-infinite-scroller/index.js',
      ),
      'react-tweet-theme': join(
        import.meta.dirname,
        'node_modules/react-tweet/dist/twitter-theme',
      ),
    },
  },
  base: '/twi-fav/',
})
