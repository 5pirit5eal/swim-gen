import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { RAGResponse, QueryRequest, Row } from '@/types'

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
      return true
    } else {
      error.value =
        result.error?.message + (': ' + (result.error?.details || 'Unknown')) ||
        'Failed to generate plan'
      isLoading.value = false
      return false
    }
  }

  function updatePlanRow(rowIndex: number, field: keyof Row, value: string | number) {
    console.log(`Updating row ${rowIndex}, field ${field} with value:`, value)
    if (currentPlan.value && currentPlan.value.table[rowIndex]) {
      const row = currentPlan.value.table[rowIndex]
      ;(row[field] as string | number) = value

      // Recalculate Sum if Amount or Distance changed
      if (field === 'Amount' || field === 'Distance') {
        row.Sum = row.Amount * row.Distance

        // Update the last row with the new sum
        if (currentPlan.value.table.length > 0) {
          const lastRowIndex = currentPlan.value.table.length - 1
          const lastRow = currentPlan.value.table[lastRowIndex]
          lastRow.Sum = currentPlan.value.table
            .slice(0, -1)
            .reduce((acc, r) => acc + (r.Sum || 0), 0)
        }
      }
    }
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
    clearPlan,
    clearError,
  }
})
