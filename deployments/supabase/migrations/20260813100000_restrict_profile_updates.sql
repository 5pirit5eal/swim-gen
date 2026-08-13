-- Profile creation and server-managed fields must not be writable by browser roles.
revoke insert, update on public.profiles from public, anon, authenticated;

-- The timestamp trigger must retain trusted write access to its server-managed column.
create or replace function public.profiles_set_updated_at()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
begin
  new.updated_at := now();
  return new;
end;
$$;

grant update (
  username,
  experience,
  preferred_language,
  preferred_strokes,
  categories,
  css_200m_seconds,
  css_400m_seconds
) on public.profiles to authenticated;

revoke execute on function public.profiles_set_updated_at() from public, anon, authenticated;
