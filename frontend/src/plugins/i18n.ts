import { createI18n } from 'vue-i18n'
import en from '../locales/en.json'
import de from '../locales/de.json'

const i18n = createI18n({
  legacy: false,
  locale: 'de', // set locale
  fallbackLocale: 'en', // set fallback locale
  messages: {
    en,
    de
  }
})

export default i18n
