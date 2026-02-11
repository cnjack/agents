/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'pixel-green': '#4a7c23',
        'pixel-brown': '#6b4423',
        'pixel-blue': '#4a90c2',
        'pixel-gold': '#d4a017',
        'pixel-red': '#c44',
      },
      fontFamily: {
        pixel: ['"Press Start 2P"', 'monospace'],
      },
    },
  },
  plugins: [],
}
