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
// Cast apiClient.query to a Mock type for TypeScript to recognize mock methods
const mockedApiQuery = apiClient.query as Mock<typeof apiClient.query>

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
})
