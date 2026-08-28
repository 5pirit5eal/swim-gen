import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import HomeView from '../HomeView.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'

// Mocks
const pushMock = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: pushMock }),
}))
// ...
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
    locale: { value: 'en' }, // Mock ref-like object
  }),
}))
vi.mock('@/tutorial/useTutorial', () => ({
  useTutorial: () => ({ startHomeTutorial: vi.fn() }),
}))

vi.mock('@/plugins/supabase', () => ({
  supabase: {
    auth: {
      onAuthStateChange: vi.fn(),
    },
    from: vi.fn().mockReturnThis(),
    select: vi.fn().mockReturnThis(),
    order: vi.fn().mockReturnThis(),
    range: vi.fn().mockReturnThis(),
    in: vi.fn().mockReturnThis(),
    eq: vi.fn().mockReturnThis(),
    single: vi.fn().mockReturnThis(),
    maybeSingle: vi.fn().mockReturnThis(),
  },
  getSupabase: vi.fn(async () => ({
    auth: {
      onAuthStateChange: vi.fn(),
    },
    from: vi.fn().mockReturnThis(),
    select: vi.fn().mockReturnThis(),
    order: vi.fn().mockReturnThis(),
    range: vi.fn().mockReturnThis(),
    in: vi.fn().mockReturnThis(),
    eq: vi.fn().mockReturnThis(),
    single: vi.fn().mockReturnThis(),
    maybeSingle: vi.fn().mockReturnThis(),
  })),
}))

describe('HomeView.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('restores anonymous plan from localStorage on mount', () => {
    const plan = { title: 'Restored Plan', table: [] }
    const query = 'restored query'
    localStorage.setItem('anonymousPlan', JSON.stringify(plan))
    localStorage.setItem('anonymousQuery', query)

    mount(HomeView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              trainingPlan: {
                currentPlan: null,
                initialQuery: '',
              },
            },
          }),
        ],
        stubs: {
          TrainingPlanForm: true,
          TrainingPlanDisplay: true,
        },
      },
    })

    const store = useTrainingPlanStore()
    expect(store.currentPlan).toEqual(plan)
    expect(store.initialQuery).toBe(query)
    expect(localStorage.getItem('anonymousPlan')).toBeNull()
  })

  it('links anonymous plan when user is logged in', async () => {
    mount(HomeView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: { user: { id: 'user-1' } },
              trainingPlan: {
                currentPlan: { plan_id: undefined, title: 'Anon' },
                initialQuery: 'query',
              },
            },
          }),
        ],
        stubs: {
          TrainingPlanForm: true,
          TrainingPlanDisplay: true,
        },
      },
    })

    const store = useTrainingPlanStore()
    // The watcher should trigger immediately
    expect(store.linkAnonymousPlan).toHaveBeenCalled()
  })

  it('renders CTA banner for unauthenticated user with current plan and handles navigation', async () => {
    const wrapper = mount(HomeView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: { user: null },
              trainingPlan: {
                currentPlan: { plan_id: undefined, title: 'My Workout', table: [] },
                initialQuery: 'Swim 2000m',
              },
            },
          }),
        ],
        stubs: {
          TrainingPlanForm: true,
          TrainingPlanDisplay: true,
        },
      },
    })

    expect(wrapper.find('.cta-banner').exists()).toBe(true)
    expect(wrapper.find('.cta-title').text()).toBe('home.banner.not_logged_in.title')

    // Click primary button (Save my plan -> register)
    await wrapper.find('.cta-button').trigger('click')
    expect(pushMock).toHaveBeenCalledWith({ name: 'login', query: { register: 'true' } })

    // Click secondary link (Already registered -> login)
    await wrapper.find('.cta-secondary-link').trigger('click')
    expect(pushMock).toHaveBeenCalledWith({ name: 'login' })
  })

  it('renders CTA banner for authenticated user with current plan and navigates to interaction', async () => {
    const wrapper = mount(HomeView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              auth: { user: { id: 'user-123' } },
              trainingPlan: {
                currentPlan: { plan_id: 'plan-123', title: 'My Workout', table: [] },
                initialQuery: 'Swim 2000m',
              },
            },
          }),
        ],
        stubs: {
          TrainingPlanForm: true,
          TrainingPlanDisplay: true,
        },
      },
    })

    expect(wrapper.find('.cta-banner').exists()).toBe(true)
    expect(wrapper.find('.cta-title').text()).toBe('home.banner.logged_in.title')

    await wrapper.find('.cta-button').trigger('click')
    expect(pushMock).toHaveBeenCalledWith({ name: 'plan', params: { id: 'plan-123' } })
  })
})
