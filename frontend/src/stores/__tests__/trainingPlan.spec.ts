// frontend/src/stores/__tests__/trainingPlan.spec.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import type { QueryRequest, RAGResponse, ApiResult, UpsertPlanResponse } from '@/types'
import { apiClient } from '@/api/client'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'
import type { Mock } from 'vitest'

// --- Mocks ---
vi.mock('@/api/client', async (importOriginal) => {
  const actual = (await importOriginal()) as typeof import('@/api/client')
  return {
    ...actual,
    apiClient: {
      query: vi.fn(),
      upsertPlan: vi.fn(),
    },
    formatError: vi.fn((error) => `${error.message}: ${error.details}`),
  }
})

vi.mock('@/plugins/supabase', () => ({
  supabase: {
    from: vi.fn().mockReturnThis(),
    select: vi.fn().mockReturnThis(),
    order: vi.fn().mockReturnThis(),
    in: vi.fn().mockReturnThis(),
    update: vi.fn().mockReturnThis(),
    eq: vi.fn().mockReturnThis(),
    limit: vi.fn().mockReturnThis(),
  },
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    user: { id: 'test-user-id' },
    getUser: vi.fn(),
  })),
}))

// --- Mock Casts ---
const mockedApiQuery = apiClient.query as Mock
const mockedApiUpsert = apiClient.upsertPlan as Mock
const mockedSupabase = supabase as unknown as {
  from: Mock
  select: Mock
  order: Mock
  in: Mock
  update: Mock
  eq: Mock
  limit: Mock
}
const mockedAuthStore = useAuthStore as unknown as Mock

// Helper to create a mock RAGResponse
const createMockPlan = (): RAGResponse => ({
  title: 'Test Plan',
  description: 'A plan for testing.',
  table: [
    {
      Amount: 1,
      Distance: 100,
      Sum: 100,
      Break: '10s',
      Content: 'Swim',
      Intensity: 'GA1',
      Multiplier: 'x',
    },
    {
      Amount: 2,
      Distance: 200,
      Sum: 400,
      Break: '20s',
      Content: 'Kick',
      Intensity: 'GA2',
      Multiplier: 'x',
    },
    {
      Amount: 0,
      Distance: 0,
      Sum: 500,
      Break: '',
      Content: 'Total',
      Intensity: '',
      Multiplier: '',
    },
  ],
})

describe('trainingPlan Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.resetAllMocks()
    mockedAuthStore.mockReturnValue({
      user: { id: 'test-user-id' },
      getUser: vi.fn(),
    })
    mockedSupabase.from.mockImplementation(() => ({
      select: vi.fn().mockReturnThis(),
      order: vi.fn().mockReturnThis(),
      limit: vi.fn().mockResolvedValue({ data: [], error: null }),
      in: vi.fn().mockResolvedValue({ data: [], error: null }),
      update: vi.fn().mockReturnThis(),
      eq: vi.fn().mockResolvedValue({ error: null }),
    }))
  })
  it('verify initial state', () => {
    const store = useTrainingPlanStore()

    expect(store.currentPlan).toBeNull()
    expect(store.isLoading).toBe(false)
    expect(store.error).toBeNull()
    expect(store.hasPlan).toBe(false)
  })
  it('generates a plan successfully and updates the store', async () => {
    const store = useTrainingPlanStore()

    // Define a mock response for apiClient.query
    const mockResponse: ApiResult<RAGResponse> = {
      success: true,
      data: {
        title: 'Generated Plan',
        description: 'A plan generated for testing.',
        table: [], // Simplified for this test
      },
    }
    mockedApiQuery.mockResolvedValue(mockResponse)

    // Mock the fetchHistory call that happens on success
    mockedSupabase.from.mockImplementation((tableName: string) => {
      if (tableName === 'history') {
        return {
          select: vi.fn().mockReturnThis(),
          order: vi.fn().mockReturnThis(),
          limit: vi.fn().mockResolvedValue({ data: [], error: null }),
        }
      }
      if (tableName === 'plans') {
        return {
          select: vi.fn().mockReturnThis(),
          in: vi.fn().mockResolvedValue({ data: [], error: null }),
        }
      }
      return {}
    })

    const requestPayload: QueryRequest = {
      content: 'test query',
      method: 'generate',
      filter: {},
      language: 'en',
    }

    // Call the action
    const result = await store.generatePlan(requestPayload)

    // Assertions
    expect(result).toBe(true) // The action should return true on success
    expect(store.isLoading).toBe(false) // Should no longer be loading
    expect(store.error).toBeNull() // Should have no error
    expect(store.currentPlan).toEqual(mockResponse.data) // The plan should be set
    expect(store.hasPlan).toBe(true) // hasPlan should be true

    // Verify that apiClient.query was called with the correct payload
    expect(mockedApiQuery).toHaveBeenCalledTimes(1)
    expect(mockedApiQuery).toHaveBeenCalledWith(requestPayload)
  })
  it('handles plan generation failure and updates the store with an error', async () => {
    const store = useTrainingPlanStore()

    // Define a mock error response for apiClient.query
    const mockErrorResponse: ApiResult<RAGResponse> = {
      success: false,
      error: {
        status: 500,
        details: 'API_ERROR',
        message: 'Failed to connect to the plan generation service.',
      },
    }

    // Tell our mocked apiClient.query to return the mockErrorResponse
    mockedApiQuery.mockResolvedValue(mockErrorResponse)

    const requestPayload: QueryRequest = {
      content: 'test query',
      method: 'generate',
      filter: {},
      language: 'en',
    }

    // Call the action
    const result = await store.generatePlan(requestPayload)

    // Assertions
    expect(result).toBe(false) // The action should return false on failure
    expect(store.isLoading).toBe(false) // Should no longer be loading
    expect(store.currentPlan).toBeNull() // Plan should remain null
    expect(store.hasPlan).toBe(false) // hasPlan should be false
    expect(store.error).toBe(
      `${mockErrorResponse.error?.message}: ${mockErrorResponse.error?.details}`,
    ) // Error message should be set

    // Verify that apiClient.query was called
    expect(mockedApiQuery).toHaveBeenCalledTimes(1)
    expect(mockedApiQuery).toHaveBeenCalledWith(requestPayload)
  })

  it('adds a row at the specified index and recalculates sum', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    const initialRowCount = store.currentPlan.table.length
    const initialSum = store.currentPlan.table[initialRowCount - 1].Sum

    // Add a new row at index 1
    store.addRow(1)

    const newRowCount = store.currentPlan.table.length
    expect(newRowCount).toBe(initialRowCount + 1)

    // Check that the new row is at the correct position and has default values
    const newRow = store.currentPlan.table[1]
    expect(newRow.Amount).toBe(0)
    expect(newRow.Content).toBe('')
    expect(newRow.Sum).toBe(0)

    // The total sum should be unchanged since the new row has a sum of 0
    const newSum = store.currentPlan.table[newRowCount - 1].Sum
    expect(newSum).toBe(initialSum)
  })

  it('removes a row at the specified index and recalculates sum', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    const initialRowCount = store.currentPlan.table.length
    const rowToRemove = store.currentPlan.table[1] // Sum is 400
    const initialSum = store.currentPlan.table[initialRowCount - 1].Sum // Sum is 500

    // Remove the row at index 1
    store.removeRow(1)

    const newRowCount = store.currentPlan.table.length
    expect(newRowCount).toBe(initialRowCount - 1)

    // Check that the correct row was removed
    expect(store.currentPlan.table.find((r) => r.Content === 'Kick')).toBeUndefined()

    // The new total sum should be the initial sum minus the sum of the removed row
    const newSum = store.currentPlan.table[newRowCount - 1].Sum
    expect(newSum).toBe(initialSum - rowToRemove.Sum) // 500 - 400 = 100
  })

  it('does not add a row if the table has 26 or more rows', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    // Fill the table with 26 rows
    store.currentPlan.table = Array.from({ length: 26 }, (_, i) => ({
      Amount: 1,
      Distance: 100,
      Sum: 100,
      Break: '10s',
      Content: `Swim ${i}`,
      Intensity: 'GA1',
      Multiplier: 'x',
    }))

    const initialRowCount = store.currentPlan.table.length
    store.addRow(1)
    const newRowCount = store.currentPlan.table.length
    expect(newRowCount).toBe(initialRowCount)
  })

  it('does not remove a row if only one exercise row is left', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = {
      title: 'Test Plan',
      description: 'A plan for testing.',
      table: [
        {
          Amount: 1,
          Distance: 100,
          Sum: 100,
          Break: '10s',
          Content: 'Swim',
          Intensity: 'GA1',
          Multiplier: 'x',
        },
        {
          Amount: 0,
          Distance: 0,
          Sum: 100,
          Break: '',
          Content: 'Total',
          Intensity: '',
          Multiplier: '',
        },
      ],
    }

    const initialRowCount = store.currentPlan.table.length
    store.removeRow(0)
    const newRowCount = store.currentPlan.table.length
    expect(newRowCount).toBe(initialRowCount)
  })

  it('moves a row up', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    const rowToMove = store.currentPlan.table[1]
    store.moveRow(1, 'up')

    expect(store.currentPlan.table[0]).toBe(rowToMove)
  })

  it('moves a row down', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    const rowToMove = store.currentPlan.table[0]
    store.moveRow(0, 'down')

    expect(store.currentPlan.table[1]).toBe(rowToMove)
  })

  it('does not move the first row up', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    const initialOrder = [...store.currentPlan.table]
    store.moveRow(0, 'up')

    expect(store.currentPlan.table).toEqual(initialOrder)
  })

  it('does not move the last exercise row down', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    const initialOrder = [...store.currentPlan.table]
    const lastExerciseRowIndex = store.currentPlan.table.length - 2
    store.moveRow(lastExerciseRowIndex, 'down')

    expect(store.currentPlan.table).toEqual(initialOrder)
  })

  describe('History Management', () => {
    it('fetches history successfully', async () => {
      const store = useTrainingPlanStore()
      const mockHistory = [{ plan_id: 'plan-1' }, { plan_id: 'plan-2' }]
      const mockPlans = [
        {
          plan_id: 'plan-1',
          title: 'Plan 1',
          description: 'Desc 1',
          table: [],
        },
        {
          plan_id: 'plan-2',
          title: 'Plan 2',
          description: 'Desc 2',
          table: [],
        },
      ]

      // Mock supabase calls for history IDs
      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'history') {
          return {
            select: vi.fn().mockReturnThis(),
            order: vi.fn().mockReturnThis(),
            limit: vi.fn().mockResolvedValue({ data: mockHistory, error: null }),
          }
        }
        if (tableName === 'plans') {
          return {
            select: vi.fn().mockReturnThis(),
            in: vi.fn().mockResolvedValue({ data: mockPlans, error: null }),
          }
        }
        return {
          select: vi.fn().mockReturnThis(),
          order: vi.fn().mockReturnThis(),
          in: vi.fn().mockReturnThis(),
        }
      })

      await store.fetchHistory()

      expect(store.isLoading).toBe(false)
      expect(store.generationHistory).toHaveLength(2)
      expect(store.generationHistory[0].title).toBe('Plan 1')
      expect(mockedSupabase.from).toHaveBeenCalledWith('history')
      expect(mockedSupabase.from).toHaveBeenCalledWith('plans')
    })

    it('does not fetch history if user is not available', async () => {
      mockedAuthStore.mockReturnValue({ user: null, getUser: vi.fn() })
      const store = useTrainingPlanStore()

      await store.fetchHistory()

      expect(mockedSupabase.from).not.toHaveBeenCalled()
      expect(store.generationHistory).toEqual([])
    })

    it('handles error when fetching history plan IDs', async () => {
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      mockedSupabase.from.mockReturnValue({
        select: vi.fn().mockReturnThis(),
        order: vi.fn().mockReturnThis(),
        limit: vi.fn().mockResolvedValue({ data: null, error: new Error('DB Error') }),
      })
      const store = useTrainingPlanStore()

      await store.fetchHistory()

      expect(store.isLoading).toBe(false)
      expect(store.generationHistory).toEqual([])
      expect(consoleErrorSpy).toHaveBeenCalledWith(new Error('DB Error'))
      consoleErrorSpy.mockRestore()
    })

    it('handles error when fetching plan details', async () => {
      const store = useTrainingPlanStore()
      const mockHistory = [{ plan_id: 'plan-1' }]
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'history') {
          return {
            select: vi.fn().mockReturnThis(),
            order: vi.fn().mockReturnThis(),
            limit: vi.fn().mockResolvedValue({ data: mockHistory, error: null }),
          }
        }
        if (tableName === 'plans') {
          return {
            select: vi.fn().mockReturnThis(),
            in: vi.fn().mockResolvedValue({ data: null, error: new Error('Plan fetch error') }),
          }
        }
        return {
          select: vi.fn().mockReturnThis(),
        }
      })

      await store.fetchHistory()

      expect(store.isLoading).toBe(false)
      expect(store.generationHistory).toEqual([])
      expect(consoleErrorSpy).toHaveBeenCalledWith(new Error('Plan fetch error'))
      consoleErrorSpy.mockRestore()
    })

    it('upserts a plan successfully', async () => {
      const store = useTrainingPlanStore()
      const planToUpsert = {
        title: 'New Plan',
        description: 'A new plan',
        table: [],
      }
      const mockResponse: ApiResult<UpsertPlanResponse> = {
        success: true,
        data: { plan_id: 'new-plan-1' },
      }
      mockedApiUpsert.mockResolvedValue(mockResponse)

      // Mock fetchHistory to be empty
      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'history') {
          return {
            select: vi.fn().mockReturnThis(),
            order: vi.fn().mockReturnThis(),
            limit: vi.fn().mockResolvedValue({ data: [], error: null }),
          }
        }
        if (tableName === 'plans') {
          return {
            select: vi.fn().mockReturnThis(),
            in: vi.fn().mockResolvedValue({ data: [], error: null }),
          }
        }
        return {}
      })

      const result = await store.upsertPlan(planToUpsert)

      expect(result).toEqual(mockResponse.data)
      expect(mockedApiUpsert).toHaveBeenCalledWith(planToUpsert)
      expect(store.isLoading).toBe(false)
      // fetchHistory is called on success
      expect(mockedSupabase.from).toHaveBeenCalledWith('history')
    })

    it('does not upsert plan if user is not available', async () => {
      mockedAuthStore.mockReturnValue({ user: null, getUser: vi.fn() })
      const store = useTrainingPlanStore()
      const planToUpsert = {
        title: 'New Plan',
        description: 'A new plan',
        table: [],
      }
      // Mock the api call to prevent failure if the user check fails
      mockedApiUpsert.mockResolvedValue({ success: false, error: 'Should not be called' })

      const result = await store.upsertPlan(planToUpsert)

      expect(result).toBeNull()
      expect(mockedApiUpsert).not.toHaveBeenCalled()
    })

    it('handles upsert plan failure', async () => {
      const store = useTrainingPlanStore()
      const planToUpsert = {
        title: 'New Plan',
        description: 'A new plan',
        table: [],
      }
      const mockErrorResponse: ApiResult<UpsertPlanResponse> = {
        success: false,
        error: { status: 500, message: 'Server Error', details: 'UPSERT_FAILED' },
      }
      mockedApiUpsert.mockResolvedValue(mockErrorResponse)
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

      const result = await store.upsertPlan(planToUpsert)

      expect(result).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(consoleErrorSpy).toHaveBeenCalled()
      consoleErrorSpy.mockRestore()
    })

    it('loads a plan from history and ensures it is a deep copy', () => {
      const store = useTrainingPlanStore()
      const planFromHistory = createMockPlan()
      store.generationHistory = [planFromHistory]

      store.loadPlanFromHistory(planFromHistory)

      expect(store.currentPlan).toEqual(planFromHistory)
      // Verify it's a deep copy by modifying the loaded plan
      store.updatePlanRow(0, 'Amount', 99)
      expect(store.currentPlan?.table[0].Amount).toBe(99)
      expect(planFromHistory.table[0].Amount).toBe(1) // Original should be unchanged
    })

    it('toggles keep_forever status for a plan', async () => {
      const store = useTrainingPlanStore()
      const planId = 'plan-1'
      store.historyMetadata = [
        { plan_id: planId, keep_forever: false, created_at: '', updated_at: '' },
      ]

      const historyMock = {
        update: vi.fn().mockReturnThis(),
        eq: vi.fn().mockResolvedValue({ error: null }),
      }

      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'history') return historyMock
        return {}
      })

      await store.toggleKeepForever(planId)

      expect(mockedSupabase.from).toHaveBeenCalledWith('history')
      expect(historyMock.update).toHaveBeenCalledWith({ keep_forever: true })
      expect(historyMock.eq).toHaveBeenCalledWith('plan_id', planId)
      expect(store.historyMetadata[0].keep_forever).toBe(true)

      // Toggle back
      await store.toggleKeepForever(planId)
      expect(historyMock.update).toHaveBeenCalledWith({ keep_forever: false })
      expect(store.historyMetadata[0].keep_forever).toBe(false)
    })

    it('does not toggle keep_forever if user is not available', async () => {
      mockedAuthStore.mockReturnValue({ user: null, getUser: vi.fn() })
      const store = useTrainingPlanStore()
      const planId = 'plan-1'

      await store.toggleKeepForever(planId)

      expect(mockedSupabase.from).not.toHaveBeenCalled()
    })

    it('handles error when toggling keep_forever', async () => {
      const store = useTrainingPlanStore()
      const planId = 'plan-1'
      store.historyMetadata = [
        { plan_id: planId, keep_forever: false, created_at: '', updated_at: '' },
      ]
      const dbError = new Error('Update failed')
      const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'history') {
          return {
            update: vi.fn().mockReturnThis(),
            eq: vi.fn().mockResolvedValue({ error: dbError }),
          }
        }
        return {}
      })

      await store.toggleKeepForever(planId)

      expect(consoleErrorSpy).toHaveBeenCalledWith(dbError)
      consoleErrorSpy.mockRestore()
    })
  })
})
