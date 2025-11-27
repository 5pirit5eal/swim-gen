import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Sidebar from '../AppSidebar.vue'
import { createTestingPinia } from '@pinia/testing'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useSidebarStore } from '@/stores/sidebar'
import i18n from '@/plugins/i18n'

const push = vi.fn()
vi.mock('vue-router', async (importOriginal) => {
  const actual = (await importOriginal()) as typeof import('vue-router')
  return {
    ...actual,
    useRouter: vi.fn(() => ({
      push,
      currentRoute: { value: { path: '/not-home' } },
    })),
  }
})

describe('AppSidebar.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    const pinia = createTestingPinia({
      createSpy: vi.fn,
    })

    const trainingPlanStore = useTrainingPlanStore(pinia)
    trainingPlanStore.generationHistory = [
      { plan_id: '1', title: 'Test Plan 1', description: 'Desc 1', table: [] },
    ]
    trainingPlanStore.historyMetadata = [
      {
        plan_id: '1',
        keep_forever: false,
        created_at: '2025-01-01T00:00:00Z',
        updated_at: '2025-01-01T00:00:00Z',
      },
    ]

    const sidebarStore = useSidebarStore(pinia)
    sidebarStore.isOpen = true
  })

  it('renders the sidebar when open', () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    expect(wrapper.find('.sidebar').classes()).toContain('is-open')
    expect(wrapper.find('.sidebar-header h3').text()).toBe(i18n.global.t('sidebar.history'))
    expect(wrapper.find('.plan-title span').text()).toBe('Test Plan 1')
  })

  it('closes the sidebar when the close button is clicked', async () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    const sidebarStore = useSidebarStore()
    await wrapper.find('.close-btn').trigger('click')
    expect(sidebarStore.close).toHaveBeenCalledTimes(1)
  })

  it('loads a plan when a plan is clicked', async () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    const trainingPlanStore = useTrainingPlanStore()

    await wrapper.find('.plan-title').trigger('click')
    expect(trainingPlanStore.loadPlanFromHistory).toHaveBeenCalledTimes(1)
    expect(push).toHaveBeenCalledWith('/plan/1')
  })
})
