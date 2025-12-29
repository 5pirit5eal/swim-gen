/**
 * Drill exercise types
 * Based on the data structure from data/drills/en.json
 */

export interface Drill {
  slug: string
  targets: string[]
  short_description: string
  img_name: string
  img_description: string
  title: string
  description: string[]
  video_url: string[]
  styles: string[]
  difficulty: string // Localized: Easy/Medium/Hard in EN, Leicht/Mittel/Schwer in DE
  target_groups: string[]
}

/**
 * Request parameters for fetching a drill
 */
export interface DrillRequest {
  id: string // img_name without extension
  lang: string
}

/**
 * Minimal drill info for preview cards
 */
export interface DrillPreview {
  img_name: string
  title: string
  short_description: string
  difficulty: string
  target?: string
  style?: string
}

export interface DrillSearchParams {
  lang: string
  q?: string
  difficulty?: string
  target_groups?: string[]
  styles?: string[]
  page: number
  limit: number
}

export interface DrillSearchResult {
  drills: Drill[]
  total: number
  page: number
  limit: number
}

export interface DrillFilterOptions {
  styles: string[]
  target_groups: string[]
  targets: string[]
  difficulties: string[]
}
