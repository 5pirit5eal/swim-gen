-- The pre-request hook executes under the PostgREST request role. Keep the
-- function security-definer, but allow API roles to invoke it.
revoke execute on function public.check_request() from public;
grant execute on function public.check_request() to anon, authenticated, service_role;
