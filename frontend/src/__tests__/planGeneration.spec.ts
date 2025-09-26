// frontend/src/__tests__/planGeneration.spec.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '@/App.vue'
import { apiClient } from '@/api/client'
import type { Mock } from 'vitest'
import type { RAGResponse, ApiResult } from '@/types'
import router from '@/router' // Import the router instance
import i18n from '@/plugins/i18n' // Import the i18n instance

// Mock the apiClient module for integration tests
vi.mock('@/api/client', () => ({
  apiClient: {
    query: vi.fn(),
    exportPDF: vi.fn(),
  },
}))

// Cast apiClient.query to a Mock type for TypeScript
const mockedApiQuery = apiClient.query as Mock<typeof apiClient.query>

describe('Plan Generation End-to-End Workflow', () => {
  beforeEach(async () => {
    // Make beforeEach async
    vi.clearAllMocks()
    // Reset Pinia stores if necessary (Pinia setup in vitest.setup.ts handles activation)
    // For integration tests, we often want a fresh store state
    // You might need to manually reset specific stores if their state persists
    // For now, we rely on the component remounting to reset its internal state

    // Ensure router is at the base path before each test
    router.push('/')
    await router.isReady() // Wait for router to be ready
  })

  it('successfully generates and displays a training plan', async () => {
    const mockPlanResponse: ApiResult<RAGResponse> = {
      success: true,
      data: {
        title: 'My Awesome Test Plan',
        description: 'This plan was generated during an integration test.',
        table: [
          {
            Amount: 1,
            Multiplier: 'x',
            Distance: 100,
            Break: '1min',
            Content: 'Warm-up',
            Intensity: 'Easy',
            Sum: 100,
          },
          {
            Amount: 2,
            Multiplier: 'x',
            Distance: 50,
            Break: '30s',
            Content: 'Drill',
            Intensity: 'Moderate',
            Sum: 100,
          },
          {
            Amount: 1,
            Multiplier: '',
            Distance: 0,
            Break: '',
            Content: 'Total',
            Intensity: '',
            Sum: 200,
          },
        ],
      },
    }

    mockedApiQuery.mockResolvedValue(mockPlanResponse)

    // Mount the entire App component and provide the router
    const wrapper = mount(App, {
      global: {
        plugins: [router, i18n], // Provide the router and i18n instances
      },
    })

    const textarea = wrapper.find('textarea')
    const submitButton = wrapper.find('button[type="submit"]')

    await textarea.setValue('Generate a simple test plan.')
    await submitButton.trigger('submit')

    // Wait for the DOM to update after the plan is generated and displayed
    // Use a more robust wait strategy for the plan title to appear
    await wrapper.find('.plan-title') // Wait until the plan title element appears

    expect(wrapper.text()).toContain(mockPlanResponse.data!.title)
    expect(wrapper.text()).toContain(mockPlanResponse.data!.description)

    expect(mockedApiQuery).toHaveBeenCalledTimes(1)
    expect(mockedApiQuery).toHaveBeenCalledWith({
      content: 'Generate a simple test plan.',
      method: 'generate',
      filter: {},
      language: 'en-US', // Expect language to be included
      pool_length: 25, // Expect default pool length to be included
    })

    expect(wrapper.find('.loading-state').exists()).toBe(false)
  })
})
