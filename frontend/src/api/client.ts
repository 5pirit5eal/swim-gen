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
  type UpsertPlanRequest,
  type UpsertPlanResponse,
  type ShareUrlRequest,
  type ShareUrlResponse,
  type ChatRequest,
  type ChatResponsePayload,
  type FeedbackRequest,
  type DonatePlanRequest,
  type UploadedPlan,
  ApiEndpoints,
} from '@/types'
import i18n from '@/plugins/i18n'
import type { Message } from '@/types'
import { useAuthStore } from '@/stores/auth'

export function formatError(error: { message?: string; details?: string }): string {
  return `${error.message}: ${error.details ?? i18n.global.t('errors.unknown_error')}`
}

class ApiClient {
  private baseUrl: string
  public readonly DEFAULT_TIMEOUT_MS: number = 5000 // 5 seconds
  public readonly QUERY_TIMEOUT_MS: number = 60000 // 60 seconds
  public readonly PROMPT_TIMEOUT_MS: number = 10000 // 10 seconds

  constructor(baseUrl = '/api') {
    this.baseUrl = baseUrl
  }

  private async _getAuthToken(): Promise<string | null> {
    const authStore = useAuthStore()
    return authStore.session?.access_token ?? null
  }

  private async _fetch<T>(
    endpoint: string,
    options: RequestInit,
    timeout: number,
    authenticated = false,
  ): Promise<ApiResult<T>> {
    try {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), timeout)

      const headers = new Headers(options.headers)
      if (authenticated) {
        const token = await this._getAuthToken()
        if (token) {
          headers.set('Authorization', `Bearer ${token}`)
        }
      }

      let bodyLog: unknown = options.body
      if (typeof options.body === 'string') {
        try {
          bodyLog = JSON.parse(options.body)
        } catch {
          // ignore
        }
      } else if (options.body instanceof FormData) {
        bodyLog = 'FormData'
      }

      console.debug(
        `[API] Request: ${options.method || 'GET'} ${endpoint}`,
        bodyLog ? { body: bodyLog } : '',
      )

      const response = await fetch(`${this.baseUrl}/${endpoint}`, {
        ...options,
        headers,
        signal: controller.signal,
      })

      clearTimeout(timeoutId)

      console.debug(`[API] Response: ${endpoint} ${response.status}`)

      if (!response.ok) {
        return {
          success: false,
          error: {
            message: i18n.global.t('errors.api_request_failed', { endpoint }),
            status: response.status,
            details: response.statusText,
          },
        }
      }

      const contentType = response.headers.get('content-type')
      const data =
        contentType && contentType.includes('application/json')
          ? await response.json()
          : await response.text()
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
              ? i18n.global.t('errors.timeout', { time: timeout / 1000 })
              : i18n.global.t('errors.connection_failed'),
        },
      }
    }
  }

  /**
   * Check API health status
   */
  async checkHealth(): Promise<ApiResult<HealthCheckResponse>> {
    return this._fetch(ApiEndpoints.HEALTH, {}, this.DEFAULT_TIMEOUT_MS)
  }

  /**
   * Generate a random prompt for generating training plans
   */
  async generatePrompt(
    request: PromptGenerationRequest,
  ): Promise<ApiResult<PromptGenerationResponse>> {
    return this._fetch(
      ApiEndpoints.PROMPT,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.PROMPT_TIMEOUT_MS,
    )
  }

  /**
   * Query for training plans (may take up to 60 seconds)
   */
  async query(request: QueryRequest): Promise<ApiResult<RAGResponse>> {
    return this._fetch(
      ApiEndpoints.QUERY,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.QUERY_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Export training plan as PDF
   */
  async exportPDF(request: PlanToPDFRequest): Promise<ApiResult<PlanToPDFResponse>> {
    return this._fetch(
      ApiEndpoints.EXPORT_PDF,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Upsert (create or update) a training plan
   */
  async upsertPlan(plan: UpsertPlanRequest): Promise<ApiResult<UpsertPlanResponse>> {
    return this._fetch(
      ApiEndpoints.UPSERT_PLAN,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(plan),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Create a shareable URL for a training plan
   */
  async createShareUrl(request: ShareUrlRequest): Promise<ApiResult<ShareUrlResponse>> {
    return this._fetch(
      ApiEndpoints.SHARE,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Get conversation history for a plan
   */
  async getConversation(planId: string): Promise<ApiResult<Message[]>> {
    return this._fetch(
      `${ApiEndpoints.CONVERSATION}?plan_id=${planId}`,
      {
        method: 'GET',
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Add a plan to user history
   */
  async addPlanToHistory(
    plan: RAGResponse,
  ): Promise<ApiResult<{ message: string; plan_id: string }>> {
    return this._fetch(
      ApiEndpoints.ADD_PLAN_TO_HISTORY,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(plan),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Send a chat message to the AI
   */
  async chat(request: ChatRequest): Promise<ApiResult<ChatResponsePayload>> {
    return this._fetch(
      ApiEndpoints.CHAT,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.QUERY_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Submit feedback for a training plan
   */
  async submitFeedback(request: FeedbackRequest): Promise<ApiResult<string>> {
    return this._fetch(
      ApiEndpoints.FEEDBACK,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Upload a training plan
   */
  async donatePlan(request: DonatePlanRequest): Promise<ApiResult<string>> {
    return this._fetch(
      ApiEndpoints.ADD_PLAN,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Get all uploaded plans for the user
   */
  async getUploadedPlans(): Promise<ApiResult<UploadedPlan[]>> {
    return this._fetch(
      ApiEndpoints.GET_UPLOADS,
      {
        method: 'GET',
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Get a specific uploaded plan
   */
  async getUploadedPlan(planId: string): Promise<ApiResult<UploadedPlan>> {
    return this._fetch(
      `${ApiEndpoints.GET_UPLOADS}/${planId}`,
      {
        method: 'GET',
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Extract a training plan from a file (PNG, JPEG, or PDF)
   */
  async fileToPlan(file: File, language: string = 'en'): Promise<ApiResult<RAGResponse>> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('language', language)

    return this._fetch(
      ApiEndpoints.FILE_TO_PLAN,
      {
        method: 'POST',
        // Content-Type header is automatically set by the browser for FormData
        // We need to pass undefined so that the browser sets the boundary correctly
        headers: {},
        body: formData,
      },
      this.QUERY_TIMEOUT_MS, // This might take a while
      true,
    )
  }

  /**
   * Add a message to the conversation history
   */
  async addMessage(
    planId: string,
    role: 'user' | 'ai',
    content: string,
    previousMessageId?: string,
    planSnapshot?: RAGResponse,
  ): Promise<ApiResult<{ message_id: string }>> {
    return this._fetch<{ message_id: string }>(
      ApiEndpoints.ADD_MESSAGE,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          plan_id: planId,
          role,
          content,
          previous_message_id: previousMessageId,
          plan_snapshot: planSnapshot,
        }),
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Delete a training plan
   */
  async deletePlan(planId: string): Promise<ApiResult<string>> {
    return this._fetch(
      `plan/${planId}`,
      {
        method: 'DELETE',
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }

  /**
   * Delete user account and all associated data
   */
  async deleteUser(): Promise<ApiResult<string>> {
    return this._fetch(
      'user',
      {
        method: 'DELETE',
      },
      this.DEFAULT_TIMEOUT_MS,
      true,
    )
  }
}

// Export singleton instance
export const apiClient = new ApiClient()
