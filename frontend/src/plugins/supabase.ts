const supabaseUrl = import.meta.env.VITE_SUPABASE_URL
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY

let supabaseClientPromise: Promise<import('@supabase/supabase-js').SupabaseClient> | null = null

export async function getSupabase() {
  if (!supabaseClientPromise) {
    supabaseClientPromise = import('@supabase/supabase-js').then(({ createClient }) =>
      createClient(supabaseUrl, supabaseAnonKey),
    )
  }

  return supabaseClientPromise
}

export const supabase = {} as import('@supabase/supabase-js').SupabaseClient
