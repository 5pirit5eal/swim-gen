import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import type { Mock } from 'vitest'
import { useDrillsStore } from '@/stores/drills'
import { apiClient } from '@/api/client'
import type { ApiResult, Drill, DrillSearchResult } from '@/types'

vi.mock('@/api/client', async (importOriginal) => {
  const actual = (await importOriginal()) as typeof import('@/api/client')
  return {
    ...actual,
    apiClient: {
      searchDrills: vi.fn(),
      getDrill: vi.fn(),
      getDrillOptions: vi.fn(),
    },
  }
})

const mockedSearchDrills = apiClient.searchDrills as Mock

function makeDrill(i: number): Drill {
  return {
    slug: `drill-${i}`,
    targets: ['Technique'],
    short_description: `Short description ${i}`,
    img_name: `drill-${i}.webp`,
    img_description: `Drill ${i} illustration`,
    title: `Drill ${i} Title`,
    description: [`Full description ${i}`],
    video_url: [],
    styles: ['Freestyle'],
    difficulty: 'Easy',
    target_groups: ['Beginner'],
  }
}

function successResult(drills: Drill[], total = drills.length): ApiResult<DrillSearchResult> {
  return {
    success: true,
    data: { drills, total, page: 1, limit: 4 },
  }
}

describe('drills store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('fetchFeaturedDrills', () => {
    it('requests a random sample with the given limit', async () => {
      const drills = Array.from({ length: 4 }, (_, i) => makeDrill(i + 1))
      mockedSearchDrills.mockResolvedValue(successResult(drills, 44))

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('en', 4)

      expect(mockedSearchDrills).toHaveBeenCalledWith({
        lang: 'en',
        page: 1,
        limit: 4,
        random: true,
      })
      expect(store.featuredDrills).toEqual(drills)
      expect(store.featuredTotal).toBe(44)
    })

    it('caps the featured drills at the limit and caches them', async () => {
      const drills = Array.from({ length: 6 }, (_, i) => makeDrill(i + 1))
      mockedSearchDrills.mockResolvedValue(successResult(drills, 44))

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('de', 4)

      expect(mockedSearchDrills).toHaveBeenCalledWith({
        lang: 'de',
        page: 1,
        limit: 4,
        random: true,
      })
      expect(store.featuredDrills).toHaveLength(4)
      expect(store.drillCache.size).toBe(6)
      expect(store.drillCache.get('drill-1.webp:de')).toEqual(drills[0])
    })

    it('keeps featured drills empty and stops loading when the request fails', async () => {
      mockedSearchDrills.mockResolvedValue({
        success: false,
        error: { message: 'boom', status: 500 },
      })

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('en', 4)

      expect(store.featuredDrills).toEqual([])
      expect(store.featuredTotal).toBe(0)
      expect(store.isFeaturedLoading).toBe(false)
    })
  })
})
