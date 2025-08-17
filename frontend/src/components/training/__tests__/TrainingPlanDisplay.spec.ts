// frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TrainingPlanDisplay from '../TrainingPlanDisplay.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import i18n from '@/plugins/i18n' // Import the i18n instance
describe('TrainingPlanDisplay.vue', () => {
  it('renders a placeholder when there is no plan', () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n],
      },
    })

    // Check for the placeholder text
    expect(wrapper.text()).toContain(i18n.global.t('trainingPlanDisplay.noPlanPlaceholder'))
  })
  it('renders the training plan when it exists', async () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n],
      },
    })
    const store = useTrainingPlanStore()

    const mockPlan = {
      title: 'Mock Training Plan',
      description: 'This is a mock training plan for testing.',
      table: [
        { Amount: 1, Multiplier: 'x', Distance: 100, Break: '30s', Content: 'Warm-up', Intensity: 'Easy', Sum: 100 },
        { Amount: 1, Multiplier: '', Distance: 0, Break: '', Content: 'Total', Intensity: '', Sum: 100 }
      ]
    }

    store.currentPlan = mockPlan

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain(mockPlan.title)
    expect(wrapper.text()).toContain(mockPlan.description)

    expect(wrapper.text()).not.toContain(i18n.global.t('trainingPlanDisplay.noPlanPlaceholder'))
  })
})
