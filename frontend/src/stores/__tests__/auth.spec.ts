import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'
import { supabase } from '@/plugins/supabase'
import type { Mock } from 'vitest'

// Mock the supabase client
vi.mock('@/plugins/supabase', () => ({
  supabase: {
    auth: {
      getSession: vi.fn(),
      getUser: vi.fn(),
      onAuthStateChange: vi.fn(),
      signInWithPassword: vi.fn(),
      signUp: vi.fn(),
      signOut: vi.fn(),
      refreshSession: vi.fn(),
    },
    from: vi.fn(() => ({
      select: vi.fn(() => ({
        eq: vi.fn(() => ({
          maybeSingle: vi.fn().mockResolvedValue({ data: null, error: null }),
        })),
      })),
    })),
  },
}))

const mockedGetSession = supabase.auth.getSession as Mock
const mockedGetUser = supabase.auth.getUser as Mock
const mockedOnAuthStateChange = supabase.auth.onAuthStateChange as Mock
const mockedSignInWithPassword = supabase.auth.signInWithPassword as Mock
const mockedSignUp = supabase.auth.signUp as Mock
const mockedSignOut = supabase.auth.signOut as Mock
const mockedRefreshSession = supabase.auth.refreshSession as Mock

describe('auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('initializes with no session or user', () => {
    mockedGetSession.mockResolvedValue({ data: { session: null } })
    mockedGetUser.mockResolvedValue({ data: { user: null } })

    const store = useAuthStore()
    expect(store.session).toBe(null)
    expect(store.user).toBe(null)
  })

  it('signInWithPassword successfully signs in a user', async () => {
    const store = useAuthStore()
    const mockData = { user: { id: '123' }, session: { access_token: 'abc' } }
    mockedSignInWithPassword.mockResolvedValue({ data: mockData, error: null })
    mockedRefreshSession.mockResolvedValue({ data: { session: mockData.session }, error: null })
    mockedGetUser.mockResolvedValue({ data: { user: mockData.user }, error: null })

    const result = await store.signInWithPassword('test@example.com', 'password')

    expect(result).toEqual(mockData)
    expect(mockedSignInWithPassword).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password',
    })
  })

  it('signInWithPassword throws an error on failure', async () => {
    const store = useAuthStore()
    const mockError = new Error('Sign in failed')
    mockedSignInWithPassword.mockResolvedValue({ data: null, error: mockError })

    await expect(store.signInWithPassword('test@example.com', 'password')).rejects.toThrow(
      mockError,
    )
  })

  it('signUp successfully signs up a user', async () => {
    const store = useAuthStore()
    const mockData = { user: { id: '456' }, session: { access_token: 'def' } }
    mockedSignUp.mockResolvedValue({ data: mockData, error: null })

    const result = await store.signUp('new@example.com', 'newpassword', 'newuser')

    expect(result).toEqual(mockData)
    expect(mockedSignUp).toHaveBeenCalledWith({
      email: 'new@example.com',
      password: 'newpassword',
      options: {
        data: {
          username: 'newuser',
        },
        emailRedirectTo: `${window.location.origin}/login`,
      },
    })
  })

  it('signUp throws an error on failure', async () => {
    const store = useAuthStore()
    const mockError = new Error('Sign up failed')
    mockedSignUp.mockResolvedValue({ data: null, error: mockError })

    await expect(store.signUp('new@example.com', 'newpassword', 'newuser')).rejects.toThrow(
      mockError,
    )
  })

  it('signOut successfully signs out a user', async () => {
    const store = useAuthStore()
    mockedSignOut.mockResolvedValue({ error: null })

    await store.signOut()

    expect(mockedSignOut).toHaveBeenCalled()
  })

  it('signOut throws an error on failure', async () => {
    const store = useAuthStore()
    const mockError = new Error('Sign out failed')
    mockedSignOut.mockResolvedValue({ error: mockError })

    await expect(store.signOut()).rejects.toThrow(mockError)
  })

  it('updates session and user on onAuthStateChange', () => {
    const store = useAuthStore()
    expect(mockedOnAuthStateChange).toHaveBeenCalled()

    const callback = mockedOnAuthStateChange.mock.calls[0][0]
    const newSession = {
      access_token: 'new-token',
      user: { id: 'user-id' },
    }

    // Simulate an auth state change
    callback('SIGNED_IN', newSession)

    expect(store.session).toEqual(newSession)
    expect(store.user).toEqual(newSession.user)

    // Simulate signing out
    callback('SIGNED_OUT', null)

    expect(store.session).toBe(null)
    expect(store.user).toBe(null)
  })
})
