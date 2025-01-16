/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,ts}'],
  darkMode: 'selector',
  theme: {
    extend: {
      width: {
        'cx-75': '50%'
      }
    }
  },
  plugins: [require('tailwindcss-primeui')]
};
