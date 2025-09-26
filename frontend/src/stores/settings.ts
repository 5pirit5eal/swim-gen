import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Filter } from '@/types'

export const useSettingsStore = defineStore('settings', () => {
  // Existing settings
  const dataDonationOptOut = ref(true)
  const poolLength = ref<25 | 50 | 'Freiwasser'>(25)
  const preferredMethod = ref<'choose' | 'generate'>('generate')

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
    filters,
    // Actions
    updateStrokeFilter,
    clearFilters,
  }
})
