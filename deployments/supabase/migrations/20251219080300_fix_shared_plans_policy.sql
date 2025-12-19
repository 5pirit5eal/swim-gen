-- Drop the existing policy and recreate it with 'public' access to allow anonymous users to view shared plans
drop policy if exists "User can query plans in any shared plans" on public.plans;

create policy "Anyone can query plans in any shared plans" on public.plans
for select
to public
using (
    exists (
        select 1
        from public.shared_plans
        where shared_plans.plan_id = plans.plan_id
    )
);
