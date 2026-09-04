import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { createVuetify } from 'vuetify'

export default createVuetify({
  theme: {
    defaultTheme: 'wfrpDark',
    themes: {
      wfrpDark: {
        dark: true,
        colors: {
          primary: '#C99700', // Gold/Brass accent
          secondary: '#7A1C1C', // Dark Crimson / Burgundy
          accent: '#00A896',
          error: '#E63946',
          info: '#457B9D',
          success: '#2A9D8F',
          warning: '#E76F51',
          background: '#121316',
          surface: '#1E2026',
          'surface-variant': '#2A2D36',
          'on-background': '#E2E8F0',
          'on-surface': '#F1F5F9',
        },
      },
      wfrpLight: {
        dark: false,
        colors: {
          primary: '#A37200',
          secondary: '#8C1D1D',
          accent: '#028090',
          background: '#F8FAF9',
          surface: '#FFFFFF',
          'surface-variant': '#E2E8F0',
        },
      },
    },
  },
})
