export interface HistoryMetadata {
  plan_id: string
  keep_forever: boolean
  created_at: string
  updated_at: string
  shared?: boolean
  shared_count?: number
  feedback_rating?: number
  was_swam?: boolean
  difficulty_rating?: number
  exported_at?: string
}
