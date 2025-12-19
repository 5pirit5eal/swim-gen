import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ButtonSharePlan from '../ButtonSharePlan.vue'
import { createTestingPinia } from '@pinia/testing'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import i18n from '@/plugins/i18n'
import type { RAGResponse } from '@/types'

// Mock toast
const { mockToastSuccess, mockToastError } = vi.hoisted(() => ({
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
}))

vi.mock('vue3-toastify', () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}))

// Mock apiClient
const { mockCreateShareUrl } = vi.hoisted(() => ({
  mockCreateShareUrl: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  apiClient: {
    createShareUrl: mockCreateShareUrl,
  },
}))

// Mock platform detection
const { getMockIsIOS, setMockIsIOS } = vi.hoisted(() => {
  let mockIsIOS = false
  return {
    getMockIsIOS: () => mockIsIOS,
    setMockIsIOS: (value: boolean) => {
      mockIsIOS = value
    },
  }
})

vi.mock('@/utils/platform', () => ({
  isIOS: () => getMockIsIOS(),
}))

describe('ButtonSharePlan.vue', () => {
  let mockStore: ReturnType<typeof useTrainingPlanStore>
  let mockClipboard: { writeText: ReturnType<typeof vi.fn> }

  beforeEach(() => {
    vi.clearAllMocks()
    setMockIsIOS(false)

    // Mock clipboard API
    mockClipboard = {
      writeText: vi.fn().mockResolvedValue(undefined),
    }
    Object.defineProperty(navigator, 'clipboard', {
      value: mockClipboard,
      configurable: true,
    })

    // Setup pinia
    const pinia = createTestingPinia({
      createSpy: vi.fn,
    })

    mockStore = useTrainingPlanStore(pinia)
    mockStore.currentPlan = {
      plan_id: 'test-plan-id',
      title: 'Test Plan',
      description: 'Test Description',
      table: [],
    } as RAGResponse
    mockStore.keepForever = vi.fn().mockResolvedValue(undefined)

    // Mock API response
    mockCreateShareUrl.mockResolvedValue({
      success: true,
      data: { url_hash: 'abc123' },
    })
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('Non-iOS behavior (one-step share)', () => {
    beforeEach(() => {
      setMockIsIOS(false)
    })

    it('renders share button with correct initial state', () => {
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      expect(wrapper.find('.share-btn').exists()).toBe(true)
      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.share_plan'))
    })

    it('creates share URL and copies immediately on click', async () => {
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      await wrapper.find('.share-btn').trigger('click')

      // Should call keepForever
      expect(mockStore.keepForever).toHaveBeenCalledWith('test-plan-id')

      // Should create share URL
      expect(mockCreateShareUrl).toHaveBeenCalledWith({
        plan_id: 'test-plan-id',
        method: 'link',
      })

      await wrapper.vm.$nextTick()

      // Should copy to clipboard immediately
      expect(mockClipboard.writeText).toHaveBeenCalledWith(
        `${window.location.origin}/shared/abc123`,
      )

      // Should show success toast
      expect(mockToastSuccess).toHaveBeenCalledWith(i18n.global.t('share.copied'))
    })

    it('shows success state after copying', async () => {
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      await wrapper.find('.share-btn').trigger('click')
      await flushPromises()

      // Should show "Copied" text
      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.copied'))
      expect(wrapper.find('.share-btn').classes()).toContain('success')
    })

    it('resets to initial state after 2 seconds', async () => {
      vi.useFakeTimers()
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      await wrapper.find('.share-btn').trigger('click')
      await flushPromises()

      // After 2 seconds, should reset
      await vi.advanceTimersByTimeAsync(2000)
      await flushPromises()

      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.share_plan'))
    })

    it('handles API error gracefully', async () => {
      mockCreateShareUrl.mockResolvedValueOnce({
        success: false,
        error: { message: 'API Error' },
      })

      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      expect(mockToastError).toHaveBeenCalledWith(i18n.global.t('share.create_error'))
      expect(mockClipboard.writeText).not.toHaveBeenCalled()
    })

    it('handles clipboard error gracefully', async () => {
      mockClipboard.writeText.mockRejectedValueOnce(new Error('Clipboard error'))

      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      expect(mockToastError).toHaveBeenCalledWith(i18n.global.t('share.copy_error'))
    })
  })

  describe('iOS behavior (two-step share)', () => {
    beforeEach(() => {
      setMockIsIOS(true)
    })

    it('first click creates share URL and shows copy button', async () => {
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      // Should create share URL
      expect(mockCreateShareUrl).toHaveBeenCalledWith({
        plan_id: 'test-plan-id',
        method: 'link',
      })

      // Should NOT copy to clipboard yet
      expect(mockClipboard.writeText).not.toHaveBeenCalled()

      // Should show "Copy" text
      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.copy'))
    })

    it('second click copies URL to clipboard', async () => {
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      // First click
      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      expect(mockClipboard.writeText).not.toHaveBeenCalled()

      // Second click
      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      // Should copy to clipboard
      expect(mockClipboard.writeText).toHaveBeenCalledWith(
        `${window.location.origin}/shared/abc123`,
      )

      // Should show success toast
      expect(mockToastSuccess).toHaveBeenCalledWith(i18n.global.t('share.copied'))

      // Should show "Copied" text
      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.copied'))
    })

    it('resets after 2 seconds on iOS', async () => {
      vi.useFakeTimers()
      const wrapper = mount(ButtonSharePlan, {
        props: { store: mockStore },
        global: {
          plugins: [i18n],
        },
      })

      // First click - create URL
      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      // Second click - copy URL
      await wrapper.find('.share-btn').trigger('click')
      await wrapper.vm.$nextTick()

      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.copied'))

      // After 2 seconds, should reset
      vi.advanceTimersByTime(2000)
      await wrapper.vm.$nextTick()

      expect(wrapper.find('.share-btn').text()).toContain(i18n.global.t('share.share_plan'))
    })
  })

  it('disables button while loading', async () => {
    const wrapper = mount(ButtonSharePlan, {
      props: { store: mockStore },
      global: {
        plugins: [i18n],
      },
    })

    // Make API slow
    mockCreateShareUrl.mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve({ success: true }), 1000)),
    )

    const button = wrapper.find('.share-btn')
    await button.trigger('click')

    expect(button.attributes('disabled')).toBeDefined()
  })

  it('does nothing if no current plan', async () => {
    mockStore.currentPlan = null

    const wrapper = mount(ButtonSharePlan, {
      props: { store: mockStore },
      global: {
        plugins: [i18n],
      },
    })

    await wrapper.find('.share-btn').trigger('click')

    expect(mockStore.keepForever).not.toHaveBeenCalled()
    expect(mockCreateShareUrl).not.toHaveBeenCalled()
  })
})
