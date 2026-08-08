import { describe, expect, it, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import ButtonExportPlan from '../ButtonExportPlan.vue'
import i18n from '@/plugins/i18n'
import type { PlanStore } from '@/types'

const mockExportToPDF = vi.hoisted(() => vi.fn())
vi.mock('@/stores/export', () => ({
  useExportStore: () => ({ exportToPDF: mockExportToPDF }),
}))

const plan = {
  plan_id: 'plan-1',
  title: 'Test plan',
  description: 'Notes',
  table: [
    {
      Amount: 1,
      Multiplier: 'x',
      Distance: 100,
      Break: '20s',
      Content: 'Warmup',
      Intensity: 'GA1',
      Sum: 100,
    },
    {
      Amount: 0,
      Multiplier: '',
      Distance: 0,
      Break: '',
      Content: 'Total',
      Intensity: '',
      Sum: 100,
    },
  ],
}

const store = { currentPlan: plan, hasPlan: true, isLoading: false } as PlanStore

describe('ButtonExportPlan.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockExportToPDF.mockResolvedValue('https://example.test/plan.pdf')
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText: vi.fn().mockResolvedValue(undefined) },
    })
  })

  it('opens the general export menu in the requested order', async () => {
    const wrapper = mount(ButtonExportPlan, { props: { store }, global: { plugins: [i18n] } })

    await wrapper.find('.main-action').trigger('click')
    const menu = wrapper.find('.dropdown-menu')
    const buttons = menu.findAll('button')

    expect(menu.exists()).toBe(true)
    expect(buttons[0]!.text()).toContain(i18n.global.t('display.export_pdf'))
    expect(menu.findAll('input')).toHaveLength(2)
    expect(buttons[1]!.text()).toContain(i18n.global.t('text_export.export_text'))
  })

  it('opens text export and copies the active representation', async () => {
    const wrapper = mount(ButtonExportPlan, { props: { store }, global: { plugins: [i18n] } })

    await wrapper.find('.main-action').trigger('click')
    await wrapper.find('.dropdown-menu button:last-child').trigger('click')
    await flushPromises()

    expect(wrapper.findComponent({ name: 'TextExportModal' }).exists()).toBe(true)
    await document.body.querySelector<HTMLButtonElement>('.copy-text-button')?.click()
    await flushPromises()

    expect(navigator.clipboard.writeText).toHaveBeenCalledWith(expect.stringContaining('- '))
    expect(document.body.querySelector('.copy-text-button')?.classList.contains('success')).toBe(
      true,
    )
  })
})
