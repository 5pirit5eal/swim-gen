/**
 * Message type matching backend models.Message structure
 */
export type Role = 'user' | 'ai'

export interface Message {
    id: string
    plan_id: string
    user_id: string
    role: Role
    content: string
    previous_message_id: string | null
    next_message_id: string | null
    plan_snapshot?: {
        plan_id: string
        title: string
        description: string
        table: {
            Amount: number
            Multiplier: string
            Distance: number
            Break: string
            Content: string
            Intensity: string
            Sum: number
        }[]
    }
    created_at: string
}
