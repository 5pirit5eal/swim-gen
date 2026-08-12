-- Restrict shared-plan access to validated RPCs and read-only recipient access.

-- Remove the policies that exposed share rows or granted broad plan access.
drop policy if exists "Anyone can view shared plans." on public.shared_plans;
drop policy if exists "Users can insert their own shared plans." on public.shared_plans;
drop policy if exists "Users can delete their own shared plans." on public.shared_plans;
drop policy if exists "Anyone can update share_count on shared plans." on public.shared_plans;

drop policy if exists "User can query Plans in his History" on public.plans;
drop policy if exists "User can query Donations for Plans" on public.plans;
drop policy if exists "User can query plans in any shared plans" on public.plans;
drop policy if exists "Anyone can query plans in any shared plans" on public.plans;
drop policy if exists "User can query plans in his shared history" on public.plans;

-- Remove direct client writes that could manufacture ownership relationships.
drop policy if exists "Users can insert their own plan history." on public.history;
drop policy if exists "Users can delete their own plan history." on public.history;
drop policy if exists donations_insert_own on public.donations;
drop policy if exists donations_delete_own on public.donations;

-- Keep only metadata updates for a user's own history row.
drop policy if exists "Users can update their own plan history." on public.history;
create policy "Users can update their own plan history."
on public.history
for update
to authenticated
using ((select auth.uid()) = user_id)
with check ((select auth.uid()) = user_id);

-- Plans are readable through validated ownership or recipient relationships only.
create policy "Users can read plans in their history"
on public.plans
for select
to authenticated
using (
  exists (
    select 1
    from public.history
    where history.plan_id = plans.plan_id
      and history.user_id = (select auth.uid())
  )
);

create policy "Users can read donated plans"
on public.plans
for select
to authenticated
using (
  exists (
    select 1
    from public.donations
    where donations.plan_id = plans.plan_id
      and donations.user_id = (select auth.uid())
  )
);

create policy "Recipients can read shared plans"
on public.plans
for select
to authenticated
using (
  exists (
    select 1
    from public.shared_history
    where shared_history.plan_id = plans.plan_id
      and shared_history.user_id = (select auth.uid())
  )
);

-- Shared history remains readable and removable only by the recipient. New rows
-- are created by record_shared_plan(), which derives all relationship fields.
drop policy if exists "Users can insert their own shared history." on public.shared_history;
drop policy if exists "Users can delete their own shared history." on public.shared_history;
create policy "Users can delete their own shared history."
on public.shared_history
for delete
to authenticated
using ((select auth.uid()) = user_id);

-- Remove stale or forged relationships left by the old client-writable policies.
delete from public.shared_history sh
where not exists (
  select 1
  from public.shared_plans sp
  where sp.plan_id = sh.plan_id
    and sp.user_id = sh.shared_by
    and (
      exists (
        select 1
        from public.history h
        where h.plan_id = sp.plan_id
          and h.user_id = sp.user_id
      )
      or exists (
        select 1
        from public.donations d
        where d.plan_id = sp.plan_id
          and d.user_id = sp.user_id
      )
    )
);

delete from public.shared_plans sp
where not exists (
  select 1
  from public.history h
  where h.plan_id = sp.plan_id
    and h.user_id = sp.user_id
)
and not exists (
  select 1
  from public.donations d
  where d.plan_id = sp.plan_id
    and d.user_id = sp.user_id
);

-- Keep share_count equal to the number of current recipient relationships.
update public.shared_plans sp
set share_count = (
  select count(*)::integer
  from public.shared_history sh
  where sh.plan_id = sp.plan_id
);

-- Trigger functions are trusted implementation details, not public RPCs.
create or replace function public.history_set_updated_at()
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

create or replace function public.shared_plans_set_updated_at()
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

create or replace function public.increment_plan_generations()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
begin
  update public.profiles
  set monthly_generations = monthly_generations + 1,
      overall_generations = overall_generations + 1,
      updated_at = now()
  where user_id = new.user_id;
  return new;
end;
$$;

create or replace function public.increment_shared_plan_count()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
declare
  affected_plan_id uuid;
begin
  if tg_op = 'DELETE' then
    affected_plan_id := old.plan_id;
  else
    affected_plan_id := new.plan_id;
  end if;

  update public.shared_plans
  set share_count = (
    select count(*)::integer
    from public.shared_history
    where plan_id = affected_plan_id
  )
  where plan_id = affected_plan_id;

  if tg_op = 'DELETE' then
    return old;
  end if;
  return new;
end;
$$;

drop trigger if exists increment_shared_plan_count_trg on public.shared_history;
create trigger increment_shared_plan_count_trg
after insert or delete on public.shared_history
for each row execute procedure public.increment_shared_plan_count();

-- Resolve a share hash without exposing shared_plans or plans as enumerable tables.
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
  from public.shared_plans sp
  join public.plans plan on plan.plan_id = sp.plan_id
  left join public.profiles profile on profile.user_id = sp.user_id
  where sp.url_hash = p_url_hash;
$$;

-- Record a recipient relationship from the bearer hash. The owner and plan are
-- always derived from the canonical share row, never accepted from the client.
create or replace function public.record_shared_plan(
  p_url_hash uuid,
  p_share_method text default 'link'
)
returns table (
  plan_id uuid,
  shared_by uuid,
  share_method text,
  created_at timestamptz
)
language plpgsql
security definer
set search_path = ''
as $$
declare
  recipient_id uuid := (select auth.uid());
  shared_plan_id uuid;
  owner_id uuid;
begin
  if recipient_id is null or p_share_method is distinct from 'link' then
    return;
  end if;

  select sp.plan_id, sp.user_id
  into shared_plan_id, owner_id
  from public.shared_plans sp
  where sp.url_hash = p_url_hash;

  if not found or recipient_id = owner_id then
    return;
  end if;

  insert into public.shared_history (user_id, plan_id, share_method, shared_by)
  values (recipient_id, shared_plan_id, p_share_method, owner_id)
  on conflict on constraint shared_history_pkey do nothing;

  return query
  select sh.plan_id, sh.shared_by, sh.share_method, sh.created_at
  from public.shared_history sh
  where sh.user_id = recipient_id
    and sh.plan_id = shared_plan_id;
end;
$$;

-- Existing default privileges granted broad access to every public object. Keep
-- future objects deny-by-default and expose only the two intended RPCs.
alter default privileges for role postgres in schema public
  revoke all on tables from anon, authenticated;
alter default privileges for role postgres in schema public
  revoke all on sequences from anon, authenticated;
alter default privileges for role postgres in schema public
  revoke all on functions from anon, authenticated;

revoke all on public.shared_plans from anon, authenticated;
revoke all on public.plans from anon, authenticated;
revoke all on public.history from anon;
revoke insert, update, delete on public.history from authenticated;
revoke insert, update, delete on public.donations from anon, authenticated;
revoke all on public.shared_history from anon;
revoke insert, update on public.shared_history from authenticated;

grant select on public.plans to authenticated;
grant select, update (keep_forever) on public.history to authenticated;
grant select on public.donations to authenticated;
grant select, delete on public.shared_history to authenticated;

revoke execute on function public.history_set_updated_at() from public, anon, authenticated;
revoke execute on function public.shared_plans_set_updated_at() from public, anon, authenticated;
revoke execute on function public.increment_plan_generations() from public, anon, authenticated;
revoke execute on function public.increment_shared_plan_count() from public, anon, authenticated;

revoke execute on function public.get_shared_plan_by_hash(uuid) from public, anon, authenticated;
grant execute on function public.get_shared_plan_by_hash(uuid) to anon, authenticated;

revoke execute on function public.record_shared_plan(uuid, text) from public, anon, authenticated;
grant execute on function public.record_shared_plan(uuid, text) to authenticated;

-- Harden the existing pre-request definer while retaining its authenticator call.
create or replace function public.check_request()
returns void
language plpgsql
security definer
set search_path = ''
as $$
declare
  req_method text := current_setting('request.method', true);
  req_ip inet := split_part(
    current_setting('request.headers', true)::json->>'x-forwarded-for',
    ',', 1
  )::inet;
  count_in_five_mins integer;
begin
  if req_method = 'GET' or req_method = 'HEAD' or req_method is null then
    return;
  end if;

  select count(*)
  into count_in_five_mins
  from private.rate_limits
  where ip = req_ip
    and request_at between now() - interval '5 minutes' and now();

  if count_in_five_mins > 100 then
    raise sqlstate 'PGRST' using
      message = json_build_object(
        'message', 'Rate limit exceeded, try again after a while'
      )::text,
      detail = json_build_object(
        'status', 420,
        'status_text', 'Enhance Your Calm'
      )::text;
  end if;

  insert into private.rate_limits (ip, request_at)
  values (req_ip, now());
end;
$$;

revoke execute on function public.check_request() from public, anon, authenticated;
grant execute on function public.check_request() to authenticator;

-- Do not let the cleanup job delete plans that still have an active share.
select cron.unschedule('clean-up-old-plans');

select cron.schedule(
  'clean-up-old-plans',
  '0 0 * * *',
  $$
  DO $do$
  BEGIN
    DELETE FROM public.history
    WHERE created_at < now() - interval '2 days'
      AND keep_forever = false;

    DELETE FROM public.plans p
    WHERE NOT EXISTS (SELECT 1 FROM public.history h WHERE h.plan_id = p.plan_id)
      AND NOT EXISTS (SELECT 1 FROM public.donations d WHERE d.plan_id = p.plan_id)
      AND NOT EXISTS (SELECT 1 FROM public.feedback f WHERE f.plan_id = p.plan_id)
      AND NOT EXISTS (SELECT 1 FROM public.scraped s WHERE s.plan_id = p.plan_id)
      AND NOT EXISTS (SELECT 1 FROM public.shared_plans sp WHERE sp.plan_id = p.plan_id);
  END
  $do$;
  $$
);
