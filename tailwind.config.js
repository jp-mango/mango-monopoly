/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./ui/**/*.tmpl', './cmd/web/**/*.go'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['JetBrains Mono', 'monospace'],
      },
    },
  },
  plugins: [],
};
