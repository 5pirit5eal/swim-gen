-- Keep the profile table private while exposing only the username needed by
-- public registration and shared-plan flows.
drop policy if exists "Public profiles are viewable by everyone." on public.profiles;
drop policy if exists "Users can read their own profile." on public.profiles;

revoke select on public.profiles from public, anon, authenticated;

create policy "Users can read their own profile."
on public.profiles
for select
to authenticated
using ((select auth.uid()) = user_id);

grant select (
  user_id,
  updated_at,
  username,
  experience,
  preferred_language,
  preferred_strokes,
  categories,
  overall_generations,
  monthly_generations,
  exports,
  css_200m_seconds,
  css_400m_seconds
) on public.profiles to authenticated;

create or replace view public.public_profiles as
select username
from public.profiles;

revoke all on public.public_profiles from public, anon, authenticated;
grant select on public.public_profiles to anon, authenticated;

-- Resolve one username without exposing any other profile columns.
create or replace function public.get_public_profile_username(p_user_id uuid)
returns table (username text)
language sql
security definer
set search_path = ''
as $$
  select profile.username
  from public.profiles as profile
  where profile.user_id = p_user_id;
$$;

revoke all on function public.get_public_profile_username(uuid)
  from public, anon, authenticated, service_role;
grant execute on function public.get_public_profile_username(uuid)
  to anon, authenticated, service_role;

-- Keep the existing shared-plan response columns and order while resolving the
-- sharer username through the narrow lookup above.
create or replace function public.get_shared_plan_by_hash(p_url_hash uuid)
returns table (
  plan_id uuid,
  sharer_id uuid,
  sharer_username text,
  title text,
  description text,
  plan_table jsonb
)
language sql
security definer
set search_path = ''
as $$
  select
    sp.plan_id,
    sp.user_id,
    profile.username,
    plan.title,
    plan.description,
    plan.plan_table
  from public.shared_plans as sp
  join public.plans as plan on plan.plan_id = sp.plan_id
  left join lateral public.get_public_profile_username(sp.user_id) as profile on true
  where sp.url_hash = p_url_hash;
$$;

revoke execute on function public.get_shared_plan_by_hash(uuid)
  from public, anon, authenticated;
grant execute on function public.get_shared_plan_by_hash(uuid)
  to anon, authenticated;
