import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUploadFormStore } from '@/stores/uploadForm'

import { apiClient } from '@/api/client'

// Mock dependencies
vi.mock('@/api/client', () => ({
  apiClient: {
    donatePlan: vi.fn(),
  },
  formatError: vi.fn((error) => error.message),
}))

const mockFetchUploadedPlans = vi.fn()

vi.mock('@/stores/uploads', () => ({
  useDonationStore: vi.fn(() => ({
    fetchUploadedPlans: mockFetchUploadedPlans,
  })),
}))

vi.mock('@/plugins/i18n', () => ({
  default: {
    global: {
      t: (key: string) => key,
      locale: { value: 'en' },
    },
  },
}))

describe('uploadForm Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('initializes with default state', () => {
    const store = useUploadFormStore()
    expect(store.currentPlan).toBeNull()
    expect(store.isLoading).toBe(false)
    expect(store.error).toBeNull()
    expect(store.hasPlan).toBe(false)
  })

  it('initializes a new plan', () => {
    const store = useUploadFormStore()
    store.initNewPlan()
    expect(store.currentPlan).not.toBeNull()
    expect(store.currentPlan?.title).toBe('donation.newPlan.title')
    expect(store.currentPlan?.table).toHaveLength(2)
  })

  it('uploads a plan successfully', async () => {
    const store = useUploadFormStore()

    store.initNewPlan()
    // Mock successful API response
    ;(apiClient.donatePlan as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({
      success: true,
    })

    const result = await store.uploadCurrentPlan()

    expect(result).toBe(true)
    expect(apiClient.donatePlan).toHaveBeenCalled()
    expect(mockFetchUploadedPlans).toHaveBeenCalled()
    expect(store.isLoading).toBe(false)
  })

  it('handles upload failure', async () => {
    const store = useUploadFormStore()

    store.initNewPlan()
    // Mock failed API response
    ;(apiClient.donatePlan as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({
      success: false,
      error: { message: 'Upload failed' },
    })

    const result = await store.uploadCurrentPlan()

    expect(result).toBe(false)
    expect(store.error).toBe('Upload failed')
    expect(store.isLoading).toBe(false)
  })

  it('clears the state', () => {
    const store = useUploadFormStore()
    store.initNewPlan()
    store.clear()
    expect(store.currentPlan).toBeNull()
    expect(store.error).toBeNull()
    expect(store.isLoading).toBe(false)
  })
})
