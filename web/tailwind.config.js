module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    fontFamily: {
      'sans': ['Open Sans', 'ui-sans-serif', 'system-ui'],
      'serif': ['Lora', 'serif'],
      'yomogi': ['Yomogi', 'cursive']
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
