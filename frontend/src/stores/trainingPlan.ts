import { apiClient, formatError } from '@/api/client'
import i18n from '@/plugins/i18n'
import type {
  QueryRequest,
  RAGResponse,
  Row,
  HistoryMetadata,
  Message,
  FeedbackRequest,
} from '@/types'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'

export const useTrainingPlanStore = defineStore('trainingPlan', () => {
  // --- STATE ---
  const currentPlan = ref<RAGResponse | null>(null)
  const isLoading = ref(false)
  const isFetchingHistory = ref(false)
  const isFetchingConversation = ref(false)
  const error = ref<string | null>(null)
  const initialQuery = ref<string>('')
  const generationHistory = ref<RAGResponse[]>([])
  const historyMetadata = ref<HistoryMetadata[]>([])
  const conversation = ref<Message[]>([])
  const userStore = useAuthStore()

  // Pagination state
  const PAGE_SIZE = 10
  const historyPage = ref(0)
  const historyHasMore = ref(true)
  const isLoadingMore = ref(false)

  // Search state
  const SEARCH_LIMIT = 20
  const searchQuery = ref('')
  const searchResults = ref<(RAGResponse & HistoryMetadata)[]>([])
  const isSearching = ref(false)
  const searchHitLimit = computed(() => searchResults.value.length >= SEARCH_LIMIT)

  // --- COMPUTED ---
  const hasPlan = computed(() => currentPlan.value !== null)
  const planHistory = computed(() => {
    // If searching, show only search results
    if (searchQuery.value.trim()) {
      return searchResults.value
    }
    // Otherwise show normal history
    const combined = historyMetadata.value.map((metadata) => {
      const plan = generationHistory.value.find((p) => p.plan_id === metadata.plan_id)
      return {
        ...plan,
        ...metadata,
      }
    })
    return combined.filter(
      (plan): plan is RAGResponse & HistoryMetadata =>
        !!plan.title && !!plan.description && !!plan.table,
    )
  })

  watch(
    () => userStore.user?.id ?? null,
    async (newUserId) => {
      if (newUserId) {
        await fetchHistory()
      } else {
        generationHistory.value = []
        historyMetadata.value = []
      }
    },
    { immediate: true },
  )

  // --- ACTIONS ---

  // Fetches the user's plan history with pagination
  async function fetchHistory(reset = true) {
    console.debug('[TrainingPlanStore] fetchHistory', { reset, page: historyPage.value })
    if (!userStore.user) return
    if (reset) {
      historyPage.value = 0
      historyHasMore.value = true
      generationHistory.value = []
      historyMetadata.value = []
    }

    isFetchingHistory.value = true
    const offset = historyPage.value * PAGE_SIZE
    const { data, error } = await supabase
      .from('history')
      .select('plan_id, keep_forever, created_at, updated_at, exported_at')
      .order('created_at', { ascending: false })
      .range(offset, offset + PAGE_SIZE - 1)

    if (error) {
      console.error(error)
    } else if (data) {
      console.debug('[TrainingPlanStore] fetchHistory success', { count: data.length })
      // Check if there are more results
      historyHasMore.value = data.length === PAGE_SIZE

      const planIds = data.map((entry) => entry.plan_id)

      // Fetch feedback for these plans
      const { data: feedbackData, error: feedbackError } = await supabase
        .from('feedback')
        .select('plan_id, rating, was_swam, difficulty_rating')
        .in('plan_id', planIds)
        .eq('user_id', userStore.user.id)

      if (feedbackError) {
        console.error('Error fetching feedback:', feedbackError)
      }

      const newMetadata = data.map((entry) => {
        const feedback = feedbackData?.find((f) => f.plan_id === entry.plan_id)
        return {
          plan_id: entry.plan_id,
          keep_forever: entry.keep_forever,
          created_at: entry.created_at,
          updated_at: entry.updated_at,
          exported_at: entry.exported_at,
          feedback_rating: feedback?.rating,
          was_swam: feedback?.was_swam,
          difficulty_rating: feedback?.difficulty_rating,
        }
      })

      if (reset) {
        historyMetadata.value = newMetadata
      } else {
        historyMetadata.value = [...historyMetadata.value, ...newMetadata]
      }

      const { data: plansData, error: plansError } = await supabase
        .from('plans')
        .select('plan_id, title, description, plan_table')
        .in('plan_id', planIds)
      if (plansError) {
        console.error(plansError)
      } else if (plansData) {
        const newPlans = plansData.map((plan) => ({
          plan_id: plan.plan_id,
          title: plan.title,
          description: plan.description,
          table: plan.plan_table,
        }))
        if (reset) {
          generationHistory.value = newPlans
        } else {
          generationHistory.value = [...generationHistory.value, ...newPlans]
        }
      }
    }
    isFetchingHistory.value = false
  }

  // Fetches more history entries (pagination)
  async function fetchMoreHistory() {
    if (!historyHasMore.value || isLoadingMore.value) return
    isLoadingMore.value = true
    historyPage.value += 1
    await fetchHistory(false)
    isLoadingMore.value = false
  }

  // Searches plans by title or description
  async function searchPlans(query: string) {
    console.debug('[TrainingPlanStore] searchPlans', { query })
    searchQuery.value = query
    if (!query.trim()) {
      searchResults.value = []
      return
    }
    if (!userStore.user) return

    isSearching.value = true
    try {
      // Search in history table joined with plans table
      const { data: historyData, error: historyError } = await supabase
        .from('history')
        .select('plan_id, keep_forever, created_at, updated_at, exported_at')
        .order('created_at', { ascending: false })

      if (historyError) {
        console.error('Search history error:', historyError)
        return
      }

      if (!historyData || historyData.length === 0) {
        searchResults.value = []
        return
      }

      const planIds = historyData.map((entry) => entry.plan_id)

      // Search plans matching the query
      const searchPattern = `%${query}%`
      const { data: plansData, error: plansError } = await supabase
        .from('plans')
        .select('plan_id, title, description, plan_table')
        .in('plan_id', planIds)
        .or(`title.ilike.${searchPattern},description.ilike.${searchPattern}`)
        .limit(SEARCH_LIMIT)

      if (plansError) {
        console.error('Search plans error:', plansError)
        return
      }

      if (!plansData) {
        searchResults.value = []
        return
      }

      // Combine with metadata
      searchResults.value = plansData.map((plan) => {
        const metadata = historyData.find((h) => h.plan_id === plan.plan_id)
        return {
          plan_id: plan.plan_id,
          title: plan.title,
          description: plan.description,
          table: plan.plan_table,
          keep_forever: metadata?.keep_forever ?? false,
          created_at: metadata?.created_at ?? '',
          updated_at: metadata?.updated_at ?? '',
          exported_at: metadata?.exported_at,
        }
      })
      console.debug('[TrainingPlanStore] searchPlans success', {
        count: searchResults.value.length,
      })
    } finally {
      isSearching.value = false
    }
  }

  // Clears search results
  function clearSearch() {
    searchQuery.value = ''
    searchResults.value = []
  }

  // Update keep_forever for plan in history
  async function toggleKeepForever(planId: string | undefined) {
    if (!userStore.user || !planId) return
    const metadataEntry = historyMetadata.value.find((entry) => entry.plan_id === planId)
    if (!metadataEntry) return
    const newKeepForever = !metadataEntry.keep_forever

    const { error } = await supabase
      .from('history')
      .update({ keep_forever: newKeepForever })
      .eq('plan_id', planId)
    if (error) {
      console.error(error)
    } else {
      metadataEntry.keep_forever = newKeepForever
    }
  }

  // Lets the user keep a plan forever, does nothing if the plan is already kept forever
  async function keepForever(planId: string) {
    if (!userStore.user) return
    const metadataEntry = historyMetadata.value.find((entry) => entry.plan_id === planId)
    if (!metadataEntry) return
    if (metadataEntry.keep_forever) return
    await toggleKeepForever(planId)
  }

  // Generates a new training plan
  async function generatePlan(request: QueryRequest): Promise<boolean> {
    console.debug('[TrainingPlanStore] generatePlan', { request })
    isLoading.value = true
    error.value = null
    if (!userStore.user) {
      initialQuery.value = request.content
    }

    const result = await apiClient.query(request)

    if (result.success && result.data) {
      currentPlan.value = result.data
      ensureRowIds(currentPlan.value.table)
      recalculateTotalSum()
      await fetchHistory() // Refresh history after generating a new plan
      isLoading.value = false
      if (result.data.plan_id) await fetchConversation(result.data.plan_id)
      return true
    } else {
      error.value = result.error
        ? formatError(result.error)
        : i18n.global.t('errors.training_plan_failed')
      isLoading.value = false
      return false
    }
  }

  // Upserts the current plan
  async function upsertCurrentPlan(): Promise<string> {
    console.debug('[TrainingPlanStore] upsertCurrentPlan', { planId: currentPlan.value?.plan_id })
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
      await fetchHistory() // Refresh history after upserting
      return result.data.plan_id
    } else {
      console.error(result.error ? formatError(result.error) : 'Unknown error during upsertPlan')
      error.value = result.error
        ? formatError(result.error)
        : i18n.global.t('errors.training_plan_failed')
      throw new Error(error.value)
    }
  }

  // Links an anonymous plan to the user's history
  async function linkAnonymousPlan() {
    console.debug('Linking anonymous plan to user account')
    if (!userStore.user || !currentPlan.value || !initialQuery.value) return

    // 1. Add plan to history (gets a new plan_id)
    // We need to strip _id from table rows
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const tableWithoutIds = currentPlan.value.table.map(({ _id, ...rest }) => rest)

    const addPlanResult = await apiClient.addPlanToHistory({
      title: currentPlan.value.title,
      description: currentPlan.value.description,
      table: tableWithoutIds,
    })

    if (!addPlanResult.success || !addPlanResult.data) {
      console.error('Failed to add plan to history:', addPlanResult.error)
      return
    }

    const newPlanId = addPlanResult.data.plan_id
    currentPlan.value.plan_id = newPlanId

    // 2. Add User Message
    const userMsgResult = await apiClient.addMessage(newPlanId, 'user', initialQuery.value)

    if (!userMsgResult.success || !userMsgResult.data) {
      console.error('Failed to add user message:', userMsgResult.error)
      return
    }

    // 3. Add Assistant Message (with plan snapshot)
    await apiClient.addMessage(
      newPlanId,
      'ai',
      currentPlan.value.description,
      userMsgResult.data.message_id,
      {
        ...currentPlan.value,
        plan_id: newPlanId,
        table: tableWithoutIds,
      },
    )

    // Refresh history and conversation
    await fetchHistory()
    await fetchConversation(newPlanId)

    // Clear initialQuery to prevent re-linking
    initialQuery.value = ''
  }

  // Loads a plan from history into the editor
  async function loadPlanFromHistory(plan: RAGResponse) {
    if (!userStore.user) return
    if (!plan.plan_id) return
    await fetchConversation(plan.plan_id)
    currentPlan.value = JSON.parse(JSON.stringify(plan)) // Deep copy to prevent accidental edits
    if (currentPlan.value) {
      recalculateTotalSum()
      ensureRowIds(currentPlan.value.table)
    }
  }

  // --- Plan Table Manipulations ---

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

  function recalculateTotalSum() {
    if (currentPlan.value && currentPlan.value.table.length > 0) {
      const lastRowIndex = currentPlan.value.table.length - 1
      const lastRow = currentPlan.value.table[lastRowIndex]!
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
    table.splice(newIndex, 0, movedRow!)
  }

  function clearError() {
    error.value = null
  }

  function clear() {
    currentPlan.value = null
    error.value = null
    isLoading.value = false
    conversation.value = []
  }

  // Fetches the conversation history for a plan
  async function fetchConversation(planId: string) {
    if (!userStore.user) return
    isFetchingConversation.value = true
    conversation.value = []

    const { data, error: fetchError } = await apiClient.getConversation(planId)
    if (fetchError) {
      console.error(fetchError)
      return
    }

    if (data !== null && Array.isArray(data)) {
      conversation.value = data
    } else {
      conversation.value = []
    }
    isFetchingConversation.value = false
  }

  // Saves a snapshot to user history
  async function saveSnapshot(plan: RAGResponse) {
    if (!userStore.user) return

    const result = await apiClient.addPlanToHistory(plan)
    if (result.success) {
      await fetchHistory() // Refresh history after saving snapshot
    } else {
      console.error('Failed to save snapshot:', result.error)
    }
  }

  // Sends a message to the AI
  async function sendMessage(message: string) {
    if (!currentPlan.value?.plan_id || !userStore.user) return

    // Optimistic update: add user message
    const userMsg: Message = {
      id: crypto.randomUUID(),
      role: 'user',
      content: message,
      created_at: new Date().toISOString(),
      plan_id: currentPlan.value.plan_id,
      user_id: userStore.user.id,
      previous_message_id: null,
      next_message_id: null,
    }
    conversation.value.push(userMsg)

    isLoading.value = true
    error.value = null

    const result = await apiClient.chat({
      plan_id: currentPlan.value.plan_id,
      message: message,
    })

    if (result.success && result.data) {
      // Add AI response
      const aiMsg: Message = {
        id: crypto.randomUUID(),
        role: 'ai',
        content: result.data.response,
        created_at: new Date().toISOString(),
        plan_id: currentPlan.value.plan_id,
        user_id: userStore.user.id,
        previous_message_id: null,
        next_message_id: null,
        plan_snapshot: result.data.table
          ? {
              plan_id: result.data.plan_id,
              title: result.data.title || '',
              description: result.data.description || '',
              table: result.data.table,
            }
          : undefined,
      }
      conversation.value.push(aiMsg)

      // Update current plan if changed
      if (result.data.table) {
        currentPlan.value = {
          ...currentPlan.value,
          title: result.data.title || currentPlan.value.title,
          description: result.data.description || currentPlan.value.description,
          table: result.data.table,
          plan_id: result.data.plan_id,
        }
        ensureRowIds(currentPlan.value.table)
        recalculateTotalSum()
        await fetchHistory() // Refresh history after updating plan
      }
    } else {
      error.value = result.error ? formatError(result.error) : i18n.global.t('errors.unknown_error')
    }
    isLoading.value = false
  }

  // Submits feedback for a plan
  async function submitFeedback(payload: FeedbackRequest): Promise<boolean> {
    if (!userStore.user) return false

    const result = await apiClient.submitFeedback({
      plan_id: payload.plan_id,
      rating: payload.rating,
      was_swam: payload.was_swam,
      difficulty_rating: payload.difficulty_rating,
      comment: payload.comment,
    })

    if (result.success) {
      // Optimistically update history
      const metadata = historyMetadata.value.find((p) => p.plan_id === payload.plan_id)
      if (metadata) {
        metadata.feedback_rating = payload.rating
        metadata.was_swam = payload.was_swam
        metadata.difficulty_rating = payload.difficulty_rating
      }
      return true
    } else {
      console.error('Failed to submit feedback:', result.error)
      return false
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
    currentPlan,
    isLoading,
    isFetchingHistory,
    isFetchingConversation,
    error,
    generationHistory,
    historyMetadata,
    conversation,
    initialQuery,
    // Pagination state
    historyHasMore,
    isLoadingMore,
    // Search state
    searchQuery,
    searchResults,
    isSearching,
    searchHitLimit,
    // Computed
    hasPlan,
    planHistory,
    // Actions
    generatePlan,
    updatePlanRow,
    addRow,
    removeRow,
    moveRow,
    clearError,
    clear,
    fetchHistory,
    fetchMoreHistory,
    searchPlans,
    clearSearch,
    upsertCurrentPlan,
    loadPlanFromHistory,
    keepForever,
    toggleKeepForever,
    fetchConversation,
    saveSnapshot,
    sendMessage,
    submitFeedback,
    linkAnonymousPlan,
  }
})
