import { describe, it, expect, vi, beforeEach, type Mock } from 'vitest'
import { mount, RouterLinkStub } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { toast } from 'vue3-toastify'
import LoginView from '../LoginView.vue'
import { useAuthStore } from '@/stores/auth'

const { mockPush, mockRoute } = vi.hoisted(() => ({
  mockPush: vi.fn(),
  mockRoute: { query: {} as Record<string, string> },
}))

// Mock dependencies
vi.mock('@/plugins/supabase', () => ({
  supabase: {
    auth: {
      getSession: vi.fn(() => Promise.resolve({ data: { session: null }, error: null })),
      getUser: vi.fn(() => Promise.resolve({ data: { user: null }, error: null })),
      onAuthStateChange: vi.fn(),
      signInWithPassword: vi.fn(),
      signUp: vi.fn(),
      signOut: vi.fn(),
    },
    from: vi.fn().mockReturnThis(),
    select: vi.fn().mockReturnThis(),
    eq: vi.fn().mockReturnThis(),
    order: vi.fn().mockReturnThis(),
    in: vi.fn().mockReturnThis(),
    insert: vi.fn().mockReturnThis(),
    limit: vi.fn().mockResolvedValue({ data: [], error: null }),
    single: vi.fn().mockResolvedValue({ data: null, error: null }),
  },
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
  useRoute: () => mockRoute,
}))
vi.mock('@/router', () => ({
  default: {
    push: mockPush,
    replace: mockPush,
  },
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('vue3-toastify', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

describe('LoginView.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockRoute.query = {}
  })

  it('renders the login form by default', () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
    expect(wrapper.find('h1').text()).toBe('login.login')
    expect(wrapper.find('input#username').exists()).toBe(false)
  })

  it('renders the sign-up form when register query is true', () => {
    mockRoute.query = { register: 'true' }
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
    expect(wrapper.find('h1').text()).toBe('login.signUp')
    expect(wrapper.find('input#username').exists()).toBe(true)
  })

  it('disables the submit button until all fields are filled', async () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })

    const button = wrapper.find('button[type="submit"]')
    expect(button.attributes('disabled')).toBeDefined()

    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password')

    expect(button.attributes('disabled')).toBeUndefined()
  })

  it('calls signInWithPassword on login form submission', async () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
    const auth = useAuthStore()

    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password')
    await wrapper.find('form').trigger('submit.prevent')

    expect(auth.signInWithPassword).toHaveBeenCalledWith('test@example.com', 'password')
    expect(toast.success).toHaveBeenCalledWith('login.loginSuccess')
    expect(mockPush).toHaveBeenCalledWith('/')
  })

  it('handles login failure', async () => {
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
    const auth = useAuthStore()
      ; (auth.signInWithPassword as Mock).mockRejectedValue(new Error('Invalid login credentials'))

    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password')
    await wrapper.find('form').trigger('submit.prevent')

    expect(toast.error).toHaveBeenCalledWith('login.invalidLogin')
  })

  it('calls signUp on sign-up form submission', async () => {
    mockRoute.query = { register: 'true' }
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
    const auth = useAuthStore()
      ; (auth.signUp as Mock).mockResolvedValue({ user: { identities: [{}] } })

    await wrapper.find('input#username').setValue('newuser')
    await wrapper.find('input#email').setValue('new@example.com')
    await wrapper.find('input#password').setValue('newpassword')
    await wrapper.find('form').trigger('submit.prevent')

    expect(auth.signUp).toHaveBeenCalledWith('new@example.com', 'newpassword', 'newuser')
    expect(toast.success).toHaveBeenCalledWith('login.registrationSuccess')
    expect(mockPush).toHaveBeenCalledWith('/login')
  })

  it('handles sign-up failure when user exists', async () => {
    mockRoute.query = { register: 'true' }
    const wrapper = mount(LoginView, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    })
    const auth = useAuthStore()
      ; (auth.signUp as Mock).mockResolvedValue({ user: { identities: [] } })
      ; (auth.signInWithPassword as Mock).mockResolvedValue(undefined)

    await wrapper.find('input#username').setValue('existinguser')
    await wrapper.find('input#email').setValue('existing@example.com')
    await wrapper.find('input#password').setValue('password')
    await wrapper.find('form').trigger('submit.prevent')

    // Wait for promises
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(auth.signInWithPassword).toHaveBeenCalledWith('existing@example.com', 'password')
    expect(toast.success).toHaveBeenCalledWith('login.userExistsLoginSuccess')
    expect(mockPush).toHaveBeenCalledWith('/')
  })
})
