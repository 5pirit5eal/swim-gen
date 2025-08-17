import { setActivePinia, createPinia } from 'pinia'
import { beforeAll } from 'vitest'
import { createApp } from 'vue' // Import createApp
import i18n from './src/plugins/i18n' // Adjust path as necessary

beforeAll(() => {
  setActivePinia(createPinia())

  // Create a mock app and install vue-i18n
  const app = createApp({})
  app.use(i18n)
})
