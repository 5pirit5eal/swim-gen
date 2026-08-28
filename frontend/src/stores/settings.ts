import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Filter, AudienceType } from '@/types'

export const useSettingsStore = defineStore('settings', () => {
  // Existing settings
  const dataDonationOptOut = ref(true)
  const poolLength = ref<25 | 50 | 'Freiwasser'>(25)
  const preferredMethod = ref<'choose' | 'generate'>('generate')
  const useProfilePreferences = ref(true)
  const selectedAudience = ref<AudienceType | null>(null)

  // Load audience from localStorage on init
  const storedAudience = localStorage.getItem('swim-gen-audience')
  if (
    storedAudience &&
    ['beginner', 'triathlete', 'competitive_swimmer', 'hobby'].includes(storedAudience)
  ) {
    selectedAudience.value = storedAudience as AudienceType
  }

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

  function setAudience(audience: AudienceType | null) {
    selectedAudience.value = audience
    try {
      if (audience) {
        localStorage.setItem('swim-gen-audience', audience)
      } else {
        localStorage.removeItem('swim-gen-audience')
      }
    } catch (e) {
      console.error('Failed to persist audience to localStorage', e)
    }
  }

  function toggleAudience(audience: AudienceType) {
    if (selectedAudience.value === audience) {
      setAudience(null)
    } else {
      setAudience(audience)
    }
  }

  return {
    // State
    dataDonationOptOut,
    poolLength,
    preferredMethod,
    useProfilePreferences,
    selectedAudience,
    filters,
    tutorials,
    // Actions
    updateStrokeFilter,
    clearFilters,
    markTutorialSeen,
    setAudience,
    toggleAudience,
  }
})
