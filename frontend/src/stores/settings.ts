import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Filter } from '@/types'

export const useSettingsStore = defineStore('settings', () => {
  // Existing settings
  const dataDonationOptOut = ref(false)
  const poolLength = ref<25 | 50>(25)
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

  // Actions for existing settings
  function updateDataDonation(optOut: boolean) {
    dataDonationOptOut.value = optOut
  }

  function updatePoolLength(length: 25 | 50) {
    poolLength.value = length
  }

  function updateMethod(method: 'choose' | 'generate') {
    preferredMethod.value = method
  }

  // Actions for filters
  function updateStrokeFilter(
    stroke: keyof Pick<Filter, 'freistil' | 'brust' | 'ruecken' | 'delfin' | 'lagen'>,
    value: boolean | undefined,
  ) {
    filters.value[stroke] = value
  }

  function updateDifficultyFilter(difficulty: Filter['schwierigkeitsgrad']) {
    filters.value.schwierigkeitsgrad = difficulty
  }

  function updateTrainingTypeFilter(type: Filter['trainingstyp']) {
    filters.value.trainingstyp = type
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
    updateDataDonation,
    updatePoolLength,
    updateMethod,
    updateStrokeFilter,
    updateDifficultyFilter,
    updateTrainingTypeFilter,
    clearFilters,
  }
})
