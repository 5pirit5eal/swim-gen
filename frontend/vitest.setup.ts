import { setActivePinia, createPinia } from 'pinia'
import { beforeAll, vi } from 'vitest'
import { createApp } from 'vue' // Import createApp
import i18n from './src/plugins/i18n' // Adjust path as necessary

// Mock the supabase client for all tests
const _globalSupabaseMock = {
  auth: {
    getSession: vi.fn(() => Promise.resolve({ data: { session: null } })),
    getUser: vi.fn(() => Promise.resolve({ data: { user: null } })),
    onAuthStateChange: vi.fn(() => ({
      data: { subscription: { unsubscribe: vi.fn() } },
    })),
    signInWithPassword: vi.fn(),
    signUp: vi.fn(),
    signOut: vi.fn(),
    refreshSession: vi.fn(() => Promise.resolve({ data: { session: null } })),
  },
}

vi.mock('@/plugins/supabase', () => ({
  supabase: _globalSupabaseMock,
  getSupabase: vi.fn(async () => _globalSupabaseMock),
}))

beforeAll(() => {
  setActivePinia(createPinia())

  // Create a mock app and install vue-i18n
  const app = createApp({})
  app.use(i18n)
})
