import { ref } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { PlanToPDFRequest } from '@/types'

export const useExportStore = defineStore('export', () => {
  // State
  const isExporting = ref(false)
  const exportError = ref<string | null>(null)

  // Actions
  async function exportToPDF(request: PlanToPDFRequest): Promise<string | null> {
    console.debug('[ExportStore] exportToPDF', { request })
    if (request.table.length <= 1) {
      exportError.value = 'Cannot export an empty plan.'
      return null
    }

    isExporting.value = true
    exportError.value = null

    const result = await apiClient.exportPDF(request)

    if (result.success && result.data) {
      isExporting.value = false
      return result.data.uri // Return the PDF URL
    } else {
      exportError.value = result.error?.message || 'Failed to export PDF'
      isExporting.value = false
      return null
    }
  }

  function clearError() {
    exportError.value = null
  }

  return {
    // State
    isExporting,
    exportError,
    // Actions
    exportToPDF,
    clearError,
  }
})
