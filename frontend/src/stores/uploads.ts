import { apiClient, formatError } from '@/api/client'
import i18n from '@/plugins/i18n'
import type {
    DonatedPlan,
    DonatePlanRequest,
    RAGResponse,
    Row,
} from '@/types'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const useDonationStore = defineStore('donation', () => {
    // --- STATE ---
    const currentPlan = ref<RAGResponse | null>(null)
    const isLoading = ref(false)
    const error = ref<string | null>(null)
    const uploadedPlans = ref<DonatedPlan[]>([])
    const isFetchingUploads = ref(false)
    const userStore = useAuthStore()

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
        { immediate: true }
    )

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
                }
            ]
        }
    }

    // Fetch all uploaded plans for the user
    async function fetchUploadedPlans() {
        if (!userStore.user) return
        isFetchingUploads.value = true
        const result = await apiClient.getUploadedPlans()
        if (result.success && Array.isArray(result.data)) {
            uploadedPlans.value = result.data
        } else {
            console.error(result.error ? formatError(result.error) : 'Failed to fetch uploaded plans')
        }
        isFetchingUploads.value = false
    }

    // Fetch a specific uploaded plan
    async function fetchUploadedPlan(planId: string): Promise<boolean> {
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
            ensureRowIds(currentPlan.value.table)
            recalculateTotalSum()
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


    // Submit the current plan as a donation
    async function uploadCurrentPlan(): Promise<boolean> {
        if (!currentPlan.value) return false
        isLoading.value = true
        error.value = null

        // Strip _id from table rows before sending
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const tableWithoutIds = currentPlan.value.table.slice(0, -1).map(({ _id, ...rest }) => rest)

        const request: DonatePlanRequest = {
            title: currentPlan.value.title,
            description: currentPlan.value.description,
            table: tableWithoutIds,
            language: i18n.global.locale.value,
        }

        const result = await apiClient.donatePlan(request)
        if (result.success) {
            await fetchUploadedPlans()
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
        // Not applicable for donation store, but required by interface
        console.log('keepForever not implemented for donation store', planId)
    }

    async function upsertCurrentPlan() {
        // Not applicable for donation store in the same way,
        // we don't auto-save to backend on every edit, only on explicit upload
        // But we can use this to trigger local validation or updates if needed
        recalculateTotalSum()
    }

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
        // Computed
        hasPlan,
        // Actions
        initNewPlan,
        fetchUploadedPlans,
        fetchUploadedPlan,
        loadPlanFromHistory,
        uploadCurrentPlan,
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
