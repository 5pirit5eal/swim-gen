import { apiClient, formatError } from '@/api/client'
import i18n from '@/plugins/i18n'
import type { QueryRequest, RAGResponse, Row } from '@/types'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useTrainingPlanStore = defineStore('trainingPlan', () => {
  // State
  const currentPlan = ref<RAGResponse | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const hasPlan = computed(() => currentPlan.value !== null)

  // Actions
  async function generatePlan(request: QueryRequest): Promise<boolean> {
    isLoading.value = true
    error.value = null

    const result = await apiClient.query(request)

    if (result.success && result.data) {
      currentPlan.value = result.data
      isLoading.value = false
      recalculateTotalSum()
      return true
    } else {
      error.value = result.error
        ? formatError(result.error)
        : i18n.global.t('errors.training_plan_failed')
      isLoading.value = false
      return false
    }
  }

  function updatePlanRow(rowIndex: number, field: keyof Row, value: string | number) {
    console.log(`Updating row ${rowIndex}, field ${field} with value:`, value)
    if (currentPlan.value && currentPlan.value.table[rowIndex]) {
      const row = currentPlan.value.table[rowIndex]
        ; (row[field] as string | number) = value

      // Recalculate Sum if Amount or Distance changed
      if (field === 'Amount' || field === 'Distance') {
        row.Sum = row.Amount * row.Distance
        recalculateTotalSum()
      }
    }
  }

  function recalculateTotalSum() {
    if (currentPlan.value && currentPlan.value.table.length > 0) {
      const lastRowIndex = currentPlan.value.table.length - 1
      const lastRow = currentPlan.value.table[lastRowIndex]
      lastRow.Sum = currentPlan.value.table.slice(0, -1).reduce((acc, r) => acc + (r.Sum || 0), 0)
    }
  }

  function addRow(rowIndex: number) {
    if (currentPlan.value && currentPlan.value.table.length < 26) {
      const newRow: Row = {
        Amount: 0,
        Break: '',
        Content: '',
        Distance: 0,
        Intensity: '',
        Multiplier: 'x',
        Sum: 0,
      }
      currentPlan.value.table.splice(rowIndex, 0, newRow)
      recalculateTotalSum()
    }
  }

  function removeRow(rowIndex: number) {
    // Ensure we don't remove the total row and at least one exercise row remains
    if (
      currentPlan.value &&
      currentPlan.value.table.length > 2 &&
      rowIndex < currentPlan.value.table.length - 1
    ) {
      currentPlan.value.table.splice(rowIndex, 1)
      recalculateTotalSum()
    }
  }

  function moveRow(rowIndex: number, direction: 'up' | 'down') {
    if (!currentPlan.value) return

    const table = currentPlan.value.table
    const isMovingUp = direction === 'up'
    const isMovingDown = direction === 'down'

    // Prevent moving the first row up or the last exercise row down
    if ((isMovingUp && rowIndex === 0) || (isMovingDown && rowIndex === table.length - 2)) {
      return
    }

    const newIndex = isMovingUp ? rowIndex - 1 : rowIndex + 1
    const [movedRow] = table.splice(rowIndex, 1)
    table.splice(newIndex, 0, movedRow)
  }

  function clearPlan() {
    currentPlan.value = null
    error.value = null
  }

  function clearError() {
    error.value = null
  }

  return {
    // State
    currentPlan,
    isLoading,
    error,
    // Computed
    hasPlan,
    // Actions
    generatePlan,
    updatePlanRow,
    addRow,
    removeRow,
    moveRow,
    clearPlan,
    clearError,
  }
})
