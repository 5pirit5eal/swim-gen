import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUploadStore } from '@/stores/uploads'

import { apiClient } from '@/api/client'

// Mock dependencies
vi.mock('@/api/client', () => ({
    apiClient: {
        getUploadedPlans: vi.fn().mockResolvedValue({ success: true, data: [] }),
        getUploadedPlan: vi.fn().mockResolvedValue({ success: true, data: {} }),
    },
    formatError: vi.fn((error) => error.message),
}))

vi.mock('@/stores/auth', () => ({
    useAuthStore: vi.fn(() => ({
        user: { id: 'user-1' },
    })),
}))

describe('uploads Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    it('fetches uploaded plans', async () => {
        const store = useUploadStore()
        const mockPlans = [{ plan_id: '1', title: 'Plan 1' }]
            ; (apiClient.getUploadedPlans as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({ success: true, data: mockPlans })

        await store.fetchUploadedPlans()

        expect(store.uploadedPlans).toEqual(mockPlans)
        expect(store.isFetchingUploads).toBe(false)
    })

    it('fetches a single uploaded plan', async () => {
        const store = useUploadStore()
        const mockPlan = { plan_id: '1', title: 'Plan 1', table: [] }
            ; (apiClient.getUploadedPlan as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({ success: true, data: mockPlan })

        const result = await store.fetchUploadedPlan('1')

        expect(result).toBe(true)
        expect(store.currentPlan?.title).toBe('Plan 1')
        expect(store.isLoading).toBe(false)
    })

    it('handles fetch error', async () => {
        const store = useUploadStore()
            ; (apiClient.getUploadedPlan as unknown as ReturnType<typeof vi.fn>).mockResolvedValue({
                success: false,
                error: { message: 'Not found' }
            })

        const result = await store.fetchUploadedPlan('1')

        expect(result).toBe(false)
        expect(store.error).toBe('Not found')
    })
})
