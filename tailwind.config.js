/** @type {import('tailwindcss').Config} */
const theme = require('./tailwind.theme.json');
module.exports = {
  content: ["./input.css", "./content/*.md", "./content/**/*.md", "./content/**/**/*.md", "./layouts/**/*.html"],
  theme: theme,
  plugins: [
    require('@tailwindcss/typography'),
    require("@tailwindcss/forms"),
  ],
}

