import { apiClient, formatError } from '@/api/client'
import type { UploadedPlan, RAGResponse, Row } from '@/types'
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
      recalculateTotalSum()
      ensureRowIds(currentPlan.value.table)
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
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const tableWithoutIds = currentPlan.value.table.map(({ _id, ...rest }) => rest)

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

  function updatePlanRow(rowIndex: number, field: keyof Row, value: string | number) {
    if (currentPlan.value && currentPlan.value.table[rowIndex]) {
      const row = currentPlan.value.table[rowIndex]
      ;(row[field] as string | number) = value

      if (field === 'Amount' || field === 'Distance') {
        row.Sum = row.Amount * row.Distance
        recalculateTotalSum()
      }
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
        _id: crypto.randomUUID(),
      }
      currentPlan.value.table.splice(rowIndex, 0, newRow)
      recalculateTotalSum()
    }
  }

  function removeRow(rowIndex: number) {
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

    if ((isMovingUp && rowIndex === 0) || (isMovingDown && rowIndex === table.length - 2)) {
      return
    }

    const newIndex = isMovingUp ? rowIndex - 1 : rowIndex + 1
    const [movedRow] = table.splice(rowIndex, 1)
    table.splice(newIndex, 0, movedRow)
  }

  function recalculateTotalSum() {
    if (currentPlan.value && currentPlan.value.table.length > 0) {
      const lastRowIndex = currentPlan.value.table.length - 1
      const lastRow = currentPlan.value.table[lastRowIndex]
      lastRow.Sum = currentPlan.value.table.slice(0, -1).reduce((acc, r) => acc + (r.Sum || 0), 0)
    }
  }

  function ensureRowIds(table: Row[]) {
    table.forEach((row) => {
      if (!row._id) {
        row._id = crypto.randomUUID()
      }
    })
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
    addRow,
    removeRow,
    moveRow,
    clear,
  }
})
