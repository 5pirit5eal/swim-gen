// frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import TrainingPlanDisplay from '../TrainingPlanDisplay.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import i18n from '@/plugins/i18n' // Import the i18n instance

describe('TrainingPlanDisplay.vue', () => {
  const mockPlan = {
    title: 'Mock Training Plan',
    description: 'This is a mock training plan for testing.',
    table: [
      {
        Amount: 1,
        Multiplier: 'x',
        Distance: 100,
        Break: '30s',
        Content: 'Warm-up',
        Intensity: 'Easy',
        Sum: 100,
      },
      {
        Amount: 1,
        Multiplier: '',
        Distance: 0,
        Break: '',
        Content: 'Total',
        Intensity: '',
        Sum: 100,
      },
    ],
  }

  it('renders a placeholder when there is no plan', () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
      },
    })

    // Check for the placeholder text
    expect(wrapper.text()).toContain(i18n.global.t('display.no_plan_placeholder'))
  })

  it('renders the training plan when it exists', async () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
      },
    })
    const store = useTrainingPlanStore()
    store.currentPlan = mockPlan

    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain(mockPlan.title)
    expect(wrapper.text()).toContain(mockPlan.description)
    expect(wrapper.text()).not.toContain(i18n.global.t('display.no_plan_placeholder'))
  })

  describe('Editing Training Plan', () => {
    it('allows editing the Amount field with a valid number', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(mockPlan)) // Deep copy
      await wrapper.vm.$nextTick() // Wait for the DOM to update

      // Enable editing
      await wrapper.find('button.export-btn').trigger('click')
      // @ts-expect-error: isEditing is not typed on the wrapper vm
      expect(wrapper.vm.isEditing).toBe(true)

      // Find the first Amount cell and click to start editing
      const amountCell = wrapper.find('tbody tr:first-child td:first-child')
      await amountCell.trigger('click')

      // Find the input, set a new value, and trigger blur
      const input = amountCell.find('input')
      await input.setValue('5')
      await input.trigger('blur')

      // Check if the store's update function was called correctly
      expect(store.updatePlanRow).toHaveBeenCalledWith(0, 'Amount', 5)
    })

    it('reverts to original value when editing Amount with invalid input', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(mockPlan)) // Deep copy
      await wrapper.vm.$nextTick() // Wait for the DOM to update

      // Enable editing
      await wrapper.find('button.export-btn').trigger('click')

      // Find the first Amount cell and click to start editing
      const amountCell = wrapper.find('tbody tr:first-child td:first-child')
      await amountCell.trigger('click')

      // Find the input, set an invalid value, and trigger blur
      const input = amountCell.find('input')
      await input.setValue('abc')
      await input.trigger('blur')

      // Check that updatePlanRow was called with the original value
      expect(store.updatePlanRow).toHaveBeenCalledWith(0, 'Amount', mockPlan.table[0].Amount)
    })
  })
})
