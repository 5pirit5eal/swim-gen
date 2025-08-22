/**
 * API Client for Swim Gen/RAG Backend
 * Handles HTTP requests with proper TypeScript typing
 */
import {
  type ApiResult,
  type HealthCheckResponse,
  type PlanToPDFRequest,
  type PlanToPDFResponse,
  type PromptGenerationRequest,
  type PromptGenerationResponse,
  type QueryRequest,
  type RAGResponse,
  ApiEndpoints,
} from '@/types';
import i18n from '@/plugins/i18n';

export function formatError(error: { message?: string; details?: string }): string {
  return `${error.message}: ${error.details ?? i18n.global.t('errors.unknown_error')}`
}

class ApiClient {
  private baseUrl: string
  public readonly DEFAUTL_TIMEOUT_MS: number = 5000 // 5 seconds
  public readonly QUERY_TIMEOUT_MS: number = 60000 // 60 seconds
  public readonly PROMPT_TIMEOUT_MS: number = 10000 // 10 seconds

  constructor(baseUrl = '/api') {
    this.baseUrl = baseUrl
  }

  /**
   * Check API health status
   */
  async checkHealth(): Promise<ApiResult<HealthCheckResponse>> {
    try {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), this.DEFAUTL_TIMEOUT_MS)

      const response = await fetch(`${this.baseUrl}/${ApiEndpoints.HEALTH}`, {
        signal: controller.signal,
      })

      clearTimeout(timeoutId)

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: i18n.global.t('errors.health_check_failed'),
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
          message: error instanceof Error ? error.message : i18n.global.t('errors.unknown_error'),
          status: 0,
          details: i18n.global.t('errors.connection_failed'),
        },
      }
    }
  }

  /**
   * Generate a random prompt for generating training plans
   */
  async generatePrompt(
    request: PromptGenerationRequest,
  ): Promise<ApiResult<PromptGenerationResponse>> {
    try {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), this.PROMPT_TIMEOUT_MS)

      const response = await fetch(`${this.baseUrl}/${ApiEndpoints.PROMPT}`, {
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
            message: i18n.global.t('errors.failed_to_generate_prompt'),
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
          message: error instanceof Error ? error.message : i18n.global.t('errors.unknown_error'),
          status: 0,
          details:
            error instanceof Error && error.name === 'AbortError'
              ? i18n.global.t('errors.timeout', { time: this.PROMPT_TIMEOUT_MS / 1000 })
              : i18n.global.t('errors.connection_failed'),
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
      const timeoutId = setTimeout(() => controller.abort(), this.QUERY_TIMEOUT_MS)

      const response = await fetch(`${this.baseUrl}/${ApiEndpoints.QUERY}`, {
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
            message: i18n.global.t('errors.training_plan_failed'),
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
          message: error instanceof Error ? error.message : i18n.global.t('errors.unknown_error'),
          status: 0,
          details:
            error instanceof Error && error.name === 'AbortError'
              ? i18n.global.t('errors.timeout', { time: this.QUERY_TIMEOUT_MS / 1000 })
              : i18n.global.t('errors.connection_failed'),
        },
      }
    }
  }

  /**
   * Export training plan as PDF
   */
  async exportPDF(request: PlanToPDFRequest): Promise<ApiResult<PlanToPDFResponse>> {
    try {
      const response = await fetch(`${this.baseUrl}/${ApiEndpoints.EXPORT_PDF}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: i18n.global.t('errors.failed_to_export_plan'),
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
              ? i18n.global.t('errors.timeout', { time: this.DEFAUTL_TIMEOUT_MS / 1000 })
              : i18n.global.t('errors.connection_failed'),
        },
      }
    }
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
