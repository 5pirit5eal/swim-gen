import { apiClient, formatError } from '@/api/client'
import type { UploadedPlan, RAGResponse, Row } from '@/types'
import {
  ensureRowIds,
  normalizeRows,
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
import { computed, ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const useUploadStore = defineStore('upload', () => {
  // --- STATE ---
  const currentPlan = ref<RAGResponse | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const uploadedPlans = ref<UploadedPlan[]>([])
  const allUploadedPlans = ref<UploadedPlan[]>([])
  const isFetchingUploads = ref(false)
  const userStore = useAuthStore()

  // Pagination state
  const PAGE_SIZE = 10
  const displayCount = ref(PAGE_SIZE)
  const historyHasMore = computed(() => displayCount.value < allUploadedPlans.value.length)
  const isLoadingMore = ref(false)

  // --- COMPUTED ---
  const hasPlan = computed(() => currentPlan.value !== null)

  watch(
    () => userStore.user?.id ?? null,
    async (userId) => {
      if (userId) {
        await fetchUploadedPlans()
      } else {
        uploadedPlans.value = []
      }
    },
    { immediate: true },
  )

  // --- ACTIONS ---

  // Fetch all uploaded plans for the user
  async function fetchUploadedPlans(reset = true) {
    console.debug('[UploadStore] fetchUploadedPlans', { reset })
    if (!userStore.user) return
    if (reset) {
      displayCount.value = PAGE_SIZE
    }
    isFetchingUploads.value = true
    const result = await apiClient.getUploadedPlans()
    if (result.success && Array.isArray(result.data)) {
      allUploadedPlans.value = result.data.sort((a, b) => {
        return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      })
      // Display only the first PAGE_SIZE items initially
      uploadedPlans.value = allUploadedPlans.value.slice(0, displayCount.value)
    } else {
      console.error(result.error ? formatError(result.error) : 'Failed to fetch uploaded plans')
    }
    isFetchingUploads.value = false
  }

  // Fetches more uploaded plans (pagination - client side)
  function fetchMoreUploadedPlans() {
    if (!historyHasMore.value || isLoadingMore.value) return
    isLoadingMore.value = true
    displayCount.value += PAGE_SIZE
    uploadedPlans.value = allUploadedPlans.value.slice(0, displayCount.value)
    isLoadingMore.value = false
  }

  // Fetch a specific uploaded plan
  async function fetchUploadedPlan(planId: string): Promise<boolean> {
    console.debug('[UploadStore] fetchUploadedPlan', { planId })
    isLoading.value = true
    error.value = null
    const result = await apiClient.getUploadedPlan(planId)
    if (result.success && result.data) {
      // Convert DonatedPlan to RAGResponse format for display
      currentPlan.value = {
        plan_id: result.data.plan_id,
        title: result.data.title,
        description: result.data.description,
        table: result.data.table,
      }
      recalculateAllSums(currentPlan.value.table)
      ensureRowIds(currentPlan.value.table)
      normalizeRows(currentPlan.value.table)
      isLoading.value = false
      return true
    } else {
      error.value = result.error ? formatError(result.error) : 'Failed to fetch uploaded plan'
      isLoading.value = false
      return false
    }
  }

  // Loads a plan from history into the editor
  async function loadPlanFromHistory(plan_id: string) {
    if (!userStore.user) return
    if (currentPlan.value?.plan_id === plan_id) return
    await fetchUploadedPlan(plan_id)
    if (currentPlan.value) {
      ensureRowIds(currentPlan.value.table)
      normalizeRows(currentPlan.value.table)
    }
  }

  // --- PlanStore Implementation ---

  async function keepForever(planId: string) {
    // Not applicable for donation store, but required by interface
    console.log('keepForever not implemented for donation store', planId)
  }

  async function upsertCurrentPlan(): Promise<string> {
    console.debug('[UploadStore] upsertCurrentPlan', { planId: currentPlan.value?.plan_id })
    if (!userStore.user) throw new Error('User is not available')
    if (!currentPlan.value) throw new Error('No current plan to upsert')

    // Strip _id from table rows before sending to backend
    const tableWithoutIds = stripRowIds(currentPlan.value.table)

    const result = await apiClient.upsertPlan({
      plan_id: currentPlan.value.plan_id,
      title: currentPlan.value.title,
      description: currentPlan.value.description,
      table: tableWithoutIds,
    })
    if (result.success && result.data) {
      await fetchUploadedPlans() // Refresh plans after upserting
      return result.data.plan_id
    } else {
      console.error(result.error ? formatError(result.error) : 'Unknown error during upsertPlan')
      throw new Error(result.error ? formatError(result.error) : 'Upsert failed')
    }
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
    uploadedPlans,
    isFetchingUploads,
    // Pagination state
    historyHasMore,
    isLoadingMore,
    // Computed
    hasPlan,
    // Actions
    fetchUploadedPlans,
    fetchMoreUploadedPlans,
    fetchUploadedPlan,
    loadPlanFromHistory,
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
