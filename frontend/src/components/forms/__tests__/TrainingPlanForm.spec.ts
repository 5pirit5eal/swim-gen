// frontend/src/components/forms/__tests__/TrainingPlanForm.spec.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TrainingPlanForm from '../TrainingPlanForm.vue'

describe('TrainingPlanForm.vue', () => {
  it('renders correctly without errors', () => {
    const wrapper = mount(TrainingPlanForm)
    // The simplest test: does it mount?
    expect(wrapper.exists()).toBe(true)
  })
})
