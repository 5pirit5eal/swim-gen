/**
 * Training plan related TypeScript type definitions
 * Based on the backend API specification (swagger.yaml)
 */

// Backend API Row structure - represents a single exercise
export interface Row {
  Amount: number
  Break: string
  Content: string
  Distance: number
  Intensity: string
  Multiplier: string
  Sum: number
}

// Backend API QueryRequest structure
export interface QueryRequest {
  content: string
  filter?: Record<string, unknown>
  method?: 'choose' | 'generate'
}

// Backend API RAGResponse structure
export interface RAGResponse {
  title: string
  description: string
  table: Row[]
}

// Backend API PlanToPDFRequest structure
export interface PlanToPDFRequest {
  title: string
  description: string
  table: Row[]
}

// Backend API PlanToPDFResponse structure
export interface PlanToPDFResponse {
  uri: string
}

// Backend API DonatePlanRequest structure
export interface DonatePlanRequest {
  title: string
  description: string
  table: Row[]
  user_id: string
}
