import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useSettingsStore = defineStore('settings', () => {
  // Settings state
  const dataDonationOptOut = ref(false)
  const poolLength = ref<25 | 50>(25)
  const preferredMethod = ref<'choose' | 'generate'>('generate')

  // Actions
  function updateDataDonation(optOut: boolean) {
    dataDonationOptOut.value = optOut
  }

  function updatePoolLength(length: 25 | 50) {
    poolLength.value = length
  }

  function updateMethod(method: 'choose' | 'generate') {
    preferredMethod.value = method
  }

  return {
    // State
    dataDonationOptOut,
    poolLength,
    preferredMethod,
    // Actions
    updateDataDonation,
    updatePoolLength,
    updateMethod,
  }
})
