/**
 * Main types index file - exports all type definitions
 * This provides a central place to import types from
 */

// Export all training-related types (backend API structures)
export type {
  PromptGenerationRequest,
  PromptGenerationResponse,
  Row,
  QueryRequest,
  RAGResponse,
  PlanToPDFRequest,
  PlanToPDFResponse,
  DonatePlanRequest,
  Filter,
  UpsertPlanRequest,
  UpsertPlanResponse,
  PlanStore,
  ChatRequest,
  ChatResponsePayload,
  FeedbackRequest,
} from './training'

// Export filter option constants
export { DIFFICULTY_OPTIONS, TRAINING_TYPE_OPTIONS } from './training'

// Export all API-related types
export type {
  HttpMethod,
  ApiErrorResponse,
  ApiResult,
  HealthCheckResponse,
  ValidationError,
  ValidationResult,
} from './api'

// Export API endpoints enum
export { ApiEndpoints } from './api'

// Frontend-specific utility types
export interface SelectOption {
  value: string
  label: string
  disabled?: boolean
}

// Application state types
export interface AppState {
  isLoading: boolean
  error: string | null
  theme: 'light' | 'dark' | 'auto'
}

// Export user profile type
export type { Profile } from './profile'

// Export miscellaneous types
export type { HistoryMetadata } from './misc'
export type { Message, Role } from './message'

// Export all shared-plan-related types
export type {
  SharedPlan,
  SharedHistoryItem,
  ShareUrlRequest,
  ShareUrlResponse,
  SharedPlanData,
} from './shared'
