import type { RAGResponse } from './training'

export interface SharedPlan {
    user_id: string
    plan_id: string
    created_at: string
    updated_at: string
    url_hash: string
    share_count: number
}

export interface SharedHistoryItem {
    user_id: string
    plan_id: string
    share_method: string
    shared_by: string
    created_at: string
    plan: RAGResponse
}

export interface ShareUrlRequest {
    plan_id: string
    method: 'link' | 'email'
}

export interface ShareUrlResponse {
    url_hash: string
}

export interface SharedPlanData {
    plan: RAGResponse
    sharer_username: string
    sharer_id: string
}
