/**
 * API client related TypeScript type definitions
 * Based on the backend API specification (swagger.yaml)
 */

// HTTP methods supported by our API
export type HttpMethod = 'GET' | 'POST'

// Base API response structure (for error cases)
export interface ApiErrorResponse {
  error: string
  status: number
}

export interface ApiResult<T> {
  success: boolean
  data?: T
  error?: {
    message: string
    status: number
    details?: string
  }
}

// Health check response (simple string)
export type HealthCheckResponse = string

// API endpoints enum for type safety
export enum ApiEndpoints {
  HEALTH = 'health',
  PROMPT = 'prompt',
  QUERY = 'query',
  EXPORT_PDF = 'export-pdf',
  ADD_PLAN = 'add',
  SCRAPE = 'scrape',
  UPSERT_PLAN = 'upsert-plan',
  SHARE = 'share-plan',
  CONVERSATION = 'memory/conversation',
  ADD_MESSAGE = 'memory/message',
  ADD_PLAN_TO_HISTORY = 'add-plan-to-history',
  CHAT = 'chat',
  FEEDBACK = 'feedback',
  GET_UPLOADS = 'uploads',
  FILE_TO_PLAN = 'file-to-plan',
  DRILL = 'drill',
  DRILL_SEARCH = 'drills/search',
  DRILL_OPTIONS = 'drills/options',
}

// Form validation types
export interface ValidationError {
  field: string
  message: string
  code: string
}

export interface ValidationResult {
  isValid: boolean
  errors: ValidationError[]
}
