import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSidebarStore } from '../sidebar'

describe('sidebar store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should have a default state of closed', () => {
    const store = useSidebarStore()
    expect(store.isOpen).toBe(false)
  })

  it('should open the sidebar', () => {
    const store = useSidebarStore()
    store.open()
    expect(store.isOpen).toBe(true)
  })

  it('should close the sidebar', () => {
    const store = useSidebarStore()
    store.isOpen = true
    store.close()
    expect(store.isOpen).toBe(false)
  })

  it('should toggle the sidebar', () => {
    const store = useSidebarStore()
    expect(store.isOpen).toBe(false)
    store.toggle()
    expect(store.isOpen).toBe(true)
    store.toggle()
    expect(store.isOpen).toBe(false)
  })
})
