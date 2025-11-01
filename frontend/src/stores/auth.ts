import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Session, User } from '@supabase/supabase-js'
import { supabase } from '@/plugins/supabase'

export const useAuthStore = defineStore('auth', () => {
  const session = ref<Session | null>(null)
  const user = ref<User | null>(null)

  supabase.auth.getSession().then(({ data }) => {
    console.log('Fetched session on store init:', data.session)
    session.value = data.session
  })

  supabase.auth.getUser().then(({ data }) => {
    console.log('Fetched user on store init:', data.user)
    user.value = data.user ?? null
  })

  supabase.auth.onAuthStateChange((event, newSession) => {
    console.log('Auth state changed:', event)
    session.value = newSession
    user.value = newSession?.user ?? null
  })

  async function getSession() {
    if (session.value) return
    const { data, error } = await supabase.auth.getSession()
    if (error) {
      console.error('Error fetching session:', error)
      return
    }
    session.value = data.session
    console.log('Session fetched:', session.value)
  }

  async function getUser() {
    if (user.value) return
    if (!session.value) await getSession()
    const { data, error } = await supabase.auth.getUser()
    if (error) {
      console.error('Error fetching user:', error)
      return
    }
    user.value = data.user ?? null
    console.log('User fetched:', user.value)
  }

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
    getSession,
    getUser,
    signInWithPassword,
    signUp,
    signOut,
  }
})
