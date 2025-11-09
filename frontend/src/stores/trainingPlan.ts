import { apiClient, formatError } from '@/api/client'
import i18n from '@/plugins/i18n'
import type { QueryRequest, RAGResponse, Row, UpsertPlanRequest, UpsertPlanResponse } from '@/types'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'

export const useTrainingPlanStore = defineStore('trainingPlan', () => {
  // --- STATE ---
  const currentPlan = ref<RAGResponse | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const generationHistory = ref<RAGResponse[]>([])
  const userStore = useAuthStore()

  // --- COMPUTED ---
  const hasPlan = computed(() => currentPlan.value !== null)

  // --- ACTIONS ---

  // Fetches the user's plan history
  async function fetchHistory() {
    userStore.getUser()
    if (!userStore.user) {
      console.log('User is not available.')
      return
    }
    isLoading.value = true
    const { data, error } = await supabase
      .from('history')
      .select('plan_id')
      .order('created_at', { ascending: false })

    if (error) {
      console.error(error)
    } else if (data) {
      const planIds = data.map((entry) => entry.plan_id)
      const { data: plansData, error: plansError } = await supabase
        .from('plans')
        .select('plan_id, title, description, table')
        .in('plan_id', planIds)
      if (plansError) {
        console.error(plansError)
      } else if (plansData) {
        generationHistory.value = plansData.map((plan) => ({
          plan_id: plan.plan_id,
          title: plan.title,
          description: plan.description,
          table: plan.table,
        }))
      }
    }
    isLoading.value = false
  }

  // Generates a new training plan
  async function generatePlan(request: QueryRequest): Promise<boolean> {
    isLoading.value = true
    error.value = null

    const result = await apiClient.query(request)

    if (result.success && result.data) {
      currentPlan.value = result.data
      recalculateTotalSum()
      await fetchHistory() // Refresh history after generating a new plan
      isLoading.value = false
      return true
    } else {
      error.value = result.error
        ? formatError(result.error)
        : i18n.global.t('errors.training_plan_failed')
      isLoading.value = false
      return false
    }
  }

  // Upserts a plan
  async function upsertPlan(plan: UpsertPlanRequest): Promise<UpsertPlanResponse | null> {
    userStore.getUser()
    if (!userStore.user) {
      console.log('User is not available.')
      return null
    }
    isLoading.value = true
    const result = await apiClient.upsertPlan(plan)
    isLoading.value = false
    if (result.success && result.data) {
      await fetchHistory() // Refresh history after upserting
      return result.data
    } else {
      console.error(result.error ? formatError(result.error) : 'Unknown error during upsertPlan')
      return null
    }
  }

  // Loads a plan from history into the editor
  function loadPlanFromHistory(plan: RAGResponse) {
    currentPlan.value = JSON.parse(JSON.stringify(plan)) // Deep copy to prevent accidental edits
  }

  // Sets a plan to be remembered forever
  async function keepPlanForever(planId: string) {
    userStore.getUser()
    if (!userStore.user) {
      console.log('User is not available.')
      return null
    }
    const { error } = await supabase
      .from('history')
      .update({ keep_forever: true })
      .eq('plan_id', planId)
    if (error) {
      console.error(error)
    }

    await fetchHistory()
  }

  // --- Plan Table Manipulations ---

  function updatePlanRow(rowIndex: number, field: keyof Row, value: string | number) {
    if (currentPlan.value && currentPlan.value.table[rowIndex]) {
      const row = currentPlan.value.table[rowIndex]
        ; (row[field] as string | number) = value

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
    generationHistory,
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
    fetchHistory,
    upsertPlan,
    loadPlanFromHistory,
    keepPlanForever,
  }
})
