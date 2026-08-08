import { describe, expect, it } from 'vitest'
import { createBulletExport, createMarkdownExport, type TextExportLabels } from '../textExport'
import type { RAGResponse, Row } from '@/types'

const labels: TextExportLabels = {
  title: 'Text Export',
  description: 'Notes',
  set: 'Set',
  exercise: 'Exercise',
  distance: 'Distance',
  break: 'Break',
  intensity: 'Intensity',
  equipment: 'Equipment',
  total: 'Total',
  meters: 'm',
}

const row = (overrides: Partial<Row>): Row => ({
  Amount: 1,
  Multiplier: 'x',
  Distance: 100,
  Break: '20s',
  Content: 'Freestyle | drill',
  Intensity: 'GA1',
  Sum: 100,
  ...overrides,
})

const plan: RAGResponse = {
  title: 'Morning Plan',
  description: 'Keep it smooth.',
  table: [
    row({ SubRows: [row({ Content: 'Kick', Equipment: ['Fins'] })] }),
    row({ Content: 'Total', Amount: 0, Distance: 0, Break: '', Intensity: '', Sum: 200 }),
  ],
}

const equipment = (item: Row) => item.Equipment?.join(', ') ?? ''

describe('text export formatters', () => {
  it('creates a markdown table with nested rows and escaped cells', () => {
    const text = createMarkdownExport(plan, labels, equipment)

    expect(text).toContain('# Morning Plan')
    expect(text).toContain('Freestyle \\| drill')
    expect(text).toContain('| 1.1 | 1x Kick | 100 m | 20s | GA1 | Fins |')
    expect(text).toContain('**Total:** 200 m')
    expect(text).not.toContain('| Total |')
  })

  it('creates indented compact bullet points', () => {
    const text = createBulletExport(plan, labels, equipment)

    expect(text).toContain('- 1x Freestyle | drill (100 m | 20s | GA1)')
    expect(text).toContain('  - 1x Kick (100 m | 20s | GA1 | Fins)')
    expect(text).toContain('Total: 200 m')
  })
})
