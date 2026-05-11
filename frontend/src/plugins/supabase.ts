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

/**
 * @deprecated Do not use this export directly. Always call `await getSupabase()` instead.
 * This Proxy exists solely for backwards-compatibility with test mocks and will throw
 * a descriptive error at runtime if any code accidentally accesses it outside of tests.
 */
export const supabase = new Proxy({} as import('@supabase/supabase-js').SupabaseClient, {
  get(_target, prop) {
    throw new Error(
      `[supabase] Attempted to access .${String(prop)} on the legacy 'supabase' export. ` +
        `Use 'await getSupabase()' instead.`,
    )
  },
})
