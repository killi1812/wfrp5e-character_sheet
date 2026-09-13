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
          secondary: '#4A90E2', // Imperial / Royal Blue
          accent: '#00A896',
          error: '#E63946',
          info: '#457B9D',
          success: '#2A9D8F',
          warning: '#EAB308', // Warm Golden Yellow
          background: '#121316',
          surface: '#1E2026',
          'surface-variant': '#2A2D36',
          'on-background': '#F1F5F9',
          'on-surface': '#F8FAFC',
          'on-surface-variant': '#CBD5E1',
          'on-primary': '#121316',
          'on-secondary': '#FFFFFF',
          'on-error': '#FFFFFF',
          'on-warning': '#121316',
        },
      },
      wfrpLight: {
        dark: false,
        colors: {
          primary: '#855B00',
          secondary: '#2563EB', // Royal Blue
          accent: '#028090',
          error: '#D32F2F',
          info: '#1976D2',
          success: '#2E7D32',
          warning: '#CA8A04', // Warm Amber Gold
          background: '#F4F4F6',
          surface: '#FFFFFF',
          'surface-variant': '#E2E8F0',
          'on-background': '#0F172A',
          'on-surface': '#1E293B',
          'on-surface-variant': '#475569',
          'on-primary': '#FFFFFF',
          'on-secondary': '#FFFFFF',
          'on-error': '#FFFFFF',
          'on-warning': '#FFFFFF',
        },
      },
    },
  },
})
