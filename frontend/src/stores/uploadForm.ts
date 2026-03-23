import { apiClient, formatError } from '@/api/client'
import i18n from '@/plugins/i18n'
import type { DonatePlanRequest, RAGResponse, Row } from '@/types'
import {
  stripRowIds,
  updateRowField,
  addRowAtPath,
  addSubRow as addSubRowHelper,
  removeRowAtPath,
  moveRowAtPath,
  recalculateAllSums,
  updateRowEquipment,
} from '@/utils/rowHelpers'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useUploadStore } from '@/stores/uploads'

export const useUploadFormStore = defineStore('uploadForm', () => {
  // --- STATE ---
  const currentPlan = ref<RAGResponse | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const donationStore = useUploadStore()

  // --- COMPUTED ---
  const hasPlan = computed(() => currentPlan.value !== null)

  // --- ACTIONS ---

  // Initialize a new empty plan for donation
  function initNewPlan() {
    currentPlan.value = {
      title: i18n.global.t('donation.newPlan.title'),
      description: i18n.global.t('donation.newPlan.description'),
      table: [
        {
          Amount: 1,
          Multiplier: 'x',
          Distance: 100,
          Break: '20',
          Content: i18n.global.t('donation.newPlan.warmup'),
          Intensity: 'GA1',
          Sum: 100,
          _id: crypto.randomUUID(),
        },
        {
          Amount: 0,
          Multiplier: '',
          Distance: 0,
          Break: '',
          Content: i18n.global.t('donation.newPlan.total'),
          Intensity: '',
          Sum: 100,
          _id: crypto.randomUUID(),
        },
      ],
    }
  }

  // Submit the current plan as a donation
  async function uploadCurrentPlan(allowSharing: boolean = false): Promise<boolean> {
    if (!currentPlan.value) return false
    isLoading.value = true
    error.value = null

    // Strip _id from table rows before sending
    const tableWithoutIds = stripRowIds(currentPlan.value.table.slice(0, -1))

    const request: DonatePlanRequest = {
      title: currentPlan.value.title,
      description: currentPlan.value.description,
      table: tableWithoutIds,
      language: i18n.global.locale.value,
      allow_sharing: allowSharing,
    }

    const result = await apiClient.donatePlan(request)
    if (result.success) {
      // Refresh the list of uploaded plans in the donation store
      await donationStore.fetchUploadedPlans()
      isLoading.value = false
      return true
    } else {
      error.value = result.error ? formatError(result.error) : 'Failed to upload plan'
      isLoading.value = false
      return false
    }
  }

  // --- PlanStore Implementation ---

  async function keepForever(planId: string) {
    // Not applicable for upload form
    console.log('keepForever not implemented for upload form', planId)
  }

  async function upsertCurrentPlan(): Promise<string> {
    // Not applicable for upload form in the same way,
    // we don't auto-save to backend on every edit, only on explicit upload
    // But we can use this to trigger local validation or updates if needed
    if (currentPlan.value) {
      recalculateAllSums(currentPlan.value.table)
    }
    return '' // Upload form doesn't create plans in history
  }

  function updatePlanRow(path: number[], field: keyof Row, value: string | number) {
    if (!currentPlan.value) return
    updateRowField(currentPlan.value.table, path, field, value)
  }

  function updatePlanRowEquipment(path: number[], equipment: string[]) {
    if (!currentPlan.value) return
    updateRowEquipment(currentPlan.value.table, path, equipment)
  }

  function addRow(path: number[]) {
    if (!currentPlan.value) return
    addRowAtPath(currentPlan.value.table, path)
  }

  function addSubRow(path: number[], depth: number) {
    if (!currentPlan.value) return
    addSubRowHelper(currentPlan.value.table, path, depth)
  }

  function removeRow(path: number[]) {
    if (!currentPlan.value) return
    removeRowAtPath(currentPlan.value.table, path)
  }

  function moveRow(path: number[], direction: 'up' | 'down') {
    if (!currentPlan.value) return
    moveRowAtPath(currentPlan.value.table, path, direction)
  }

  function clear() {
    currentPlan.value = null
    error.value = null
    isLoading.value = false
  }

  return {
    // State
    currentPlan,
    isLoading,
    error,
    // Computed
    hasPlan,
    // Actions
    initNewPlan,
    uploadCurrentPlan,
    // PlanStore Implementation
    keepForever,
    upsertCurrentPlan,
    updatePlanRow,
    updatePlanRowEquipment,
    addRow,
    addSubRow,
    removeRow,
    moveRow,
    clear,
  }
})
