import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Filter } from '@/types'

export const useSettingsStore = defineStore('settings', () => {
  // Existing settings
  const dataDonationOptOut = ref(true)
  const poolLength = ref<25 | 50 | 'Freiwasser'>(25)
  const preferredMethod = ref<'choose' | 'generate'>('generate')
  const useProfilePreferences = ref(true)

  // Tutorial settings
  const tutorials = ref({
    home: false,
    interaction: false,
    sidebar: false,
  })

  // Load from localStorage on init
  const storedTutorials = localStorage.getItem('swim-gen-tutorials')
  if (storedTutorials) {
    try {
      tutorials.value = { ...tutorials.value, ...JSON.parse(storedTutorials) }
    } catch (e) {
      console.error('Failed to parse stored tutorials settings', e)
    }
  }

  // Filter settings
  const filters = ref<Filter>({
    freistil: undefined,
    brust: undefined,
    ruecken: undefined,
    delfin: undefined,
    lagen: undefined,
    schwierigkeitsgrad: undefined,
    trainingstyp: undefined,
  })

  // Actions for tutorials
  function markTutorialSeen(tutorial: keyof typeof tutorials.value) {
    tutorials.value[tutorial] = true
    try {
      localStorage.setItem('swim-gen-tutorials', JSON.stringify(tutorials.value))
    } catch (e) {
      console.error('Failed to persist tutorial state to localStorage for', tutorial, e)
    }
  }

  // Actions for filters
  function updateStrokeFilter(
    stroke: keyof Pick<Filter, 'freistil' | 'brust' | 'ruecken' | 'delfin' | 'lagen'>,
    value: boolean | undefined,
  ) {
    filters.value[stroke] = value
  }

  function clearFilters() {
    filters.value = {
      freistil: undefined,
      brust: undefined,
      ruecken: undefined,
      delfin: undefined,
      lagen: undefined,
      schwierigkeitsgrad: undefined,
      trainingstyp: undefined,
    }
  }

  return {
    // State
    dataDonationOptOut,
    poolLength,
    preferredMethod,
    useProfilePreferences,
    filters,
    tutorials,
    // Actions
    updateStrokeFilter,
    clearFilters,
    markTutorialSeen,
  }
})
