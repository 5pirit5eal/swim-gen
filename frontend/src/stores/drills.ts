import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiClient, formatError } from '@/api/client'
import type { Drill, DrillPreview } from '@/types'

export const useDrillsStore = defineStore('drills', () => {
  // State
  const currentDrill = ref<Drill | null>(null)
  const drillCache = ref<Map<string, Drill>>(new Map())
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const hasCurrentDrill = computed(() => currentDrill.value !== null)

  const drillPreview = computed<DrillPreview | null>(() => {
    if (!currentDrill.value) return null
    return {
      img_name: currentDrill.value.img_name,
      title: currentDrill.value.title,
      short_description: currentDrill.value.short_description,
      difficulty: currentDrill.value.difficulty,
    }
  })

  // Actions
  function getCacheKey(id: string, lang: string): string {
    return `${id}:${lang}`
  }

  async function fetchDrill(id: string, lang: string): Promise<Drill | null> {
    const cacheKey = getCacheKey(id, lang)

    // Check cache first
    if (drillCache.value.has(cacheKey)) {
      const cached = drillCache.value.get(cacheKey)!
      currentDrill.value = cached
      return cached
    }

    isLoading.value = true
    error.value = null

    try {
      const result = await apiClient.getDrill(id, lang)
      if (result.success && result.data) {
        drillCache.value.set(cacheKey, result.data)
        currentDrill.value = result.data
        return result.data
      } else {
        error.value = result.error ? formatError(result.error) : 'Failed to fetch drill'
        return null
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function fetchDrillPreview(id: string, lang: string): Promise<DrillPreview | null> {
    const drill = await fetchDrill(id, lang)
    if (!drill) return null
    return {
      img_name: drill.img_name,
      title: drill.slug,
      short_description: drill.short_description,
      difficulty: drill.difficulty,
      target: drill.targets?.[0],
      style: drill.styles?.[0],
    }
  }

  function clearCurrentDrill(): void {
    currentDrill.value = null
    error.value = null
  }

  function clearCache(): void {
    drillCache.value.clear()
  }

  function $reset(): void {
    currentDrill.value = null
    drillCache.value.clear()
    isLoading.value = false
    error.value = null
  }

  return {
    // State
    currentDrill,
    drillCache,
    isLoading,
    error,
    // Getters
    hasCurrentDrill,
    drillPreview,
    // Actions
    fetchDrill,
    fetchDrillPreview,
    clearCurrentDrill,
    clearCache,
    $reset,
  }
})
