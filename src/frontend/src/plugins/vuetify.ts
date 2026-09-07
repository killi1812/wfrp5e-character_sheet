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
          secondary: '#38BDF8', // Steel Blue (non-red)
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
          'on-error': '#FFFFFF',
        },
      },
      wfrpLight: {
        dark: false,
        colors: {
          primary: '#855B00',
          secondary: '#475569', // Slate Steel (non-red)
          accent: '#028090',
          error: '#D32F2F',
          info: '#1976D2',
          success: '#2E7D32',
          warning: '#ED6C02',
          background: '#F4F4F6',
          surface: '#FFFFFF',
          'surface-variant': '#E2E8F0',
          'on-background': '#0F172A',
          'on-surface': '#1E293B',
          'on-surface-variant': '#475569',
          'on-primary': '#FFFFFF',
          'on-secondary': '#FFFFFF',
          'on-error': '#FFFFFF',
        },
      },
    },
  },
})
