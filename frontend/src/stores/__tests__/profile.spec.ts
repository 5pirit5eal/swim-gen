import { setActivePinia, createPinia } from 'pinia'
import { useProfileStore } from '../profile'
import { useAuthStore } from '../auth'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { supabase } from '@/plugins/supabase'
import type { Mock } from 'vitest'

const { mockedSupabaseClient } = vi.hoisted(() => ({
  mockedSupabaseClient: {
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

const mockedSupabase = supabase as unknown as {
  from: Mock
}

vi.mock('@/plugins/supabase', () => ({
  supabase: mockedSupabaseClient,
  getSupabase: vi.fn(async () => mockedSupabaseClient),
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

    const select = vi.fn().mockReturnThis()
    mockedSupabase.from.mockReturnValue({
      select,
      eq: vi.fn().mockReturnThis(),
      single: vi.fn().mockResolvedValue({ data: mockProfile, error: null }),
    })

    await profileStore.fetchProfile()

    expect(profileStore.profile).toEqual(mockProfile)
    expect(select).toHaveBeenCalledWith(
      'user_id, updated_at, username, experience, preferred_language, preferred_strokes, categories, overall_generations, monthly_generations, exports, css_200m_seconds, css_400m_seconds',
    )
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
      exports: 999999,
      user_id: 'forged-user',
    } as unknown as Parameters<typeof profileStore.updateProfile>[0]

    const mockUpdatedProfile = {
      user_id: '123',
      username: 'testuser',
      experience: 'Intermediate',
      preferred_strokes: ['Freestyle'],
      categories: ['Swimmer'],
    }

    const update = vi.fn().mockReturnThis()
    const select = vi.fn().mockReturnThis()
    mockedSupabase.from.mockReturnValue({
      update,
      eq: vi.fn().mockReturnThis(),
      select,
      single: vi.fn().mockResolvedValue({ data: mockUpdatedProfile, error: null }),
    })

    await profileStore.updateProfile(updatedProfileData)

    expect(profileStore.profile).toEqual(mockUpdatedProfile)
    expect(update).toHaveBeenCalledWith({ experience: 'Intermediate' })
    expect(select).toHaveBeenCalledWith(
      'user_id, updated_at, username, experience, preferred_language, preferred_strokes, categories, overall_generations, monthly_generations, exports, css_200m_seconds, css_400m_seconds',
    )
  })
})
