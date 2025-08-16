// frontend/vitest.setup.ts
import { setActivePinia, createPinia } from 'pinia'
import { beforeAll } from 'vitest'

beforeAll(() => {
  setActivePinia(createPinia())
})
