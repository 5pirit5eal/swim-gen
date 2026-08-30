// frontend/src/components/forms/__tests__/TrainingPlanForm.spec.ts
import { describe, it, expect, beforeAll, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import TrainingPlanForm from '../TrainingPlanForm.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan.ts'
import { useSettingsStore } from '@/stores/settings.ts'
import { useAuthStore } from '@/stores/auth.ts'
import { useProfileStore } from '@/stores/profile.ts'
import type { User } from '@supabase/supabase-js'
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

// Mock supabase
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

    const authStore = useAuthStore()
    authStore.user = null

    const profileStore = useProfileStore()
    profileStore.profile = null

    const settingsStore = useSettingsStore()
    settingsStore.selectedAudience = null
    settingsStore.clearFilters()
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
    expect(promptButton.text()).toContain(i18n.global.t('form.suggest_workout'))

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
    expect(promptButton.text()).toContain(i18n.global.t('form.suggest_workout'))

    // And the textarea should be updated with the new prompt
    expect(textarea.element.value).toBe(mockSuccessResponse.data?.prompt)
  })

  it('renders audience selector and adapts placeholder and settings when audience changes', async () => {
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })
    const settingsStore = useSettingsStore()

    const audienceButtons = wrapper.findAll('.audience-btn')
    expect(audienceButtons.length).toBe(4)

    // Button labels should match Beginner, Triathlete, Competitive Swimmer, Hobby
    expect(audienceButtons[0]!.text()).toBe(i18n.global.t('form.audience_beginner'))
    expect(audienceButtons[1]!.text()).toBe(i18n.global.t('form.audience_triathlete'))
    expect(audienceButtons[2]!.text()).toBe(i18n.global.t('form.audience_competitive_swimmer'))
    expect(audienceButtons[3]!.text()).toBe(i18n.global.t('form.audience_hobby'))

    let textarea = wrapper.find('textarea')
    expect(textarea.attributes('placeholder')).toBe(i18n.global.t('form.example_placeholder'))

    // Click 'beginner' button (index 0)
    await audienceButtons[0]!.trigger('click')
    expect(settingsStore.selectedAudience).toBe('beginner')
    expect(settingsStore.filters.schwierigkeitsgrad).toBe('Anfaenger')
    expect(settingsStore.filters.trainingstyp).toBe('Techniktraining')
    await wrapper.vm.$nextTick()

    textarea = wrapper.find('textarea')
    expect(textarea.attributes('placeholder')).toBe(
      i18n.global.t('form.example_placeholder_beginner'),
    )

    // Click 'triathlete' button (index 1)
    await audienceButtons[1]!.trigger('click')
    expect(settingsStore.selectedAudience).toBe('triathlete')
    expect(settingsStore.filters.schwierigkeitsgrad).toBe('Fortgeschritten')
    expect(settingsStore.filters.trainingstyp).toBe('Grundlagenausdauer')
    await wrapper.vm.$nextTick()

    textarea = wrapper.find('textarea')
    expect(textarea.attributes('placeholder')).toBe(
      i18n.global.t('form.example_placeholder_triathlete'),
    )

    // Click 'competitive_swimmer' button (index 2)
    await audienceButtons[2]!.trigger('click')
    expect(settingsStore.selectedAudience).toBe('competitive_swimmer')
    expect(settingsStore.filters.schwierigkeitsgrad).toBe('Leistungsschwimmer')
    expect(settingsStore.filters.trainingstyp).toBe('Wettkampfvorbereitung')
    await wrapper.vm.$nextTick()

    textarea = wrapper.find('textarea')
    expect(textarea.attributes('placeholder')).toBe(
      i18n.global.t('form.example_placeholder_competitive_swimmer'),
    )

    // Click 'hobby' button (index 3)
    await audienceButtons[3]!.trigger('click')
    expect(settingsStore.selectedAudience).toBe('hobby')
    expect(settingsStore.filters.schwierigkeitsgrad).toBeUndefined()
    expect(settingsStore.filters.trainingstyp).toBeUndefined()
    await wrapper.vm.$nextTick()

    textarea = wrapper.find('textarea')
    expect(textarea.attributes('placeholder')).toBe(i18n.global.t('form.example_placeholder_hobby'))
  })

  it('shows tooltip on audience button hover after delay', async () => {
    vi.useFakeTimers()
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
      attachTo: document.body,
    })

    const beginnerBtn = wrapper.findAll('.audience-btn')[0]!
    await beginnerBtn.trigger('mouseenter')

    // Tooltip should not be displayed immediately
    expect(document.getElementById('audience-hover-tooltip')).toBeNull()

    // Fast forward timer past 500ms
    vi.advanceTimersByTime(550)
    await wrapper.vm.$nextTick()

    // Tooltip should now be visible in DOM
    const tooltip = document.getElementById('audience-hover-tooltip')
    expect(tooltip).not.toBeNull()
    expect(tooltip?.textContent).toContain(i18n.global.t('form.audience_hint_beginner'))

    // Mouse leave should hide tooltip
    await beginnerBtn.trigger('mouseleave')
    await wrapper.vm.$nextTick()
    expect(document.getElementById('audience-hover-tooltip')).toBeNull()

    vi.useRealTimers()
    wrapper.unmount()
  })

  it('hides audience selector when user is logged in and uses profile category for prompt generation', async () => {
    const authStore = useAuthStore()
    // Mock logged in user
    authStore.user = { id: 'test-user', email: 'test@example.com' } as unknown as User

    const profileStore = useProfileStore()
    profileStore.profile = {
      user_id: 'test-user',
      updated_at: '',
      username: 'test',
      experience: null,
      preferred_language: null,
      preferred_strokes: [],
      categories: ['Triathlet'],
      overall_generations: 0,
      monthly_generations: 0,
      exports: 0,
      css_200m_seconds: null,
      css_400m_seconds: null,
    }

    const mockSuccessResponse: ApiResult<PromptGenerationResponse> = {
      success: true,
      data: { prompt: 'Triathlete prompt' },
    }
    vi.mocked(apiClient.generatePrompt).mockResolvedValue(mockSuccessResponse)

    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })

    // Audience section should be hidden
    expect(wrapper.find('.audience-section').exists()).toBe(false)

    // Click the feel lucky / prompt generation button
    const promptButton = wrapper.findAll('.toggle-settings-btn')[1]!
    await promptButton.trigger('click')

    expect(apiClient.generatePrompt).toHaveBeenCalledWith({
      language: 'en',
      audience: 'triathlete',
    })
  })

  it('shows tooltip when hovering over disabled submit button and hides when enabled', async () => {
    vi.useFakeTimers()
    const wrapper = mount(TrainingPlanForm, {
      global: {
        plugins: [i18n],
      },
    })

    const submitWrapper = wrapper.find('.submit-btn-wrapper')
    const textarea = wrapper.find('textarea')

    // Form is empty -> submit is disabled
    expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()

    // Hover submit wrapper
    await submitWrapper.trigger('mouseenter')
    vi.advanceTimersByTime(250)
    await wrapper.vm.$nextTick()

    // Tooltip should be visible
    const tooltip = document.getElementById('submit-disabled-tooltip')
    expect(tooltip).not.toBeNull()
    expect(tooltip?.textContent).toContain(
      i18n.global.t('form.generate_training_plan_disabled_tooltip'),
    )

    // Leave hover
    await submitWrapper.trigger('mouseleave')
    await wrapper.vm.$nextTick()
    expect(document.getElementById('submit-disabled-tooltip')).toBeNull()

    // When valid text is provided, entering hover should not show tooltip
    await textarea.setValue('Valid training plan request')
    await submitWrapper.trigger('mouseenter')
    vi.advanceTimersByTime(250)
    await wrapper.vm.$nextTick()
    expect(document.getElementById('submit-disabled-tooltip')).toBeNull()

    vi.useRealTimers()
  })
})
