import { setActivePinia, createPinia } from 'pinia'
import { useProfileStore } from '../profile'
import { useAuthStore } from '../auth'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { supabase } from '@/plugins/supabase'
import type { Mock } from 'vitest'

const mockedSupabase = supabase as unknown as {
  from: Mock
}

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
              id: '123',
            },
          },
        },
      }),
      getUser: vi.fn().mockResolvedValue({
        data: {
          user: {
            id: '123',
          },
        },
      }),
      onAuthStateChange: vi.fn(),
    },
  },
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
      created_at: new Date().toISOString(),
    }

    const mockProfile = {
      user_id: '123',
      username: 'testuser',
      experience: 'Beginner',
      preferred_strokes: ['Freestyle'],
      categories: ['Swimmer'],
    }

    mockedSupabase.from.mockReturnValue({
      select: vi.fn().mockReturnThis(),
      eq: vi.fn().mockReturnThis(),
      single: vi.fn().mockResolvedValue({ data: mockProfile, error: null }),
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
      created_at: new Date().toISOString(),
    }

    const updatedProfileData = {
      experience: 'Intermediate',
    }

    const mockUpdatedProfile = {
      user_id: '123',
      username: 'testuser',
      experience: 'Intermediate',
      preferred_strokes: ['Freestyle'],
      categories: ['Swimmer'],
    }

    mockedSupabase.from.mockReturnValue({
      update: vi.fn().mockReturnThis(),
      eq: vi.fn().mockReturnThis(),
      select: vi.fn().mockReturnThis(),
      single: vi.fn().mockResolvedValue({ data: mockUpdatedProfile, error: null }),
    })

    await profileStore.updateProfile(updatedProfileData)

    expect(profileStore.profile).toEqual(mockUpdatedProfile)
  })
})
