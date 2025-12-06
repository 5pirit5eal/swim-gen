import { apiClient, formatError } from '@/api/client'
import i18n from '@/plugins/i18n'
import type { SharedPlanData, SharedHistoryItem, Row, RAGResponse } from '@/types'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import router from '@/router'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import { useUploadStore } from '@/stores/uploads'

export const useSharedPlanStore = defineStore('sharedPlan', () => {
  const authStore = useAuthStore()

  // --- STATE ---
  const sharedPlan = ref<SharedPlanData | null>(null)
  const sharedHistory = ref<SharedHistoryItem[]>([])
  const isLoading = ref(false)
  const isFetchingHistory = ref(false)
  const error = ref<string | null>(null)
  const isForked = ref(false)

  // Pagination state
  const PAGE_SIZE = 5
  const historyPage = ref(0)
  const historyHasMore = ref(true)
  const isLoadingMore = ref(false)

  // --- COMPUTED ---
  const currentPlan = computed(() => sharedPlan.value?.plan || null)
  const hasPlan = computed(() => currentPlan.value !== null)

  watch(
    () => authStore.user?.id ?? null,
    async (newUserId) => {
      if (newUserId) {
        await fetchSharedHistory()
      }
    },
    { immediate: true },
  )

  // --- ACTIONS ---

  // Required by the PlanStore interface, but not applicable for shared plans
  async function keepForever() {
    // No-op
    return
  }

  // Fetches a shared plan by its hash
  async function fetchSharedPlanByHash(hash: string): Promise<string | null> {
    isLoading.value = true
    error.value = null
    isForked.value = false

    // Check if we are loading a plan from history
    if (hash === '' && sharedPlan.value !== null) {
      isLoading.value = false
      return null
    } else if (hash === '') {
      error.value = i18n.global.t('errors.fetch_shared_plan_failed')
      isLoading.value = false
      return null
    }

    try {
      // 1. Get the plan ID and sharer ID from the shared_plans table using the hash
      const { data: sharedPlanData, error: sharedPlanError } = await supabase
        .from('shared_plans')
        .select('plan_id, user_id')
        .eq('url_hash', hash)
        .single()

      if (sharedPlanError) throw sharedPlanError
      if (!sharedPlanData) throw new Error('Plan not found')

      // Check if the plan is already loaded
      if (
        sharedPlan.value &&
        sharedPlan.value.plan.plan_id === sharedPlanData.plan_id &&
        sharedPlan.value.sharer_id === sharedPlanData.user_id
      ) {
        isLoading.value = false
        return null
      }

      // Check if the user is trying to load their own shared plan
      if (authStore.user && sharedPlanData.user_id === authStore.user.id) {
        const trainingPlanStore = useTrainingPlanStore()
        if (!trainingPlanStore.planHistory.length) {
          if (!trainingPlanStore.isFetchingHistory) {
            await trainingPlanStore.fetchHistory()
          } else {
            while (trainingPlanStore.isFetchingHistory) {
              await new Promise((resolve) => setTimeout(resolve, 100))
            }
          }
        }
        const ownPlan = trainingPlanStore.planHistory.find(
          (plan) => plan.plan_id === sharedPlanData.plan_id,
        )
        if (ownPlan) {
          trainingPlanStore.loadPlanFromHistory(ownPlan)
          isLoading.value = false
          router.push({ name: 'plan', params: { id: ownPlan.plan_id } })
          return 'own_plan'
        }

        // Check if it is an uploaded plan
        const uploadStore = useUploadStore()
        if (!uploadStore.uploadedPlans.length) {
          if (!uploadStore.isFetchingUploads) {
            await uploadStore.fetchUploadedPlans()
          } else {
            while (uploadStore.isFetchingUploads) {
              await new Promise((resolve) => setTimeout(resolve, 100))
            }
          }
        }
        const ownUploadedPlan = uploadStore.uploadedPlans.find(
          (plan) => plan.plan_id === sharedPlanData.plan_id,
        )
        if (ownUploadedPlan) {
          isLoading.value = false
          router.push({ name: 'uploaded', params: { planId: ownUploadedPlan.plan_id } })
          return 'own_plan'
        }
      }

      // 2. Fetch the plan details
      const { data: planData, error: planError } = await supabase
        .from('plans')
        .select('plan_id, title, description, plan_table')
        .eq('plan_id', sharedPlanData.plan_id)
        .single()

      if (planError) throw planError
      if (!planData) throw new Error('Plan details not found')

      // 3. Fetch the sharer's username (if available)
      let sharerUsername = i18n.global.t('shared.unknown_user')
      if (sharedPlanData.user_id) {
        const { data: profileData, error: profileError } = await supabase
          .from('profiles')
          .select('username')
          .eq('user_id', sharedPlanData.user_id)
          .single()

        if (!profileError && profileData) {
          sharerUsername = profileData.username || i18n.global.t('shared.unknown_user')
        }
      }

      sharedPlan.value = {
        plan: {
          plan_id: planData.plan_id,
          title: planData.title,
          description: planData.description,
          table: planData.plan_table,
        },
        sharer_username: sharerUsername,
        sharer_id: sharedPlanData.user_id,
      }
      if (sharedPlan.value?.plan) {
        ensureRowIds(sharedPlan.value.plan.table)
      }

      // Add to history if user is logged in and the plan is not their own
      if (authStore.user) {
        await addPlanToHistory(sharedPlanData.plan_id, sharedPlanData.user_id)
        await fetchSharedHistory()
      }
    } catch (e) {
      console.error(e)
      error.value = i18n.global.t('errors.fetch_shared_plan_failed')
    } finally {
      isLoading.value = false
    }
    return null
  }

  // Fetches the user's shared plan history with pagination
  async function fetchSharedHistory(reset = true) {
    if (!authStore.user) return

    if (reset) {
      historyPage.value = 0
      historyHasMore.value = true
      sharedHistory.value = []
    }

    isFetchingHistory.value = true
    try {
      const offset = historyPage.value * PAGE_SIZE
      const { data: historyData, error: historyError } = await supabase
        .from('shared_history')
        .select('user_id, plan_id, share_method, shared_by, created_at')
        .eq('user_id', authStore.user.id)
        .order('created_at', { ascending: false })
        .range(offset, offset + PAGE_SIZE - 1)

      if (historyError) throw historyError

      if (historyData) {
        // Check if there are more results
        historyHasMore.value = historyData.length === PAGE_SIZE

        const planIds = historyData.map((item) => item.plan_id)

        // Fetch all plan details in one query
        const { data: plansData } = await supabase
          .from('plans')
          .select('plan_id, title, description, plan_table')
          .in('plan_id', planIds)

        if (plansData === null) {
          if (reset) sharedHistory.value = []
          console.info('No plans data found for shared history')
          return
        }

        // Create a map for easier lookup
        const plansMap = new Map<
          string,
          { plan_id: string; title: string; description: string; plan_table: Row[] }
        >()
        if (plansData) {
          plansData.forEach((plan) => {
            plansMap.set(plan.plan_id, plan)
          })
        }

        // Combine history items with plan details
        const rawSharedPlanHistory = historyData.map((item) => {
          const planData = plansMap.get(item.plan_id)
          if (!planData) {
            console.warn(`Plan data not found for plan_id: ${item.plan_id}`)
            return undefined
          }
          return {
            user_id: item.user_id,
            plan_id: item.plan_id,
            share_method: item.share_method,
            shared_by: item.shared_by,
            created_at: item.created_at,
            plan: {
              plan_id: planData.plan_id,
              title: planData.title,
              description: planData.description,
              table: planData.plan_table,
            } as RAGResponse,
          }
        })
        // Filter out any undefined entries due to missing plan data
        const newItems = rawSharedPlanHistory.filter(
          (item): item is SharedHistoryItem => item !== undefined,
        )
        if (reset) {
          sharedHistory.value = newItems
        } else {
          sharedHistory.value = [...sharedHistory.value, ...newItems]
        }
      }
    } catch (e) {
      console.error(e)
      // Don't set global error for history fetch failure to avoid blocking UI
    } finally {
      isFetchingHistory.value = false
    }
  }

  // Fetches more shared history entries (pagination)
  async function fetchMoreSharedHistory() {
    if (!historyHasMore.value || isLoadingMore.value) return
    isLoadingMore.value = true
    historyPage.value += 1
    await fetchSharedHistory(false)
    isLoadingMore.value = false
  }

  // Load plan from history
  async function loadPlanFromHistory(item: SharedHistoryItem) {
    isLoading.value = true
    error.value = null
    sharedPlan.value = null
    isForked.value = false

    try {
      // 1. Get the username from the profiles table using the shared_by ID
      const { data: profileData, error: profileError } = await supabase
        .from('profiles')
        .select('username')
        .eq('user_id', item.shared_by)
        .single()

      if (profileError) throw profileError
      if (!profileData) throw new Error('User not found')
      if (!item.plan) throw new Error('Plan details not found')

      sharedPlan.value = {
        plan: item.plan,
        sharer_username: profileData.username || i18n.global.t('shared.unknown_user'),
        sharer_id: item.user_id,
      }
      if (sharedPlan.value?.plan) {
        ensureRowIds(sharedPlan.value.plan.table)
      }
    } catch (e) {
      console.error(e)
      error.value = i18n.global.t('errors.fetch_shared_plan_failed')
    } finally {
      isLoading.value = false
    }
  }

  // Adds a plan to the user's shared history
  async function addPlanToHistory(planId: string, sharedBy: string) {
    if (!authStore.user) return

    try {
      // Check if already exists to avoid duplicates (or could use upsert)
      const { data: existing } = await supabase
        .from('shared_history')
        .select('plan_id')
        .eq('user_id', authStore.user.id)
        .eq('plan_id', planId)
        .maybeSingle()

      if (!existing) {
        await supabase.from('shared_history').insert({
          user_id: authStore.user.id,
          plan_id: planId,
          shared_by: sharedBy,
          share_method: 'link', // Default for now
        })
      }
    } catch (e) {
      console.error('Failed to add to history', e)
    }
  }

  function clear() {
    sharedPlan.value = null
    error.value = null
    isForked.value = false
  }

  // --- Plan Table Manipulations (Local Only) ---

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

  // Upserts the current plan. If it's the first edit (not forked yet),
  // it saves as a new plan (forking) to the user's history.
  // If already forked, updates the forked plan.
  async function upsertCurrentPlan(): Promise<string> {
    if (!authStore.user || !currentPlan.value) {
      throw new Error('User or plan not available')
    }

    // If not yet forked, strip the plan_id to create a new plan
    const planIdToUse = isForked.value ? currentPlan.value.plan_id : undefined

    // Strip _id from table rows before sending to backend
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const tableWithoutIds = currentPlan.value.table.map(({ _id, ...rest }) => rest)

    const result = await apiClient.upsertPlan({
      plan_id: planIdToUse,
      title: currentPlan.value.title,
      description: currentPlan.value.description,
      table: tableWithoutIds,
    })

    if (result.success && result.data) {
      // Mark as forked and refresh history
      isForked.value = true
      const trainingPlanStore = useTrainingPlanStore()
      await trainingPlanStore.fetchHistory()

      return result.data.plan_id
    } else {
      console.error(result.error ? formatError(result.error) : 'Unknown error during upsertPlan')
      throw new Error(result.error ? formatError(result.error) : 'Upsert failed')
    }
  }

  function ensureRowIds(table: Row[]) {
    table.forEach((row) => {
      if (!row._id) {
        row._id = crypto.randomUUID()
      }
    })
  }

  return {
    // State
    sharedPlan,
    sharedHistory,
    isLoading,
    isFetchingHistory,
    error,
    isForked,
    // Pagination state
    historyHasMore,
    isLoadingMore,
    // Computed
    currentPlan,
    hasPlan,
    // Actions
    keepForever,
    fetchSharedPlanByHash,
    fetchSharedHistory,
    fetchMoreSharedHistory,
    loadPlanFromHistory,
    addPlanToHistory,
    clear,
    // Table Manipulation Actions
    updatePlanRow,
    addRow,
    removeRow,
    moveRow,
    upsertCurrentPlan,
  }
})
