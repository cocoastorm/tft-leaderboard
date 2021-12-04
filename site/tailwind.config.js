const colors = require('tailwindcss/colors')

module.exports = {
  purge: ['./pages/**/*.{js,ts,jsx,tsx}', './components/**/*.{js,ts,jsx,tsx}'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        sky: colors.sky,
        emerald: colors.emerald,
        teal: colors.teal,
        fuchsia: colors.fuchsia,
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
