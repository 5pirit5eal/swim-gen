/**
 * Training plan related TypeScript type definitions
 * Based on the backend API specification (swagger.yaml)
 */

// Prompt Generation API Request structure
export interface PromptGenerationRequest {
  language: string
}

// PromptGeneration API Response structure
export interface PromptGenerationResponse {
  prompt: string
}

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
  language: string
  pool_length?: 25 | 50 | 'Freiwasser'
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

// Backend filter structure (matching the Pydantic model)
export interface Filter {
  freistil?: boolean // Freestyle swimming technique
  brust?: boolean // Breaststroke swimming technique
  ruecken?: boolean // Backstroke swimming technique
  delfin?: boolean // Butterfly swimming technique
  lagen?: boolean // Individual medley swimming
  schwierigkeitsgrad?:
    | 'Nichtschwimmer'
    | 'Anfaenger'
    | 'Fortgeschritten'
    | 'Leistungsschwimmer'
    | 'Top-Athlet'
  trainingstyp?:
    | 'Techniktraining'
    | 'Leistungstest'
    | 'Grundlagenausdauer'
    | 'Recovery'
    | 'Kurzstrecken'
    | 'Langstrecken'
    | 'Atemmangel'
    | 'Wettkampfvorbereitung'
}

// Helper type for difficulty options
export const DIFFICULTY_OPTIONS = [
  { value: 'Anfaenger', label: 'form.difficulty_beginner' },
  { value: 'Fortgeschritten', label: 'form.difficulty_advanced' },
  { value: 'Leistungsschwimmer', label: 'form.difficulty_competitive_swimmer' },
] as const

// Helper type for training type options
export const TRAINING_TYPE_OPTIONS = [
  { value: 'Techniktraining', label: 'form.training_type_technique_training' },
  { value: 'Leistungstest', label: 'form.training_type_performance_test' },
  { value: 'Grundlagenausdauer', label: 'form.training_type_base_endurance' },
  { value: 'Recovery', label: 'form.training_type_recovery' },
  { value: 'Kurzstrecken', label: 'form.training_type_sprint' },
  { value: 'Langstrecken', label: 'form.training_type_distance' },
  { value: 'Atemmangel', label: 'form.training_type_breath_control' },
  { value: 'Wettkampfvorbereitung', label: 'form.training_type_race_preparation' },
] as const
