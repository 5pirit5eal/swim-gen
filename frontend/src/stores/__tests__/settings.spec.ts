// frontend/src/stores/__tests__/settings.spec.ts
import { describe, it, expect, beforeEach } from 'vitest'
import { useSettingsStore } from '@/stores/settings'

describe('settings Store', () => {
  beforeEach(() => {
    // Manually reset the state for a setup store
    const store = useSettingsStore()
    store.dataDonationOptOut = true
    store.poolLength = 25
    store.preferredMethod = 'generate'
    store.filters = {
      freistil: undefined,
      brust: undefined,
      ruecken: undefined,
      delfin: undefined,
      lagen: undefined,
      schwierigkeitsgrad: undefined,
      trainingstyp: undefined,
    }
  })

  it('has the correct initial state', () => {
    const store = useSettingsStore()

    expect(store.dataDonationOptOut).toBe(true)
    expect(store.poolLength).toBe(25)
    expect(store.preferredMethod).toBe('generate')
    expect(store.filters).toEqual({
      freistil: undefined,
      brust: undefined,
      ruecken: undefined,
      delfin: undefined,
      lagen: undefined,
      schwierigkeitsgrad: undefined,
      trainingstyp: undefined,
    })
  })
  it('updates a stroke filter correctly', () => {
    const store = useSettingsStore()

    // Test setting a stroke to true
    store.updateStrokeFilter('freistil', true)
    expect(store.filters.freistil).toBe(true)

    // Test setting another stroke to true
    store.updateStrokeFilter('brust', true)
    expect(store.filters.brust).toBe(true)

    // Test setting a stroke to undefined (clearing it)
    store.updateStrokeFilter('freistil', undefined)
    expect(store.filters.freistil).toBeUndefined()
  })
  it('clears all filters correctly', () => {
    const store = useSettingsStore()

    // First, set some filters to non-default values
    store.updateStrokeFilter('freistil', true)
    store.updateStrokeFilter('brust', true)
    store.filters.schwierigkeitsgrad = 'Anfaenger'
    store.filters.trainingstyp = 'Techniktraining'

    // Assert that they are set
    expect(store.filters.freistil).toBe(true)
    expect(store.filters.brust).toBe(true)
    expect(store.filters.schwierigkeitsgrad).toBe('Anfaenger')
    expect(store.filters.trainingstyp).toBe('Techniktraining')

    // Now, call clearFilters
    store.clearFilters()

    // Assert that all filters are reset to undefined
    expect(store.filters.freistil).toBeUndefined()
    expect(store.filters.brust).toBeUndefined()
    expect(store.filters.ruecken).toBeUndefined()
    expect(store.filters.delfin).toBeUndefined()
    expect(store.filters.lagen).toBeUndefined()
    expect(store.filters.schwierigkeitsgrad).toBeUndefined()
    expect(store.filters.trainingstyp).toBeUndefined()
  })
})
