-- Chat memory is managed by the authenticated backend, not by browser clients.
-- Keep the existing RLS policies as defense in depth for any trusted API path,
-- but remove direct table access from the Supabase API roles.
revoke all on public.memory from anon, authenticated;
