/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  variants: {
    extend: {
      tableLayout: ['hover', 'focus'],
    },
    container: {
      center: true,
    },
  },
  theme: {
    extend: {},
  },
  plugins: [],
}
