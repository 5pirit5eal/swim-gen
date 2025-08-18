// frontend/src/components/forms/__tests__/TrainingPlanForm.spec.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TrainingPlanForm from '../TrainingPlanForm.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan.ts'
import i18n from '@/plugins/i18n' // Import the i18n instance

describe('TrainingPlanForm.vue', () => {
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
  })
  it('disables the prompt generation button when generating', async () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const store = useTrainingPlanStore()

    const promptButton = wrapper.find('.toggle-settings-btn')

    // Initially, the button should be enabled
    expect(promptButton.attributes('disabled')).toBeFalsy()

    // Set isLoading to true to disable the button
    store.isLoading = true
    await wrapper.vm.$nextTick()

    expect(promptButton.attributes('disabled')).toBeDefined()
  })
})
