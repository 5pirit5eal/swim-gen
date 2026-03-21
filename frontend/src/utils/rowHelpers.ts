import type { Row } from '@/types'

/** Matches backend validateRowDepth max depth (plan.go). */
export const MAX_NESTING_DEPTH = 4

/** Matches backend EquipmentType constants (metadata.go). */
export const EQUIPMENT_TYPES = [
  'Flossen',
  'Kickboard',
  'Handpaddles',
  'Pull buoy',
  'Schnorchel',
] as const

export type EquipmentType = (typeof EQUIPMENT_TYPES)[number]

export function createEmptyRow(): Row {
  return {
    Amount: 0,
    Break: '',
    Content: '',
    Distance: 0,
    Intensity: '',
    Multiplier: 'x',
    Sum: 0,
    Equipment: [],
    SubRows: [],
    _id: crypto.randomUUID(),
  }
}

/** Recursively assigns _id to all rows and their SubRows. Mutates in place. */
export function ensureRowIds(rows: Row[]): void {
  for (const row of rows) {
    if (!row._id) {
      row._id = crypto.randomUUID()
    }
    if (row.SubRows && row.SubRows.length > 0) {
      ensureRowIds(row.SubRows)
    }
  }
}

/** Ensures SubRows and Equipment are always arrays (never undefined). Mutates in place. */
export function normalizeRows(rows: Row[]): void {
  for (const row of rows) {
    if (!row.SubRows) {
      row.SubRows = []
    }
    if (!row.Equipment) {
      row.Equipment = []
    }
    if (row.SubRows.length > 0) {
      normalizeRows(row.SubRows)
    }
  }
}

/** Recursively strips _id from all rows. Returns new array (pure). */
export function stripRowIds(rows: Row[]): Row[] {
  return rows.map((row) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { _id, ...rest } = row
    return {
      ...rest,
      SubRows: rest.SubRows && rest.SubRows.length > 0 ? stripRowIds(rest.SubRows) : [],
      Equipment: rest.Equipment ?? [],
    }
  })
}

export interface RowPathResult {
  row: Row
  parent: Row[]
  index: number
}

/**
 * Navigates to a row via index path.
 * Path example: [2] = top-level row 2, [2, 1] = row 2's SubRow 1.
 */
export function findRowByPath(table: Row[], path: number[]): RowPathResult | null {
  if (path.length === 0) return null

  let currentArray = table
  for (let i = 0; i < path.length - 1; i++) {
    const idx = path[i]!
    const row = currentArray[idx]
    if (!row || !row.SubRows) return null
    currentArray = row.SubRows
  }

  const lastIdx = path[path.length - 1]!
  const row = currentArray[lastIdx]
  if (!row) return null

  return { row, parent: currentArray, index: lastIdx }
}

/**
 * Depth-first sum recalculation matching backend UpdateSum() (plan.go):
 * - Leaf: Sum = Amount * Distance
 * - Parent: Distance = sum(child.Sum), then Sum = Amount * Distance
 */
export function recalculateRowSum(row: Row): void {
  if (row.SubRows && row.SubRows.length > 0) {
    for (const subRow of row.SubRows) {
      recalculateRowSum(subRow)
    }
    row.Distance = row.SubRows.reduce((acc, sr) => acc + (sr.Sum || 0), 0)
  }
  row.Sum = row.Amount * row.Distance
}

/** Total row (last in table) Sum = sum of all other top-level rows' Sum. */
export function recalculateTotalSum(table: Row[]): void {
  if (table.length === 0) return
  const lastRow = table[table.length - 1]!
  lastRow.Sum = table.slice(0, -1).reduce((acc, r) => acc + (r.Sum || 0), 0)
}

/** Recalculates all exercise row sums (depth-first) then the total row. */
export function recalculateAllSums(table: Row[]): void {
  for (let i = 0; i < table.length - 1; i++) {
    recalculateRowSum(table[i]!)
  }
  recalculateTotalSum(table)
}

/** Updates a field on the row at path, recalculates sums if Amount/Distance changed. Mutates table. */
export function updateRowField(
  table: Row[],
  path: number[],
  field: keyof Row,
  value: string | number,
): void {
  const result = findRowByPath(table, path)
  if (!result) return

  const { row } = result
  ;(row[field] as string | number) = value

  if (field === 'Amount' || field === 'Distance') {
    recalculateAllSums(table)
  }
}

export function updateRowEquipment(table: Row[], path: number[], equipment: string[]): void {
  const result = findRowByPath(table, path)
  if (!result) return
  result.row.Equipment = equipment
}

/**
 * Inserts a new empty row before the position specified by path.
 * Top-level: max 25 exercise rows + 1 total = 26.
 */
export function addRowAtPath(table: Row[], path: number[]): void {
  if (path.length === 0) return

  const parentPath = path.slice(0, -1)
  const insertIndex = path[path.length - 1]!

  let targetArray: Row[]
  if (parentPath.length === 0) {
    targetArray = table
    if (targetArray.length >= 26) return
  } else {
    const parentResult = findRowByPath(table, parentPath)
    if (!parentResult) return
    if (!parentResult.row.SubRows) {
      parentResult.row.SubRows = []
    }
    targetArray = parentResult.row.SubRows
  }

  targetArray.splice(insertIndex, 0, createEmptyRow())
  recalculateAllSums(table)
}

/** Appends a SubRow to the row at path. Enforces MAX_NESTING_DEPTH. */
export function addSubRow(table: Row[], path: number[], depth: number): void {
  if (depth >= MAX_NESTING_DEPTH) return

  const result = findRowByPath(table, path)
  if (!result) return

  const { row } = result
  if (!row.SubRows) {
    row.SubRows = []
  }

  row.SubRows.push(createEmptyRow())

  if (row.Amount === 0) {
    row.Amount = 1
  }

  recalculateAllSums(table)
}

/**
 * Removes row at path.
 * Top-level: requires >= 2 rows (1 exercise + 1 total), cannot remove total row.
 */
export function removeRowAtPath(table: Row[], path: number[]): void {
  if (path.length === 0) return

  const parentPath = path.slice(0, -1)
  const removeIndex = path[path.length - 1]!

  let targetArray: Row[]
  if (parentPath.length === 0) {
    targetArray = table
    if (targetArray.length <= 2) return
    if (removeIndex >= targetArray.length - 1) return
  } else {
    const parentResult = findRowByPath(table, parentPath)
    if (!parentResult) return
    if (!parentResult.row.SubRows) return
    targetArray = parentResult.row.SubRows
  }

  targetArray.splice(removeIndex, 1)
  recalculateAllSums(table)
}

/** Moves row at path up or down within its parent array. */
export function moveRowAtPath(
  table: Row[],
  path: number[],
  direction: 'up' | 'down',
): void {
  if (path.length === 0) return

  const parentPath = path.slice(0, -1)
  const rowIndex = path[path.length - 1]!

  let targetArray: Row[]
  let lastMovableIndex: number

  if (parentPath.length === 0) {
    targetArray = table
    lastMovableIndex = targetArray.length - 2
  } else {
    const parentResult = findRowByPath(table, parentPath)
    if (!parentResult || !parentResult.row.SubRows) return
    targetArray = parentResult.row.SubRows
    lastMovableIndex = targetArray.length - 1
  }

  const isMovingUp = direction === 'up'
  const isMovingDown = direction === 'down'

  if ((isMovingUp && rowIndex === 0) || (isMovingDown && rowIndex >= lastMovableIndex)) {
    return
  }

  const newIndex = isMovingUp ? rowIndex - 1 : rowIndex + 1
  const [movedRow] = targetArray.splice(rowIndex, 1)
  targetArray.splice(newIndex, 0, movedRow!)
}

/**
 * Classification helpers for display-model (TDD Wave 1).
 * These pure functions derive UI-friendly metrics without modifying Row type.
 */

export function isExerciseRow(row: Row): boolean {
  return row.Content !== 'Total'
}

export function isTotalRow(table: Row[], row: Row): boolean {
  return table.length > 0 && row === table[table.length - 1]
}

export function isParentRow(row: Row): boolean {
  return !!(row.SubRows && row.SubRows.length > 0)
}

export function isLeafRow(row: Row): boolean {
  return !row.SubRows || row.SubRows.length === 0
}

export interface DisplayRowMetrics {
  isExercise: boolean
  isTotal: boolean
  isParent: boolean
  isLeaf: boolean
  depth: number
}

export function getDisplayRowMetrics(table: Row[], row: Row, path: number[]): DisplayRowMetrics {
  return {
    isExercise: isExerciseRow(row),
    isTotal: isTotalRow(table, row),
    isParent: isParentRow(row),
    isLeaf: isLeafRow(row),
    depth: path.length,
  }
}
