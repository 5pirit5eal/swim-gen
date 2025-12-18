// frontend/src/components/forms/__tests__/TrainingPlanForm.spec.ts
import { describe, it, expect, beforeAll, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import TrainingPlanForm from '../TrainingPlanForm.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan.ts'
import i18n from '@/plugins/i18n' // Import the i18n instance
import { apiClient } from '@/api/client'
import type { ApiResult, PromptGenerationResponse } from '@/types'

// Mock the apiClient module for this test file
vi.mock('@/api/client', () => ({
  apiClient: {
    generatePrompt: vi.fn(),
    formatError: vi.fn((err) => err.message),
  },
}))

describe('TrainingPlanForm.vue', () => {
  // Set the locale to 'en' before all tests in this describe block
  beforeAll(() => {
    i18n.global.locale.value = 'en'
  })

  beforeEach(() => {
    vi.clearAllMocks()
    const store = useTrainingPlanStore()
    store.currentPlan = null
    store.isLoading = false
    store.error = null
  })

  it('renders correctly without errors', () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    // The simplest test: does it mount?
    expect(wrapper.exists()).toBe(true)
  })
  it('disables the submit button when the input is empty', () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const submitButton = wrapper.find('button[type="submit"]')

    expect(submitButton.attributes('disabled')).toBeDefined()
  })
  it('enables the submit button when the input is valid', async () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const textarea = wrapper.find('textarea')
    const submitButton = wrapper.find('button[type="submit"]')

    // Set a valid value for the textarea
    await textarea.setValue('I need a workout plan.')

    // Now, the 'disabled' attribute should be gone
    expect(submitButton.attributes('disabled')).toBeUndefined()
  })
  it('disables the submit button when the input is too long', async () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const textarea = wrapper.find('textarea')
    const submitButton = wrapper.find('button[type="submit"]')

    // First, set a valid value to enable the button
    await textarea.setValue('Valid input.')
    expect(submitButton.attributes('disabled')).toBeUndefined()

    // Now, set a value that is too long
    const longText = 'a'.repeat(3001)
    await textarea.setValue(longText)

    // The button should be disabled again
    expect(submitButton.attributes('disabled')).toBeDefined()
  })
  it('shows an error message when the input is too long', async () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const textarea = wrapper.find('textarea')

    // Set a value that is too long
    const longText = 'a'.repeat(3001)
    await textarea.setValue(longText)

    const errorMessage = wrapper.find('.form-hint.text-warning')

    // The error message should be visible
    expect(errorMessage.exists()).toBe(true)
    expect(errorMessage.text()).toContain(i18n.global.t('form.request_too_long'))
  })
  it('disables the submit button when generating', async () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const store = useTrainingPlanStore()

    const submitButton = wrapper.find('button[type="submit"]')
    const textarea = wrapper.find('textarea')

    // Set isGenerating to disable the button
    await textarea.setValue('Valid input.')
    expect(submitButton.attributes('disabled')).toBeUndefined()
    expect(textarea.attributes('disabled')).toBeUndefined()

    store.isLoading = true
    await wrapper.vm.$nextTick()

    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(textarea.attributes('disabled')).toBeDefined()
    expect(submitButton.text()).toContain(i18n.global.t('form.generating_plan'))
  })
  it('disables the prompt generation button when generating', async () => {
    const mockSuccessResponse: ApiResult<PromptGenerationResponse> = {
      success: true,
      data: { prompt: 'This is a generated prompt' },
    }

    // Set up a controllable promise
    let resolvePromise: (value: ApiResult<PromptGenerationResponse>) => void
    const promise = new Promise<ApiResult<PromptGenerationResponse>>((resolve) => {
      resolvePromise = resolve
    })

    // Mock the implementation to return our controllable promise
    vi.mocked(apiClient.generatePrompt).mockReturnValue(promise)

    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })

    const promptButton = wrapper.findAll('.toggle-settings-btn')[1]!
    const textarea = wrapper.find('textarea')

    // Initially, the button should be enabled
    expect(promptButton.attributes('disabled')).toBeFalsy()
    expect(promptButton.text()).toContain(i18n.global.t('form.i_feel_lucky'))

    // Click the button to start generation
    await promptButton.trigger('click')

    // At this point, the API call is pending.
    // The button should be disabled and show "Generating..."
    expect(promptButton.attributes('disabled')).toBeDefined()
    expect(promptButton.text()).toContain(i18n.global.t('form.generating'))

    // Now, resolve the promise to simulate the API call completing
    resolvePromise!(mockSuccessResponse)

    // Wait for the promise to be processed and the DOM to update
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    // After generation, the button should be enabled again
    expect(promptButton.attributes('disabled')).toBeFalsy()
    expect(promptButton.text()).toContain(i18n.global.t('form.i_feel_lucky'))

    // And the textarea should be updated with the new prompt
    expect(textarea.element.value).toBe(mockSuccessResponse.data?.prompt)
  })
})
