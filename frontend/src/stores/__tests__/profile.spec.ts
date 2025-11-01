import { setActivePinia, createPinia } from 'pinia'
import { useProfileStore } from '../profile'
import { useAuthStore } from '../auth'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { supabase } from '@/plugins/supabase'

vi.mock('@/plugins/supabase', () => ({
    supabase: {
        from: vi.fn().mockReturnThis(),
        select: vi.fn().mockReturnThis(),
        eq: vi.fn().mockReturnThis(),
        update: vi.fn().mockReturnThis(),
        single: vi.fn(),
        auth: {
            getSession: vi.fn().mockResolvedValue({
                data: {
                    session: {
                        user: {
                            id: '123'
                        }
                    }
                }
            }),
            getUser: vi.fn().mockResolvedValue({
                data: {
                    user: {
                        id: '123'
                    }
                }
            }),
            onAuthStateChange: vi.fn()
        }
    }
}))

describe('Profile Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('should fetch a profile', async () => {
        const profileStore = useProfileStore()
        const authStore = useAuthStore()
        authStore.user = {
            id: '123',
            app_metadata: {},
            user_metadata: {},
            aud: 'authenticated',
            created_at: new Date().toISOString()
        }

        const mockProfile = {
            user_id: '123',
            username: 'testuser',
            experience: 'Beginner',
            preferred_strokes: ['Freestyle'],
            categories: ['Swimmer']
        }

        vi.mocked(supabase.from('profiles').select().eq('user_id', '123').single).mockResolvedValueOnce({
            data: mockProfile,
            error: null,
            count: 1,
            status: 200,
            statusText: 'OK'
        })

        await profileStore.fetchProfile()

        expect(profileStore.profile).toEqual(mockProfile)
    })

    it('should update a profile', async () => {
        const profileStore = useProfileStore()
        const authStore = useAuthStore()
        authStore.user = {
            id: '123',
            app_metadata: {},
            user_metadata: {},
            aud: 'authenticated',
            created_at: new Date().toISOString()
        }

        const updatedProfileData = {
            experience: 'Intermediate'
        }

        const mockUpdatedProfile = {
            user_id: '123',
            username: 'testuser',
            experience: 'Intermediate',
            preferred_strokes: ['Freestyle'],
            categories: ['Swimmer']
        }

        vi.mocked(
            supabase.from('profiles').update(updatedProfileData).eq('user_id', '123').select().single
        ).mockResolvedValueOnce({
            data: mockUpdatedProfile,
            error: null,
            count: 1,
            status: 200,
            statusText: 'OK'
        })

        await profileStore.updateProfile(updatedProfileData)

        expect(profileStore.profile).toEqual(mockUpdatedProfile)
    })
})
