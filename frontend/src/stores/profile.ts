import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

import type { Profile } from '@/types'
import { supabase } from '@/plugins/supabase'
import { useAuthStore } from '@/stores/auth'

export const useProfileStore = defineStore('profile', () => {
  const loading = ref(false)
  const profile = ref<Profile | null>(null)
  const userStore = useAuthStore()

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
      console.log('User is not available.')
      return
    }
    const { data, error } = await supabase
      .from('profiles')
      .select('*')
      .eq('user_id', userStore.user.id)
      .single()
    if (error) {
      console.error(error)
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
    const { data, error } = await supabase
      .from('profiles')
      .update(updatedProfile)
      .eq('user_id', userStore.user.id)
      .select()
      .single()
    if (error) {
      console.error(error)
    } else {
      profile.value = data
    }
    loading.value = false
  }

  return {
    loading,
    profile,
    fetchProfile,
    updateProfile,
  }
})
