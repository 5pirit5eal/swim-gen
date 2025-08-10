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
  ApiEndpoints,
} from '@/types'

class ApiClient {
  private baseUrl: string

  constructor(baseUrl = 'http://localhost:8080') {
    this.baseUrl = baseUrl
  }

  /**
   * Check API health status
   */
  async checkHealth(): Promise<HealthCheckResponse | null> {
    try {
      const response = await fetch(`${this.baseUrl}/health`)
      if (!response.ok) return null
      return await response.text()
    } catch {
      return null
    }
  }

  /**
   * Query for training plans
   */
  async query(request: QueryRequest): Promise<RAGResponse | null> {
    try {
      const response = await fetch(`${this.baseUrl}/query`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      })

      if (!response.ok) return null
      return await response.json()
    } catch {
      return null
    }
  }

  /**
   * Export training plan as PDF
   */
  async exportPDF(request: PlanToPDFRequest): Promise<PlanToPDFResponse | null> {
    try {
      const response = await fetch(`${this.baseUrl}/export-pdf`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      })

      if (!response.ok) return null
      return await response.json()
    } catch {
      return null
    }
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
