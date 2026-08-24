export interface Profile {
  user_id: string
  updated_at: string
  username: string
  experience: string | null
  preferred_language: string | null
  preferred_strokes: string[]
  categories: string[]
  overall_generations: number
  monthly_generations: number
  exports: number
  css_200m_seconds: number | null
  css_400m_seconds: number | null
}

export interface PublicProfile {
  username: string
}

export type ProfileUpdate = Partial<
  Pick<
    Profile,
    | 'username'
    | 'experience'
    | 'preferred_language'
    | 'preferred_strokes'
    | 'categories'
    | 'css_200m_seconds'
    | 'css_400m_seconds'
  >
>
