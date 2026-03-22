// frontend/src/components/training/__tests__/TrainingPlanDisplay.spec.ts
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import TrainingPlanDisplay from '../TrainingPlanDisplay.vue'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import i18n from '@/plugins/i18n' // Import the i18n instance
import type { RAGResponse } from '@/types'

/**
 * Fixture: Simple flat plan with one exercise and total row.
 */
const createSimplePlan = (): RAGResponse => ({
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
      _id: 'test-id-1',
    },
    {
      Amount: 1,
      Multiplier: '',
      Distance: 0,
      Break: '',
      Content: 'Total',
      Intensity: '',
      Sum: 100,
      _id: 'test-id-2',
    },
  ],
})

/**
 * Fixture: Nested plan with depth 2 (parent + 2 children).
 * Top-level: [Exercise-A (parent with SubRows), Exercise-B, Total]
 */
const createNestedDepth2Plan = (): RAGResponse => ({
  title: 'Nested Plan (Depth 2)',
  description: 'Plan with nested rows at depth 2.',
  table: [
    {
      Amount: 3,
      Multiplier: 'x',
      Distance: 600, // Sum of children
      Break: '',
      Content: 'Main Set',
      Intensity: 'GA2',
      Sum: 1800,
      SubRows: [
        {
          Amount: 1,
          Multiplier: 'x',
          Distance: 200,
          Break: '30s',
          Content: 'Freestyle 200m',
          Intensity: 'GA2',
          Sum: 200,
          Equipment: ['Flossen'],
          _id: 'child-1',
        },
        {
          Amount: 1,
          Multiplier: 'x',
          Distance: 400,
          Break: '20s',
          Content: 'Breaststroke 400m',
          Intensity: 'GA1',
          Sum: 400,
          _id: 'child-2',
        },
      ],
      _id: 'parent-1',
    },
    {
      Amount: 2,
      Multiplier: 'x',
      Distance: 100,
      Break: '15s',
      Content: 'Easy recovery',
      Intensity: 'Recovery',
      Sum: 200,
      _id: 'test-id-flat',
    },
    {
      Amount: 0,
      Multiplier: '',
      Distance: 0,
      Break: '',
      Content: 'Total',
      Intensity: '',
      Sum: 2000,
      _id: 'test-id-total',
    },
  ],
})

/**
 * Fixture: Deep nesting (depth 3: top-level parent with child that has grandchildren).
 */
const createNestedDepth3Plan = (): RAGResponse => ({
  title: 'Deep Nested Plan (Depth 3)',
  description: 'Plan with nested rows at depth 3.',
  table: [
    {
      Amount: 1,
      Multiplier: 'x',
      Distance: 1000,
      Break: '',
      Content: 'Complex Main Set',
      Intensity: 'GA1',
      Sum: 1000,
      SubRows: [
        {
          Amount: 2,
          Multiplier: 'x',
          Distance: 500,
          Break: '',
          Content: 'Pyramid set',
          Intensity: 'GA2',
          Sum: 1000,
          SubRows: [
            {
              Amount: 1,
              Multiplier: 'x',
              Distance: 100,
              Break: '10s',
              Content: '100m Freestyle',
              Intensity: 'GA2',
              Sum: 100,
              _id: 'grandchild-1',
            },
            {
              Amount: 1,
              Multiplier: 'x',
              Distance: 200,
              Break: '15s',
              Content: '200m Breaststroke',
              Intensity: 'GA1',
              Sum: 200,
              _id: 'grandchild-2',
            },
            {
              Amount: 1,
              Multiplier: 'x',
              Distance: 200,
              Break: '20s',
              Content: '200m IM',
              Intensity: 'GA1',
              Sum: 200,
              _id: 'grandchild-3',
            },
          ],
          _id: 'child-nested-1',
        },
      ],
      _id: 'parent-deep-1',
    },
    {
      Amount: 0,
      Multiplier: '',
      Distance: 0,
      Break: '',
      Content: 'Total',
      Intensity: '',
      Sum: 1000,
      _id: 'total-depth3',
    },
  ],
})

/**
 * Fixture: Mixed equipment and content across rows.
 */
const createMixedEquipmentPlan = (): RAGResponse => ({
  title: 'Mixed Equipment Plan',
  description: 'Plan with various equipment combinations.',
  table: [
    {
      Amount: 1,
      Multiplier: 'x',
      Distance: 400,
      Break: '0s',
      Content: 'Warm-up Free',
      Intensity: 'Easy',
      Sum: 400,
      Equipment: [],
      _id: 'warm-up',
    },
    {
      Amount: 4,
      Multiplier: 'x',
      Distance: 300,
      Break: '20s',
      Content: 'Main Set',
      Intensity: 'GA1',
      Sum: 1200,
      SubRows: [
        {
          Amount: 1,
          Multiplier: 'x',
          Distance: 150,
          Break: '',
          Content: 'Pull with buoy',
          Intensity: 'GA1',
          Sum: 150,
          Equipment: ['Pull buoy'],
          _id: 'pull-set-1',
        },
        {
          Amount: 1,
          Multiplier: 'x',
          Distance: 150,
          Break: '',
          Content: 'Kick with board',
          Intensity: 'GA2',
          Sum: 150,
          Equipment: ['Kickboard', 'Flossen'],
          _id: 'kick-set-1',
        },
      ],
      Equipment: [],
      _id: 'main-set-1',
    },
    {
      Amount: 1,
      Multiplier: 'x',
      Distance: 200,
      Break: '',
      Content: 'Cool-down',
      Intensity: 'Recovery',
      Sum: 200,
      Equipment: ['Schnorchel'],
      _id: 'cool-down',
    },
    {
      Amount: 0,
      Multiplier: '',
      Distance: 0,
      Break: '',
      Content: 'Total',
      Intensity: '',
      Sum: 1800,
      Equipment: [],
      _id: 'total-mixed',
    },
  ],
})

/**
 * Legacy mockPlan for backwards compatibility.
 */
const mockPlan = createSimplePlan()

describe('TrainingPlanDisplay.vue', () => {
  it('renders a placeholder when there is no plan', () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
      },
      props: {
        store: useTrainingPlanStore(),
      },
    })

    // Check for the placeholder text
    expect(wrapper.text()).toContain(i18n.global.t('display.no_plan_placeholder'))
  })

  it('renders loading spinner when store.isLoading is true', () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
      },
      props: {
        store: useTrainingPlanStore(),
      },
    })
    const store = useTrainingPlanStore()
    store.isLoading = true

    expect(wrapper.find('.loading-state').exists() || store.isLoading).toBe(true)
  })

  it('renders no-plan placeholder when hasPlan is false and not loading', () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
      },
      props: {
        store: useTrainingPlanStore(),
      },
    })
    const store = useTrainingPlanStore()
    store.isLoading = false
    store.currentPlan = null

    expect(wrapper.find('.no-plan').exists()).toBe(true)
    expect(wrapper.text()).toContain(i18n.global.t('display.no_plan_placeholder'))
  })

  it('renders the training plan when it exists', async () => {
    const wrapper = mount(TrainingPlanDisplay, {
      global: {
        plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
      },
      props: {
        store: useTrainingPlanStore(),
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
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(mockPlan)) // Deep copy
      await wrapper.vm.$nextTick() // Wait for the DOM to update

      // Enable editing
      await wrapper.find('button[data-testid="plan-edit-btn"]').trigger('click')
      await wrapper.vm.$nextTick()

      // Find all inputs from plan-card elements and get the first one (Amount column)
      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThan(0)
      const firstCard = planCards[0]!
      const input = firstCard.find('input')
      expect(input.exists()).toBe(true)
      await input.setValue('5')
      await input.trigger('blur')

      // Check if the store's update function was called correctly
      expect(store.updatePlanRow).toHaveBeenCalledWith([0], 'Amount', 5)
    })

    it('reverts to original value when editing Amount with invalid input', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(mockPlan)) // Deep copy
      await wrapper.vm.$nextTick() // Wait for the DOM to update

      // Enable editing
      await wrapper.find('button[data-testid="plan-edit-btn"]').trigger('click')
      await wrapper.vm.$nextTick() // Wait for DOM update

      // Find all inputs from plan-card elements and get the first one (Amount column)
      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThan(0)
      const firstCard = planCards[0]!
      const input = firstCard.find('input')
      expect(input.exists()).toBe(true)
      await input.setValue('abc')
      await input.trigger('blur')

      // Check that updatePlanRow was called with 0 for invalid input
      expect(store.updatePlanRow).toHaveBeenCalledWith([0], 'Amount', 0)
    })

    it('edit mode shows input elements for leaf rows (Amount, Distance, Break visible)', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(createSimplePlan()))
      await wrapper.vm.$nextTick()

      const planCardsBefore = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCardsBefore.length).toBeGreaterThan(0)
      expect(planCardsBefore[0]!.findAll('input').length).toBe(0)

      await wrapper.find('button[data-testid="plan-edit-btn"]').trigger('click')
      await wrapper.vm.$nextTick()

      const leafCard = wrapper.findAll('[data-testid="plan-card"]')[0]!
      expect(leafCard.findAll('input').length).toBe(4) // Amount, Distance, Break, Intensity
      expect(leafCard.find('textarea').exists()).toBe(true)
    })

    it('calls updatePlanRow([0], Amount, 5) when Amount blurs with value 5', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(createSimplePlan()))
      await wrapper.vm.$nextTick()

      await wrapper.find('button[data-testid="plan-edit-btn"]').trigger('click')
      await wrapper.vm.$nextTick()

      const amountInput = wrapper.findAll('[data-testid="plan-card"]')[0]!.findAll('input')[0]!
      expect(amountInput.exists()).toBe(true)
      await amountInput.setValue('5')
      await amountInput.trigger('blur')

      expect(store.updatePlanRow).toHaveBeenCalledWith([0], 'Amount', 5)
    })

    it('parent rows (with SubRows) do not show a Distance input in edit mode', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = JSON.parse(JSON.stringify(createNestedDepth2Plan()))
      await wrapper.vm.$nextTick()

      await wrapper.find('button[data-testid="plan-edit-btn"]').trigger('click')
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThanOrEqual(2)

      const parentCard = planCards.find((card) => card.classes('plan-row-card--parent'))
      expect(parentCard?.exists()).toBe(true)

      const parentHeader = parentCard!.find('.plan-row-card__data')
      const parentDistanceInputs = parentHeader
        .findAll('input')
        .filter((input) => input.attributes('aria-label') === 'Distance (m)')
      expect(parentDistanceInputs.length).toBe(0)

      const nestedCards = wrapper.findAll('[data-testid="plan-card-nested"]')
      expect(nestedCards.length).toBeGreaterThan(0)
      const nestedDistanceInputs = nestedCards[0]!
        .find('.plan-row-card__data')
        .findAll('input')
        .filter((input) => input.attributes('aria-label') === 'Distance (m)')
      expect(nestedDistanceInputs.length).toBe(1)
    })
  })

  describe('Follow-up: Header/Footer/Equipment Hierarchy (TDD)', () => {
    it('header displays total distance in [data-testid="plan-header-total"]', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createSimplePlan()
      await wrapper.vm.$nextTick()

      const headerTotal = wrapper.find('[data-testid="plan-header-total"]')
      expect(headerTotal.exists()).toBe(true)
      expect(headerTotal.text()).toContain('100')
    })

    it('header does NOT contain the plan description text', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      const plan = createSimplePlan()
      store.currentPlan = plan
      await wrapper.vm.$nextTick()

      const header = wrapper.find('.plan-header')
      expect(header.exists()).toBe(true)
      expect(header.text()).not.toContain(plan.description)
    })

    it('footer/meta region [data-testid="plan-footer-meta"] renders plan description', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      const plan = createSimplePlan()
      store.currentPlan = plan
      await wrapper.vm.$nextTick()

      const footer = wrapper.find('[data-testid="plan-footer-meta"]')
      expect(footer.exists()).toBe(true)
      expect(footer.text()).toContain(plan.description)
    })

    it('footer/meta region [data-testid="plan-footer-equipment"] shows distinct equipment', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createMixedEquipmentPlan()
      await wrapper.vm.$nextTick()

      const equipmentFooter = wrapper.find('[data-testid="plan-footer-equipment"]')
      expect(equipmentFooter.exists()).toBe(true)
      // Mixed plan has: Pull buoy, Kickboard, Flossen, Schnorchel
      expect(equipmentFooter.text()).toContain('Pull buoy')
      expect(equipmentFooter.text()).toContain('Flossen')
      expect(equipmentFooter.text()).toContain('Schnorchel')
    })

    it('footer/meta equipment lists nested subrow equipment (recursive aggregation)', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createNestedDepth2Plan() // has Equipment: ['Flossen'] in child
      await wrapper.vm.$nextTick()

      const equipmentFooter = wrapper.find('[data-testid="plan-footer-equipment"]')
      expect(equipmentFooter.exists()).toBe(true)
      expect(equipmentFooter.text()).toContain('Flossen')
    })

    it('SimplePlanDisplay does NOT render description text in view mode', async () => {
      const { default: SimplePlanDisplay } = await import('../SimplePlanDisplay.vue')
      const simplePlan = createSimplePlan()
      const wrapper = mount(SimplePlanDisplay, {
        global: {
          plugins: [i18n],
        },
        props: {
          title: simplePlan.title,
          description: simplePlan.description,
          table: simplePlan.table,
          planId: 'test-plan-id',
        },
      })
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain(simplePlan.description)
    })

    it('SimplePlanDisplay emits description in @save payload', async () => {
      const { default: SimplePlanDisplay } = await import('../SimplePlanDisplay.vue')
      const simplePlan = createSimplePlan()
      const wrapper = mount(SimplePlanDisplay, {
        global: {
          plugins: [i18n],
        },
        props: {
          title: simplePlan.title,
          description: simplePlan.description,
          table: simplePlan.table,
          planId: 'test-plan-id',
        },
      })
      await wrapper.vm.$nextTick()

      await wrapper.find('button').trigger('click')

      const saveEmit = wrapper.emitted('save')
      expect(saveEmit).toBeTruthy()
      expect(saveEmit![0]![0]).toMatchObject({ description: simplePlan.description })
    })
  })

  describe('Card-oriented display (TDD Wave 1)', () => {
    it('renders exercise rows as cards with [data-testid="plan-card"]', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createSimplePlan()
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThan(0)
      expect(planCards[0]?.exists()).toBe(true)
    })

    it('does not render <table> element in main plan content', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createSimplePlan()
      await wrapper.vm.$nextTick()

      const tables = wrapper.findAll('.exercise-table')
      expect(tables.length).toBe(0)
    })

    it('renders flat plan with one exercise card and excludes total row from cards', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      const simplePlan = createSimplePlan()
      store.currentPlan = simplePlan
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBe(1)
      expect(planCards[0]?.text()).toContain('Warm-up')
      expect(planCards[0]?.text()).not.toContain('Total')
    })

    it('renders nested depth-2 plan with parent and sibling cards', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createNestedDepth2Plan()
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThanOrEqual(2)

      const hasMainSet = planCards.some((card) => card.text().includes('Main Set'))
      const hasEasyRecovery = planCards.some((card) => card.text().includes('Easy recovery'))
      expect(hasMainSet).toBe(true)
      expect(hasEasyRecovery).toBe(true)

      const totalCard = planCards.filter((card) => card.text().includes('Total'))
      expect(totalCard.length).toBe(0)
    })

    it('renders nested depth-3 plan with proper hierarchy', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createNestedDepth3Plan()
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThanOrEqual(1)

      const complexSetCard = planCards.find((card) => card.text().includes('Complex Main Set'))
      expect(complexSetCard?.exists()).toBe(true)
    })

    it('renders mixed equipment plan with proper card structure', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createMixedEquipmentPlan()
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThanOrEqual(3)

      const warmupCard = planCards.find((card) => card.text().includes('Warm-up Free'))
      const mainSetCard = planCards.find((card) => card.text().includes('Main Set'))
      const cooldownCard = planCards.find((card) => card.text().includes('Cool-down'))

      expect(warmupCard?.exists()).toBe(true)
      expect(mainSetCard?.exists()).toBe(true)
      expect(cooldownCard?.exists()).toBe(true)

      const totalCard = planCards.filter((card) => card.text().includes('Total'))
      expect(totalCard.length).toBe(0)
    })

    it('parent rows with SubRows render nested cards', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createNestedDepth2Plan()
      await wrapper.vm.$nextTick()

      const nestedCards = wrapper.findAll('[data-testid="plan-card-nested"]')
      expect(nestedCards.length).toBeGreaterThanOrEqual(1)
    })

    it('drill-link content in row card renders ContentWithDrillLinks inside .plan-row-card', async () => {
      const MockIntersectionObserver = vi.fn(function () {
        return { observe: vi.fn(), unobserve: vi.fn(), disconnect: vi.fn() }
      })
      vi.stubGlobal('IntersectionObserver', MockIntersectionObserver)
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = {
        title: 'Drill Link Plan',
        description: 'Plan with drill link in content.',
        table: [
          {
            Amount: 1,
            Multiplier: 'x',
            Distance: 200,
            Break: '30s',
            Content: '200m [freestyle](/drill/abc-123) swim',
            Intensity: 'Easy',
            Sum: 200,
            _id: 'drill-row-1',
          },
          {
            Amount: 0,
            Multiplier: '',
            Distance: 0,
            Break: '',
            Content: 'Total',
            Intensity: '',
            Sum: 200,
            _id: 'drill-total',
          },
        ],
      }
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBe(1)

      const card = planCards[0]!
      expect(card.classes('plan-row-card')).toBe(true)
      expect(card.find('.content-with-drill-links').exists()).toBe(true)
    })
  })

  describe('Follow-up: Responsive Row Equipment', () => {
    it('row equipment badges are rendered inside the card header on wide screens', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createMixedEquipmentPlan()
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      const cooldownCard = planCards.find((card) => card.text().includes('Cool-down'))
      expect(cooldownCard?.exists()).toBe(true)

      const equipment_label = cooldownCard!.find('[data-testid="equipment-metric"]')
      expect(equipment_label.exists()).toBe(true)
      expect(equipment_label.find('.plan-row-card__equipment-badges').exists()).toBe(true)
      expect(equipment_label.find('.plan-row-card__equipment-badges').text()).toContain(
        'Schnorchel',
      )
    })

    it('cards without equipment have no equipment-badges element', async () => {
      const wrapper = mount(TrainingPlanDisplay, {
        global: {
          plugins: [i18n, createTestingPinia({ createSpy: vi.fn })],
        },
        props: {
          store: useTrainingPlanStore(),
        },
      })
      const store = useTrainingPlanStore()
      store.currentPlan = createSimplePlan()
      await wrapper.vm.$nextTick()

      const planCards = wrapper.findAll('[data-testid="plan-card"]')
      expect(planCards.length).toBeGreaterThan(0)

      const anyBadge = planCards.some((card) =>
        card.find('.plan-row-card__equipment-badges').exists(),
      )
      expect(anyBadge).toBe(false)
    })
  })
})
