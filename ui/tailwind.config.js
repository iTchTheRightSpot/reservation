/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,ts}'],
  darkMode: 'selector',
  theme: {
    extend: {
      width: {
        'cx-75': '50%'
      },
      colors: {
        dim: 'rgba(255, 255, 255, 0.3)'
      }
    }
  },
  plugins: [require('tailwindcss-primeui')]
};
