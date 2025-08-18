import { createI18n } from 'vue-i18n'
import messages from '@intlify/unplugin-vue-i18n/messages'

console.log('i18n messages:', messages)

const i18n = createI18n({
  locale: 'de',
  fallbackLocale: 'en',
  messages: messages,
})

export default i18n
