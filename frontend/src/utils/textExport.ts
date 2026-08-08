import type { RAGResponse, Row } from '@/types'

export interface TextExportLabels {
  title: string
  description: string
  set: string
  exercise: string
  distance: string
  break: string
  intensity: string
  equipment: string
  total: string
  meters: string
}

function compactRow(row: Row, equipment: (row: Row) => string): string[] {
  const repetitions = row.Amount ? `${row.Amount}${row.Multiplier || 'x'}` : ''
  return [
    [repetitions, row.Content].filter(Boolean).join(' '),
    row.Distance ? `${row.Distance} m` : '',
    row.Break,
    row.Intensity,
    equipment(row),
  ]
}

function planRows(plan: RAGResponse): { rows: Row[]; total: Row | null } {
  const rows = plan.table ?? []
  return {
    rows: rows.length > 0 ? rows.slice(0, -1) : [],
    total: rows.length > 0 ? (rows[rows.length - 1] ?? null) : null,
  }
}

function flattenRows(rows: Row[], equipment: (row: Row) => string, prefix = ''): string[][] {
  return rows.flatMap((row, index) => {
    const position = prefix ? `${prefix}.${index + 1}` : `${index + 1}`
    const current = [position, ...compactRow(row, equipment)]
    const children = row.SubRows?.length ? flattenRows(row.SubRows, equipment, position) : []
    return [current, ...children]
  })
}

function markdownCell(value: string): string {
  return value.replace(/\|/g, '\\|').replace(/\r?\n/g, '<br>')
}

export function createMarkdownExport(
  plan: RAGResponse,
  labels: TextExportLabels,
  equipment: (row: Row) => string,
): string {
  const { rows, total } = planRows(plan)
  const headers = [labels.set, labels.exercise, labels.distance, labels.break, labels.intensity, labels.equipment]
  const lines = [`# ${plan.title}`, '', plan.description, '', `| ${headers.join(' | ')} |`, `| ${headers.map(() => '---').join(' | ')} |`]

  for (const row of flattenRows(rows, equipment)) {
    lines.push(`| ${row.map(markdownCell).join(' | ')} |`)
  }

  if (total) lines.push('', `**${labels.total}:** ${total.Sum} ${labels.meters}`)
  return lines.join('\n')
}

function bulletRows(rows: Row[], equipment: (row: Row) => string, depth = 0): string[] {
  return rows.flatMap((row) => {
    const values = compactRow(row, equipment)
    const details = values.slice(1).filter(Boolean).join(' | ')
    const line = `${'  '.repeat(depth)}- ${values[0]}${details ? ` (${details})` : ''}`
    const children = row.SubRows?.length ? bulletRows(row.SubRows, equipment, depth + 1) : []
    return [line, ...children]
  })
}

export function createBulletExport(
  plan: RAGResponse,
  labels: TextExportLabels,
  equipment: (row: Row) => string,
): string {
  const { rows, total } = planRows(plan)
  const lines = [`${plan.title}`, '', plan.description, '', ...bulletRows(rows, equipment)]
  if (total) lines.push('', `${labels.total}: ${total.Sum} ${labels.meters}`)
  return lines.join('\n')
}
