import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import HomeView from '../HomeView.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'

// Mocks
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
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
})
