import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { RAGResponse, QueryRequest } from '@/types'

export const useTrainingPlanStore = defineStore('trainingPlan', () => {
  // State
  const currentPlan = ref<RAGResponse | null>(null)
  const isLoading = ref(false)
  const isExporting = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const hasPlan = computed(() => currentPlan.value !== null)
  const isGenerating = computed(() => isLoading.value)

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
      error.value = result.error?.message || 'Failed to generate plan'
      isLoading.value = false
      return false
    }
  }

  function clearPlan() {
    currentPlan.value = null
    error.value = null
  }

  function clearError() {
    error.value = null
  }

  async function exportToPDF(): Promise<boolean> {
    if (!currentPlan.value) {
      error.value = 'No training plan to export'
      return false
    }

    isExporting.value = true
    error.value = null

    try {
      // Import the export store when needed
      const { useExportStore } = await import('@/stores/export')
      const exportStore = useExportStore()
      
      const success = await exportStore.exportToPDF(currentPlan.value)
      
      if (!success && exportStore.exportError) {
        error.value = exportStore.exportError
      }
      
      isExporting.value = false
      return !!success
    } catch (err) {
      error.value = 'Export failed: ' + (err instanceof Error ? err.message : 'Unknown error')
      isExporting.value = false
      return false
    }
  }

  function clearTrainingPlan() {
    clearPlan()
  }

  return {
    // State
    currentPlan,
    isLoading,
    isExporting,
    error,
    // Computed
    hasPlan,
    isGenerating,
    // Actions
    generatePlan,
    clearPlan,
    clearError,
    exportToPDF,
    clearTrainingPlan,
  }
})
