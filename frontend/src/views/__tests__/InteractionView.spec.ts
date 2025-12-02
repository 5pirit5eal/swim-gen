import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import InteractionView from '../InteractionView.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useRoute } from 'vue-router'

// Mock vue-i18n
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}))

// Mock dependencies
vi.mock('vue-router', () => ({
  useRoute: vi.fn(),
  useRouter: vi.fn(() => ({
    push: vi.fn(),
  })),
}))

vi.mock('@/components/training/TrainingPlanDisplay.vue', () => ({
  default: { template: '<div class="plan-display-stub"></div>' },
}))

vi.mock('@/components/training/SimplePlanDisplay.vue', () => ({
  default: { template: '<div class="simple-plan-display-stub"></div>' },
}))

vi.mock('@/components/forms/FeedbackForm.vue', () => ({
  default: { template: '<div class="feedback-form-stub"></div>', props: ['show'] },
}))

describe('InteractionView.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    ;(useRoute as unknown as ReturnType<typeof vi.fn>).mockReturnValue({
      params: { id: 'test-plan-id' },
    })
    window.scrollTo = vi.fn()
  })

  it('renders correctly', async () => {
    const wrapper = mount(InteractionView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              trainingPlan: {
                currentPlan: { title: 'Test Plan', table: [] },
                conversation: [],
              },
            },
          }),
        ],
        stubs: {
          IconSend: true,
          IconShare: true,
          IconDownload: true,
          IconStar: true,
        },
      },
    })

    await flushPromises()
    expect(wrapper.find('.interaction-view').exists()).toBe(true)
    expect(wrapper.findAll('.tab-button').length).toBe(2)
  })

  it('switches tabs', async () => {
    const wrapper = mount(InteractionView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              trainingPlan: {
                currentPlan: { title: 'Test', table: [] },
              },
            },
          }),
        ],
        stubs: {
          IconSend: true,
          IconShare: true,
          IconDownload: true,
          IconStar: true,
        },
      },
    })

    const tabs = wrapper.findAll('.tab-button')
    expect(tabs[0].classes()).toContain('active') // Plan tab default

    await tabs[1].trigger('click') // Chat tab
    expect(tabs[1].classes()).toContain('active')
    expect(wrapper.find('.chat-container').isVisible()).toBe(true)
  })

  it('sends a message', async () => {
    const wrapper = mount(InteractionView, {
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
            initialState: {
              trainingPlan: {
                currentPlan: { title: 'Test', table: [] },
                conversation: [],
              },
            },
          }),
        ],
        stubs: {
          IconSend: true,
          IconShare: true,
          IconDownload: true,
          IconStar: true,
        },
      },
    })

    const store = useTrainingPlanStore()
    store.sendMessage = vi.fn().mockResolvedValue(true)

    // Switch to chat tab
    await wrapper.findAll('.tab-button')[1].trigger('click')

    const input = wrapper.find('.chat-input')
    await input.setValue('Make it harder')
    await wrapper.find('.chat-form').trigger('submit')

    expect(store.sendMessage).toHaveBeenCalledWith('Make it harder')
    expect((input.element as HTMLInputElement).value).toBe('')
  })
})
