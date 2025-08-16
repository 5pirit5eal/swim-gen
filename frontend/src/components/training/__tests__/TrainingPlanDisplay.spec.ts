// frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TrainingPlanDisplay from '../TrainingPlanDisplay.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'

describe('TrainingPlanDisplay.vue', () => {
  it('renders a placeholder when there is no plan', () => {
    const wrapper = mount(TrainingPlanDisplay)

    // Check for the placeholder text
    expect(wrapper.text()).toContain('No training plan generated yet. Use the form above to create one!')
  })
  it('renders the training plan when it exists', async () => {
    const wrapper = mount(TrainingPlanDisplay)
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

    expect(wrapper.text()).not.toContain('No training plan generated yet. Use the form above to create one!')
  })
})
