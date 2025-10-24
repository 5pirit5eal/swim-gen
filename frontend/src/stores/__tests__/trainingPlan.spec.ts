// frontend/src/stores/__tests__/trainingPlan.spec.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useTrainingPlanStore } from '@/stores/trainingPlan'
import type { QueryRequest, RAGResponse, ApiResult } from '@/types'
import { apiClient } from '@/api/client'
import type { Mock } from 'vitest' // Import Mock type

// --- This is the mock ---
// We are telling Vitest: "Whenever someone imports from '@/api/client',
// don't give them the real module. Instead, give them this fake object."
vi.mock('@/api/client', async (importOriginal) => {
  const actual = (await importOriginal()) as typeof import('@/api/client')
  return {
    ...actual,
    apiClient: {
      query: vi.fn(),
    },
    formatError: vi.fn((error) => `${error.message}: ${error.details}`), // Mock formatError
  }
})

// Cast apiClient.query to a Mock type for TypeScript to recognize mock methods
const mockedApiQuery = apiClient.query as Mock<typeof apiClient.query>

// Helper to create a mock RAGResponse
const createMockPlan = (): RAGResponse => ({
  title: 'Test Plan',
  description: 'A plan for testing.',
  table: [
    { Amount: 1, Distance: 100, Sum: 100, Break: '10s', Content: 'Swim', Intensity: 'GA1', Multiplier: 'x' },
    { Amount: 2, Distance: 200, Sum: 400, Break: '20s', Content: 'Kick', Intensity: 'GA2', Multiplier: 'x' },
    { Amount: 0, Distance: 0, Sum: 500, Break: '', Content: 'Total', Intensity: '', Multiplier: '' }
  ]
})

describe('trainingPlan Store', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    const store = useTrainingPlanStore()
    store.currentPlan = null
    store.isLoading = false
    store.error = null
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

    // Tell our mocked apiClient.query to return the mockResponse
    mockedApiQuery.mockResolvedValue(mockResponse)

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
    expect(store.currentPlan.table.find(r => r.Content === 'Kick')).toBeUndefined()

    // The new total sum should be the initial sum minus the sum of the removed row
    const newSum = store.currentPlan.table[newRowCount - 1].Sum
    expect(newSum).toBe(initialSum - rowToRemove.Sum) // 500 - 400 = 100
  })

  it('does not add a row if the table has 26 or more rows', () => {
    const store = useTrainingPlanStore()
    store.currentPlan = createMockPlan()

    // Fill the table with 26 rows
    store.currentPlan.table = Array.from({ length: 26 }, (_, i) => ({
      Amount: 1, Distance: 100, Sum: 100, Break: '10s', Content: `Swim ${i}`, Intensity: 'GA1', Multiplier: 'x'
    }));

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
        { Amount: 1, Distance: 100, Sum: 100, Break: '10s', Content: 'Swim', Intensity: 'GA1', Multiplier: 'x' },
        { Amount: 0, Distance: 0, Sum: 100, Break: '', Content: 'Total', Intensity: '', Multiplier: '' }
      ]
    }

    const initialRowCount = store.currentPlan.table.length
    store.removeRow(0)
    const newRowCount = store.currentPlan.table.length
    expect(newRowCount).toBe(initialRowCount)
  })
})
