import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import UploadForm from '../UploadForm.vue'
import { useUploadFormStore } from '@/stores/uploadForm'
import { nextTick } from 'vue'

// Mock vue-i18n
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}))

// Mock dependencies
vi.mock('@/api/client', () => ({
  apiClient: {
    fileToPlan: vi.fn(),
  },
}))

vi.mock('@/components/training/TrainingPlanDisplay.vue', () => ({
  default: { template: '<div class="training-plan-display-stub"></div>' },
}))

vi.mock('@/components/ui/BaseModal.vue', () => ({
  default: {
    template: `
      <div v-if="show" class="base-modal-stub">
        <slot name="header"></slot>
        <slot name="body"></slot>
        <slot name="footer"></slot>
      </div>
    `,
    props: ['show'],
    emits: ['close'],
  },
}))

describe('UploadForm.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders correctly when shown', () => {
    const wrapper = mount(UploadForm, {
      props: { show: true },
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              uploadForm: {
                currentPlan: { title: 'Test Plan', table: [] },
                isLoading: false,
              },
            },
          }),
        ],
        stubs: {
          IconUpload: true,
        },
      },
    })

    expect(wrapper.find('.base-modal-stub').exists()).toBe(true)
    expect(wrapper.find('h2').text()).toBe('donation.title')
    expect(wrapper.find('.upload-image-btn').exists()).toBe(true)
  })

  it('initializes new plan on mount', () => {
    mount(UploadForm, {
      props: { show: true },
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: false,
          }),
        ],
      },
    })

    const store = useUploadFormStore()
    expect(store.initNewPlan).toHaveBeenCalled()
  })

  it('validates title before submission', async () => {
    const wrapper = mount(UploadForm, {
      props: { show: true },
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            initialState: {
              uploadForm: {
                currentPlan: { title: '', table: [] }, // Empty title
              },
            },
          }),
        ],
      },
    })

    const store = useUploadFormStore()
    await wrapper.find('.submit-btn').trigger('click')

    expect(store.uploadCurrentPlan).not.toHaveBeenCalled()
    // Ideally check for toast error here, but toast is hard to test in isolation without a real DOM
  })

  it('submits successfully when valid', async () => {
    const wrapper = mount(UploadForm, {
      props: { show: true },
      global: {
        plugins: [
          createTestingPinia({
            createSpy: vi.fn,
            stubActions: false, // We want to mock the return value, but let the action be called
          }),
        ],
      },
    })

    const store = useUploadFormStore()
    // Mock the store action implementation or return value
    store.uploadCurrentPlan = vi.fn().mockResolvedValue(true)
    store.currentPlan = {
      title: 'Valid Title',
      table: [],
    } as unknown as import('@/types').RAGResponse
    store.currentPlan = {
      title: 'Valid Title',
      table: [],
      description: '',
      plan_id: '1',
    }

    await wrapper.find('.submit-btn').trigger('click')
    await nextTick()

    expect(store.uploadCurrentPlan).toHaveBeenCalled()
    expect(store.clear).toHaveBeenCalled()
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
