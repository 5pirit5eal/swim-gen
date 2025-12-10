import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

import type { Profile } from '@/types'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'

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
    error.value = null
    const { data, error: query_error } = await supabase
      .from('profiles')
      .select('*')
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

  async function updateProfile(updatedProfile: Partial<Profile>) {
    if (!userStore.user) {
      return
    }
    loading.value = true
    error.value = null
    const { data, error: update_error } = await supabase
      .from('profiles')
      .update(updatedProfile)
      .eq('user_id', userStore.user.id)
      .select()
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
