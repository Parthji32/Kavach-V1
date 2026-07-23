/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.html",
    "./internal/handlers/**/*.go",
    "./static/js/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        kavach: {
          dark: '#0A0A14',
          darker: '#0D0B1A',
          surface: '#12101F',
          accent: '#7C3AED',
          'accent-hover': '#8B5CF6',
          cyan: '#06B6D4',
          border: '#1E1A30',
        }
      }
    }
  },
  plugins: [],
}
