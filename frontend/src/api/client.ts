/**
 * API Client for Swim RAG Backend
 * Handles HTTP requests with proper TypeScript typing
 */

import type {
  QueryRequest,
  RAGResponse,
  PlanToPDFRequest,
  PlanToPDFResponse,
  HealthCheckResponse,
  ApiResult,
} from '@/types'

class ApiClient {
  private baseUrl: string

  constructor(baseUrl = '/api') {
    this.baseUrl = baseUrl
  }

  /**
   * Check API health status
   */
  async checkHealth(): Promise<ApiResult<HealthCheckResponse>> {
    try {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), 5000)

      const response = await fetch(`${this.baseUrl}/health`, {
        signal: controller.signal,
      })

      clearTimeout(timeoutId)

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: 'Health check failed',
            status: response.status,
            details: response.statusText,
          },
        }
      }

      const data = await response.text()
      return {
        success: true,
        data,
      }
    } catch (error) {
      return {
        success: false,
        error: {
          message: error instanceof Error ? error.message : 'Network error',
          status: 0,
          details: 'Failed to connect to server',
        },
      }
    }
  }

  /**
   * Query for training plans (may take up to 60 seconds)
   */
  async query(request: QueryRequest): Promise<ApiResult<RAGResponse>> {
    try {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), 60000)

      const response = await fetch(`${this.baseUrl}/query`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
        signal: controller.signal,
      })

      clearTimeout(timeoutId)

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: 'Query failed',
            status: response.status,
            details: response.statusText,
          },
        }
      }

      const data = await response.json()
      return {
        success: true,
        data,
      }
    } catch (error) {
      return {
        success: false,
        error: {
          message: error instanceof Error ? error.message : 'Network error',
          status: 0,
          details:
            error instanceof Error && error.name === 'AbortError'
              ? 'Request timed out after 60 seconds'
              : 'Failed to connect to server',
        },
      }
    }
  }

  /**
   * Export training plan as PDF
   */
  async exportPDF(request: PlanToPDFRequest): Promise<ApiResult<PlanToPDFResponse>> {
    try {
      const response = await fetch(`${this.baseUrl}/export-pdf`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: 'Converting plan to PDF failed',
            status: response.status,
            details: response.statusText,
          },
        }
      }

      const data = await response.json()
      return {
        success: true,
        data,
      }
    } catch (error) {
      return {
        success: false,
        error: {
          message: error instanceof Error ? error.message : 'Network error',
          status: 0,
          details:
            error instanceof Error && error.name === 'AbortError'
              ? 'Request timed out after 60 seconds'
              : 'Failed to connect to server',
        },
      }
    }
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
