module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    fontFamily: {
      'sans': ['Roboto', 'ui-sans-serif', 'system-ui'],
      'serif': ['serif']
    },
    extend: {
      typography: {
        DEFAULT: {
          css: {
            lineHeight: '1.35rem'
          }
        }
      }
    }
  },
  plugins: [
    require('@tailwindcss/typography')
  ],
}
