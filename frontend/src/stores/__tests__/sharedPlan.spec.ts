import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSharedPlanStore } from '@/stores/sharedPlan'
import { apiClient } from '@/api/client'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'
import type { Mock } from 'vitest'

// --- Mocks ---
vi.mock('@/api/client', async (importOriginal) => {
  const actual = (await importOriginal()) as typeof import('@/api/client')
  return {
    ...actual,
    apiClient: {
      upsertPlan: vi.fn(),
    },
    formatError: vi.fn((error) => `${error.message}: ${error.details}`),
  }
})

vi.mock('@/plugins/supabase', () => {
  const mockSupabase = {
    from: vi.fn(),
    select: vi.fn(),
    order: vi.fn(),
    in: vi.fn(),
    eq: vi.fn(),
    single: vi.fn(),
    insert: vi.fn(),
    limit: vi.fn(),
    update: vi.fn(),
    maybeSingle: vi.fn(),
  }

  // Setup chaining
  mockSupabase.from.mockReturnValue(mockSupabase)
  mockSupabase.select.mockReturnValue(mockSupabase)
  mockSupabase.order.mockReturnValue(mockSupabase)
  mockSupabase.in.mockReturnValue(mockSupabase)
  mockSupabase.eq.mockReturnValue(mockSupabase)
  mockSupabase.single.mockReturnValue(mockSupabase)
  mockSupabase.insert.mockReturnValue(mockSupabase)
  mockSupabase.limit.mockReturnValue(mockSupabase)
  mockSupabase.update.mockReturnValue(mockSupabase)
  mockSupabase.maybeSingle.mockReturnValue(mockSupabase)

  return {
    supabase: mockSupabase,
  }
})

vi.mock('@/stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    user: { id: 'test-user-id' },
  })),
}))

vi.mock('@/stores/trainingPlan', () => ({
  useTrainingPlanStore: vi.fn(() => ({
    fetchHistory: vi.fn(),
    planHistory: [],
    loadPlanFromHistory: vi.fn(),
  })),
}))

// --- Mock Casts ---
const mockedApiUpsertPlan = apiClient.upsertPlan as Mock
const mockedSupabase = supabase as unknown as {
  from: Mock
  select: Mock
  order: Mock
  in: Mock
  eq: Mock
  single: Mock
  insert: Mock
  limit: Mock
  update: Mock
  maybeSingle: Mock
}
const mockedAuthStore = useAuthStore as unknown as Mock

describe('sharedPlan Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.resetAllMocks()
    mockedAuthStore.mockReturnValue({
      user: { id: 'test-user-id' },
    })

    // Restore default chaining
    mockedSupabase.from.mockReturnValue(mockedSupabase)
    mockedSupabase.select.mockReturnValue(mockedSupabase)
    mockedSupabase.order.mockReturnValue(mockedSupabase)
    mockedSupabase.in.mockReturnValue(mockedSupabase)
    mockedSupabase.eq.mockReturnValue(mockedSupabase)
    mockedSupabase.single.mockReturnValue(mockedSupabase)
    mockedSupabase.insert.mockReturnValue(mockedSupabase)
    mockedSupabase.limit.mockReturnValue(mockedSupabase)
    mockedSupabase.update.mockReturnValue(mockedSupabase)
    mockedSupabase.maybeSingle.mockReturnValue(mockedSupabase)
  })

  it('verify initial state', () => {
    const store = useSharedPlanStore()

    expect(store.sharedPlan).toBeNull()
    expect(store.sharedHistory).toEqual([])
    expect(store.isLoading).toBe(false)
    expect(store.error).toBeNull()
  })

  describe('fetchSharedPlanByHash', () => {
    it('fetches a shared plan successfully', async () => {
      const store = useSharedPlanStore()
      const mockSharedPlanData = { plan_id: 'test-plan-id', user_id: 'test-sharer-id' }
      const mockPlanData = {
        plan_id: 'test-plan-id',
        title: 'Test Plan',
        description: 'Desc',
        plan_table: [],
      }
      const mockProfileData = { username: 'Test User' }

      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'shared_plans') {
          return {
            select: vi.fn().mockReturnThis(),
            eq: vi.fn().mockReturnThis(),
            single: vi.fn().mockResolvedValue({ data: mockSharedPlanData, error: null }),
          }
        }
        if (tableName === 'plans') {
          return {
            select: vi.fn().mockReturnThis(),
            eq: vi.fn().mockReturnThis(),
            single: vi.fn().mockResolvedValue({ data: mockPlanData, error: null }),
          }
        }
        if (tableName === 'profiles') {
          return {
            select: vi.fn().mockReturnThis(),
            eq: vi.fn().mockReturnThis(),
            single: vi.fn().mockResolvedValue({ data: mockProfileData, error: null }),
          }
        }
        if (tableName === 'shared_history') {
          return {
            select: vi.fn().mockReturnThis(),
            eq: vi.fn().mockReturnThis(),
            single: vi.fn().mockResolvedValue({ data: null, error: null }),
            insert: vi.fn().mockReturnThis(),
          }
        }
        return {}
      })

      await store.fetchSharedPlanByHash('test-hash')

      expect(store.sharedPlan).toEqual({
        plan: {
          plan_id: mockPlanData.plan_id,
          title: mockPlanData.title,
          description: mockPlanData.description,
          table: mockPlanData.plan_table,
        },
        sharer_username: mockProfileData.username,
        sharer_id: mockSharedPlanData.user_id,
      })
      expect(store.error).toBeNull()
      expect(store.isLoading).toBe(false)
    })

    it('handles plan not found error', async () => {
      const store = useSharedPlanStore()
      mockedSupabase.from.mockReturnValue({
        select: vi.fn().mockReturnThis(),
        eq: vi.fn().mockReturnThis(),
        single: vi.fn().mockResolvedValue({ data: null, error: new Error('Not found') }),
      })

      await store.fetchSharedPlanByHash('invalid-hash')

      expect(store.sharedPlan).toBeNull()
      expect(store.error).not.toBeNull() // Should be a localized error string
      expect(store.isLoading).toBe(false)
    })
  })

  describe('fetchSharedHistory', () => {
    it('fetches shared history successfully', async () => {
      const store = useSharedPlanStore()
      const mockHistory = [
        { plan_id: 'plan-1', user_id: 'user-1', shared_by: 'sharer-1', created_at: '2023-01-01' },
      ]
      const mockPlans = [
        { plan_id: 'plan-1', title: 'Plan 1', description: 'Desc 1', plan_table: [] },
      ]

      mockedSupabase.from.mockImplementation((tableName: string) => {
        if (tableName === 'shared_history') {
          return {
            select: vi.fn().mockReturnThis(),
            eq: vi.fn().mockReturnThis(),
            order: vi.fn().mockReturnThis(),
            limit: vi.fn().mockResolvedValue({ data: mockHistory, error: null }),
          }
        }
        if (tableName === 'plans') {
          return {
            select: vi.fn().mockReturnThis(),
            in: vi.fn().mockResolvedValue({ data: mockPlans, error: null }),
          }
        }
        return {}
      })

      await store.fetchSharedHistory()

      expect(store.sharedHistory).toHaveLength(1)
      expect(store.sharedHistory[0].plan?.title).toBe('Plan 1')
      expect(store.isLoading).toBe(false)
    })
  })

  describe('upsertCurrentPlan', () => {
    it('upserts plan and returns new plan_id', async () => {
      const store = useSharedPlanStore()
      // Setup initial state
      store.sharedPlan = {
        plan: {
          plan_id: 'old-id',
          title: 'Test',
          description: 'Desc',
          table: [],
        },
        sharer_username: 'User',
        sharer_id: 'sharer-id',
      }
      store.isForked = false

      const mockResponse = {
        success: true,
        data: { plan_id: 'new-id' },
      }
      mockedApiUpsertPlan.mockResolvedValue(mockResponse)

      await store.upsertCurrentPlan()

      expect(store.sharedPlan?.plan.plan_id).toBe('new-id')
      expect(store.isForked).toBe(true)
    })
  })
})
