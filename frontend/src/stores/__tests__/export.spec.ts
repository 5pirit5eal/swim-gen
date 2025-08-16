// frontend/src/stores/__tests__/export.spec.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useExportStore } from '../export'
import { apiClient } from '@/api/client'
import type { Mock } from 'vitest'
import type { PlanToPDFRequest, ApiResult, PlanToPDFResponse } from '@/types'

// Mock the apiClient module
vi.mock('@/api/client', () => ({
  apiClient: {
    exportPDF: vi.fn(),
  },
}))

// Cast apiClient.exportPDF to a Mock type for TypeScript
const mockedApiExportPDF = apiClient.exportPDF as Mock<typeof apiClient.exportPDF>

describe('export Store', () => {
  beforeEach(() => {
    // Reset the store state manually for setup stores
    const store = useExportStore()
    store.isExporting = false
    store.exportError = null
    vi.clearAllMocks() // Clear mock history
  })

  it('has the correct initial state', () => {
    const store = useExportStore()

    expect(store.isExporting).toBe(false)
    expect(store.exportError).toBeNull()
  })
  it('exports PDF successfully and returns the URI', async () => {
    const store = useExportStore()

    // Define a mock response for apiClient.exportPDF
    const mockResponse: ApiResult<PlanToPDFResponse> = {
      success: true,
      data: {
        uri: 'http://mock.pdf/plan123.pdf',
      },
    }

    // Tell our mocked apiClient.exportPDF to return the mockResponse
    mockedApiExportPDF.mockResolvedValue(mockResponse)

    const requestPayload: PlanToPDFRequest = {
      title: 'Test Plan',
      description: 'A plan for export.',
      table: [],
    }

    // Call the action
    const result = await store.exportToPDF(requestPayload)

    // Assertions
    expect(result).toBe('http://mock.pdf/plan123.pdf') // Should return the URI
    expect(store.isExporting).toBe(false) // Should no longer be exporting
    expect(store.exportError).toBeNull() // Should have no error

    // Verify that apiClient.exportPDF was called with the correct payload
    expect(mockedApiExportPDF).toHaveBeenCalledTimes(1)
    expect(mockedApiExportPDF).toHaveBeenCalledWith(requestPayload)
  })
  it('handles PDF export failure and updates the store with an error', async () => {
    const store = useExportStore()

    // Define a mock error response for apiClient.exportPDF
    const mockErrorResponse: ApiResult<PlanToPDFResponse> = {
      success: false,
      error: {
        status: 500,
        details: 'EXPORT_ERROR',
        message: 'Failed to export PDF due to server issue.',
      },
    }

    // Tell our mocked apiClient.exportPDF to return the mockErrorResponse
    mockedApiExportPDF.mockResolvedValue(mockErrorResponse)

    const requestPayload: PlanToPDFRequest = {
      title: 'Test Plan',
      description: 'A plan for export.',
      table: [],
    }

    // Call the action
    const result = await store.exportToPDF(requestPayload)

    // Assertions
    expect(result).toBeNull() // Should return null on failure
    expect(store.isExporting).toBe(false) // Should no longer be exporting
    expect(store.exportError).toBe('Failed to export PDF due to server issue.') // Error message should be set

    // Verify that apiClient.exportPDF was called
    expect(mockedApiExportPDF).toHaveBeenCalledTimes(1)
    expect(mockedApiExportPDF).toHaveBeenCalledWith(requestPayload)
  })

})
