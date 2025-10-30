import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Session, User } from '@supabase/supabase-js'
import { supabase } from '@/plugins/supabase'

export const useAuthStore = defineStore('auth', () => {
  const session = ref<Session | null>(null)
  const user = ref<User | null>(null)

  supabase.auth.getSession().then(({ data }) => {
    session.value = data.session
  })

  supabase.auth.getUser().then(({ data }) => {
    user.value = data.user ?? null
  })

  supabase.auth.onAuthStateChange((event, newSession) => {
    session.value = newSession
    user.value = newSession?.user ?? null
  })

  async function signInWithPassword(email: string, password: string) {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    if (error) throw error
    return data
  }

  async function signUp(email: string, password: string, username: string) {
    // TODO: Check if the username is already taken, can be done once the profile table is setup

    const { data, error } = await supabase.auth.signUp({
      email,
      password,
      options: {
        data: {
          username,
        },
      },
    })
    if (error) throw error
    return data
  }

  async function signOut() {
    const { error } = await supabase.auth.signOut()
    if (error) throw error
  }

  return {
    session,
    user,
    signInWithPassword,
    signUp,
    signOut,
  }
})
