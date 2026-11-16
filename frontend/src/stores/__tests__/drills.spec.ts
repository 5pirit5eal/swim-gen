import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
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

function successResult(drills: Drill[]): ApiResult<DrillSearchResult> {
  return {
    success: true,
    data: { drills, total: drills.length, page: 1, limit: 100 },
  }
}

describe('drills store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('fetchFeaturedDrills', () => {
    it('requests the full pool and returns a random subset of the given limit', async () => {
      const pool = Array.from({ length: 20 }, (_, i) => makeDrill(i + 1))
      mockedSearchDrills.mockResolvedValue(successResult(pool))

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('en', 4)

      expect(mockedSearchDrills).toHaveBeenCalledWith({ lang: 'en', page: 1, limit: 100 })
      expect(store.featuredDrills).toHaveLength(4)
      expect(store.featuredTotal).toBe(20)
      store.featuredDrills.forEach((drill) => {
        expect(pool).toContainEqual(drill)
      })
    })

    it('caches all fetched drills', async () => {
      const pool = Array.from({ length: 10 }, (_, i) => makeDrill(i + 1))
      mockedSearchDrills.mockResolvedValue(successResult(pool))

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('de', 4)

      expect(store.drillCache.size).toBe(10)
      expect(store.drillCache.get('drill-1.webp:de')).toEqual(pool[0])
    })

    it('picks a deterministic subset when Math.random is mocked', async () => {
      vi.spyOn(Math, 'random').mockReturnValue(0)
      const pool = Array.from({ length: 6 }, (_, i) => makeDrill(i + 1))
      mockedSearchDrills.mockResolvedValue(successResult(pool))

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('en', 4)

      expect(store.featuredDrills.map((d) => d.slug)).toEqual([
        'drill-2',
        'drill-3',
        'drill-4',
        'drill-5',
      ])
    })

    it('returns fewer drills when the pool is smaller than the limit', async () => {
      const pool = [makeDrill(1), makeDrill(2)]
      mockedSearchDrills.mockResolvedValue(successResult(pool))

      const store = useDrillsStore()
      await store.fetchFeaturedDrills('en', 4)

      expect(store.featuredDrills).toHaveLength(2)
      expect(store.featuredTotal).toBe(2)
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
