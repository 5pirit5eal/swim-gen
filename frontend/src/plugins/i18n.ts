import { createI18n } from 'vue-i18n'
import messages from '@intlify/unplugin-vue-i18n/messages'

const i18n = createI18n({
  legacy: false,
  locale: navigator.language.split('-')[0] || 'en', // Default to browser language or 'en'
  fallbackLocale: 'en',
  messages,
})

export default i18n
