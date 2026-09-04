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
          primary: '#D4AF37', // Metallic Gold
          secondary: '#9E2A2B', // Dark Crimson
          accent: '#00A896',
          error: '#E63946',
          info: '#457B9D',
          success: '#2A9D8F',
          warning: '#E76F51',
          background: '#121316',
          surface: '#1E2026',
          'surface-variant': '#2A2D36',
          'on-background': '#F1F5F9',
          'on-surface': '#F8FAFC',
          'on-surface-variant': '#CBD5E1',
          'on-primary': '#121316',
          'on-secondary': '#FFFFFF',
        },
      },
      wfrpLight: {
        dark: false,
        colors: {
          primary: '#9A6B00',
          secondary: '#8C1D1D',
          accent: '#028090',
          error: '#D32F2F',
          info: '#1976D2',
          success: '#2E7D32',
          warning: '#ED6C02',
          background: '#F4F4F6',
          surface: '#FFFFFF',
          'surface-variant': '#E8ECEF',
          'on-background': '#1E293B',
          'on-surface': '#0F172A',
          'on-surface-variant': '#334155',
          'on-primary': '#FFFFFF',
          'on-secondary': '#FFFFFF',
        },
      },
    },
  },
})
