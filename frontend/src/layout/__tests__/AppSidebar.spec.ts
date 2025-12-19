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
    useRoute: vi.fn(() => ({
      name: 'plan',
      params: { id: '1' },
    })),
  }
})

// Mock apiClient
const { mockCreateShareUrl } = vi.hoisted(() => ({
  mockCreateShareUrl: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  apiClient: {
    deletePlan: vi.fn().mockResolvedValue({ success: true }),
    upsertPlan: vi.fn().mockResolvedValue({ success: true }),
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

describe('AppSidebar.vue', () => {
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

    // Mock API response for share
    mockCreateShareUrl.mockResolvedValue({
      success: true,
      data: { url_hash: 'test-hash' },
    })
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
  it('renders the create new and upload buttons with correct text', () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    const buttons = wrapper.findAll('.create-new-btn')
    expect(buttons[0]!.text()).toBe(i18n.global.t('sidebar.create_new'))
    expect(buttons[1]!.text()).toBe(i18n.global.t('sidebar.upload_plan'))
  })

  it('toggles the menu when dots icon is clicked', async () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    // Menu should be hidden initially
    expect(wrapper.find('.dropdown-menu').exists()).toBe(false)

    // Click dots button
    await wrapper.find('.menu-button').trigger('click')
    expect(wrapper.find('.dropdown-menu').exists()).toBe(true)

    // Click again to close
    await wrapper.find('.menu-button').trigger('click')
    expect(wrapper.find('.dropdown-menu').exists()).toBe(false)
  })

  it('shows correct tooltips for status icons', () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    const iconContainer = wrapper.find('.status-icon-container')
    // Plan 1 has keep_forever: false
    expect(iconContainer.attributes('title')).toBe(i18n.global.t('sidebar.tooltip_temporary'))
  })

  it('handles plan deletion', async () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
    })

    // Mock confirm
    window.confirm = vi.fn(() => true)
    const { apiClient } = await import('@/api/client')
    const trainingPlanStore = useTrainingPlanStore()

    // Open menu
    await wrapper.find('.menu-button').trigger('click')

    // Click delete
    await wrapper.find('.menu-item.delete').trigger('click')

    expect(window.confirm).toHaveBeenCalled()
    expect(apiClient.deletePlan).toHaveBeenCalledWith('1')
    expect(trainingPlanStore.fetchHistory).toHaveBeenCalled()
    expect(push).toHaveBeenCalledWith('/') // Should go home if deleting current plan
  })

  it('handles title editing', async () => {
    const wrapper = mount(Sidebar, {
      global: {
        plugins: [i18n],
      },
      attachTo: document.body, // Needed for focus/blur events sometimes
    })

    const { apiClient } = await import('@/api/client')
    const trainingPlanStore = useTrainingPlanStore()

    // Open menu
    await wrapper.find('.menu-button').trigger('click')

    // Click edit title
    const editButton = wrapper
      .findAll('.menu-item')
      .find((b) => b.text() === i18n.global.t('sidebar.menu_edit_title'))
    await editButton?.trigger('click')

    // Input should appear
    const input = wrapper.find('input.title-input')
    expect(input.exists()).toBe(true)
    expect((input.element as HTMLInputElement).value).toBe('Test Plan 1')

    // Change value and blur to save
    await input.setValue('New Title')
    await input.trigger('blur')

    expect(apiClient.upsertPlan).toHaveBeenCalledWith({
      plan_id: '1',
      title: 'New Title',
      description: 'Desc 1',
      table: [],
    })
    expect(trainingPlanStore.fetchHistory).toHaveBeenCalled()

    // Input should disappear
    expect(wrapper.find('input.title-input').exists()).toBe(false)
  })

  describe('share functionality', () => {
    describe('Non-iOS behavior (one-step share)', () => {
      beforeEach(() => {
        setMockIsIOS(false)
      })

      it('creates share URL and copies immediately for generated plans', async () => {
        const wrapper = mount(Sidebar, {
          global: {
            plugins: [i18n],
          },
        })

        const trainingPlanStore = useTrainingPlanStore()
        const { apiClient } = await import('@/api/client')

        // Open menu
        await wrapper.find('.menu-button').trigger('click')

        // Click share
        const shareButton = wrapper
          .findAll('.menu-item')
          .find((b) => b.text() === i18n.global.t('sidebar.menu_share'))
        await shareButton?.trigger('click')

        // Should call toggleKeepForever
        expect(trainingPlanStore.toggleKeepForever).toHaveBeenCalledWith('1')

        // Should create share URL
        expect(apiClient.createShareUrl).toHaveBeenCalledWith({
          plan_id: '1',
          method: 'link',
        })

        await wrapper.vm.$nextTick()

        // Should copy to clipboard immediately
        expect(mockClipboard.writeText).toHaveBeenCalledWith(
          `${window.location.origin}/shared/test-hash`,
        )
      })
    })

    describe('iOS behavior (two-step share)', () => {
      beforeEach(() => {
        setMockIsIOS(true)
      })

      it('first click creates URL, second click copies for generated plans', async () => {
        const wrapper = mount(Sidebar, {
          global: {
            plugins: [i18n],
          },
        })

        const trainingPlanStore = useTrainingPlanStore()
        const { apiClient } = await import('@/api/client')

        // Open menu
        await wrapper.find('.menu-button').trigger('click')

        // First click - create URL
        const shareButton = wrapper
          .findAll('.menu-item')
          .find((b) => b.text().includes(i18n.global.t('sidebar.menu_share')))
        await shareButton?.trigger('click')

        expect(trainingPlanStore.toggleKeepForever).toHaveBeenCalledWith('1')
        expect(apiClient.createShareUrl).toHaveBeenCalledWith({
          plan_id: '1',
          method: 'link',
        })

        await wrapper.vm.$nextTick()

        // Should NOT copy yet on iOS
        expect(mockClipboard.writeText).not.toHaveBeenCalled()

        // Button text should change to "Copy"
        await wrapper.vm.$nextTick()
        const shareButtonAfter = wrapper
          .findAll('.menu-item')
          .find((b) => b.text().includes(i18n.global.t('share.copy')))
        expect(shareButtonAfter).toBeTruthy()

        // Second click - copy URL
        await shareButtonAfter?.trigger('click')
        await wrapper.vm.$nextTick()

        expect(mockClipboard.writeText).toHaveBeenCalledWith(
          `${window.location.origin}/shared/test-hash`,
        )
      })
    })
  })
})
