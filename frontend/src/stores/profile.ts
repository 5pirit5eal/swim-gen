import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

import type { Profile, ProfileUpdate } from '@/types'
import { getSupabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'

const PROFILE_COLUMNS =
  'user_id, updated_at, username, experience, preferred_language, preferred_strokes, categories, overall_generations, monthly_generations, exports, css_200m_seconds, css_400m_seconds'
const PROFILE_UPDATE_FIELDS = [
  'username',
  'experience',
  'preferred_language',
  'preferred_strokes',
  'categories',
  'css_200m_seconds',
  'css_400m_seconds',
] as const

export const useProfileStore = defineStore('profile', () => {
  const loading = ref(false)
  const profile = ref<Profile | null>(null)
  const userStore = useAuthStore()
  const error = ref<string | null>(null)

  watch(
    () => userStore.user?.id ?? null,
    async (newUserId) => {
      if (newUserId) {
        await _fetchProfile()
      } else {
        profile.value = null
      }
    },
    { immediate: true },
  )

  async function _fetchProfile() {
    if (!userStore.user) {
      error.value = 'User is not available.'
      return
    }
    const supabase = await getSupabase()
    error.value = null
    const { data, error: query_error } = await supabase
      .from('profiles')
      .select(PROFILE_COLUMNS)
      .eq('user_id', userStore.user.id)
      .single()
    if (query_error) {
      console.error(query_error)
      error.value = query_error.message
    } else {
      profile.value = data
    }
  }

  async function fetchProfile() {
    loading.value = true
    await _fetchProfile()
    loading.value = false
  }

  async function updateProfile(updatedProfile: ProfileUpdate) {
    if (!userStore.user) {
      return
    }
    const supabase = await getSupabase()
    loading.value = true
    error.value = null
    const profileUpdate = Object.fromEntries(
      PROFILE_UPDATE_FIELDS.filter((field) =>
        Object.prototype.hasOwnProperty.call(updatedProfile, field),
      ).map((field) => [field, updatedProfile[field]]),
    ) as ProfileUpdate
    const { data, error: update_error } = await supabase
      .from('profiles')
      .update(profileUpdate)
      .eq('user_id', userStore.user.id)
      .select(PROFILE_COLUMNS)
      .single()
    if (update_error) {
      console.error(update_error)
      error.value = update_error.message
    } else {
      profile.value = data
    }
    loading.value = false
  }

  return {
    loading,
    profile,
    error,
    fetchProfile,
    updateProfile,
  }
})
