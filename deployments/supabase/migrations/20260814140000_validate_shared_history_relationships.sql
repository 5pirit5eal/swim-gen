-- Validate and enforce relationship integrity between shared_history and shared_plans.

-- 1. Remove any stale, forged, or self-shared history rows.
delete from public.shared_history sh
where sh.user_id = sh.shared_by
   or not exists (
     select 1
     from public.shared_plans sp
     where sp.plan_id = sh.plan_id
       and sp.user_id = sh.shared_by
   );

-- 2. Add unique constraint on shared_plans (plan_id, user_id) to allow composite foreign key reference.
alter table public.shared_plans
  add constraint shared_plans_plan_id_user_id_key unique (plan_id, user_id);

-- 3. Replace loose foreign keys on shared_history with a strict composite foreign key to shared_plans.
alter table public.shared_history
  drop constraint if exists shared_history_plan_id_fkey,
  drop constraint if exists shared_history_shared_by_fkey;

alter table public.shared_history
  add constraint shared_history_plan_id_shared_by_fkey
  foreign key (plan_id, shared_by)
  references public.shared_plans (plan_id, user_id)
  on delete cascade;

-- 4. Enforce that users cannot share to themselves and that share_method is valid.
alter table public.shared_history
  add constraint shared_history_no_self_share check (user_id <> shared_by),
  add constraint shared_history_share_method_check check (share_method in ('link', 'email'));
