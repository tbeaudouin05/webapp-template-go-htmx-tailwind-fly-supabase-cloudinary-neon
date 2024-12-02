module.exports = {
  content: [
    './**/*.templ', // Include your Templ components
    './**/*.go',    // Include Go files if necessary
    './node_modules/flowbite/**/*.js' // Include Flowbite components
  ],
  theme: {
    extend: {},

  },
  plugins: [
    require('flowbite/plugin'),
  ],
}
