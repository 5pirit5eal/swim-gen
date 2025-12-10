import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Session, User, EmailOtpType } from '@supabase/supabase-js'
import { supabase } from '@/plugins/supabase'

export const useAuthStore = defineStore('auth', () => {
  // --- STATE ---
  const session = ref<Session | null>(null)
  const user = ref<User | null>(null)
  const hasInitialized = ref(false)

  supabase.auth.onAuthStateChange((event, newSession) => {
    console.log(`Auth event: ${event}`)
    if (event === 'INITIAL_SESSION' && !hasInitialized.value) {
      hasInitialized.value = true
    }
    session.value = newSession
    user.value = newSession?.user ?? null
  })

  // --- COMPUTED ---

  // --- ACTIONS ---
  async function signInWithPassword(email: string, password: string) {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    if (error) throw error
    return data
  }

  async function signUp(email: string, password: string, username: string) {
    // Check if the username is already taken
    const { data: existingUser, error: existingUserError } = await supabase
      .from('profiles')
      .select('username')
      .eq('username', username)
      .single()

    if (existingUserError && existingUserError.code !== 'PGRST116') {
      throw existingUserError
    }

    if (existingUser) {
      throw new Error('Username already taken')
    }

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

  async function signInWithOAuth() {
    console.log('Signing in with OAuth and redirecting to', `${window.location.origin}/`)
    const { data, error } = await supabase.auth.signInWithOAuth({
      provider: 'google',
      options: {
        redirectTo: `${window.location.origin}/`,
      },
    })
    if (error) throw error
    return data
  }

  async function signOut() {
    const { error } = await supabase.auth.signOut()
    if (error) throw error
  }

  async function resetPassword(email: string, redirectTo: string) {
    const { error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo,
    })
    if (error) throw error
  }

  async function updatePassword(password: string) {
    const { error } = await supabase.auth.updateUser({ password })
    if (error) throw error
  }

  async function verifyOtp(token_hash: string, type: EmailOtpType) {
    const { error } = await supabase.auth.verifyOtp({ token_hash, type })
    if (error) throw error
  }

  return {
    session,
    user,
    hasInitialized,
    signInWithPassword,
    signUp,
    signInWithOAuth,
    signOut,
    resetPassword,
    updatePassword,
    verifyOtp,
  }
})
