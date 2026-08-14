-- Feedback writes are performed by the trusted backend after checking plan access.
drop policy if exists feedback_insert_own on public.feedback;
drop policy if exists feedback_update_own on public.feedback;
drop policy if exists feedback_delete_own on public.feedback;

revoke all on public.feedback from anon;
revoke insert, update, delete on public.feedback from authenticated;
grant select on public.feedback to authenticated;

-- Authorize and upsert in one trusted operation. The shared-history branch
-- requires an active share row and derives its owner from that relationship.
create or replace function public.submit_feedback(
  p_user_id uuid,
  p_plan_id uuid,
  p_rating integer,
  p_was_swam boolean,
  p_difficulty_rating integer,
  p_comment text
)
returns boolean
language sql
security definer
set search_path = ''
as $$
  with upserted as (
    insert into public.feedback (
      user_id,
      plan_id,
      rating,
      was_swam,
      difficulty_rating,
      comment,
      removed_from_history
    )
    select
      p_user_id,
      p_plan_id,
      p_rating,
      p_was_swam,
      p_difficulty_rating,
      p_comment,
      false
    where exists (
      select 1
      from public.history h
      where h.user_id = p_user_id
        and h.plan_id = p_plan_id
    )
    or exists (
      select 1
      from public.donations d
      where d.user_id = p_user_id
        and d.plan_id = p_plan_id
    )
    or exists (
      select 1
      from public.shared_history sh
      join public.shared_plans sp
        on sp.plan_id = sh.plan_id
       and sp.user_id = sh.shared_by
      where sh.user_id = p_user_id
        and sh.plan_id = p_plan_id
    )
    on conflict (user_id, plan_id) do update
      set rating = excluded.rating,
          was_swam = excluded.was_swam,
          difficulty_rating = excluded.difficulty_rating,
          comment = excluded.comment,
          removed_from_history = false
    returning 1
  )
  select exists (select 1 from upserted);
$$;

revoke execute on function public.submit_feedback(uuid, uuid, integer, boolean, integer, text)
  from public, anon, authenticated, service_role;
grant execute on function public.submit_feedback(uuid, uuid, integer, boolean, integer, text)
  to postgres;
