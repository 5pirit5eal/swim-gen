-- Create table 'history' to manage user-created plans with reference to plans table
create table history (
  user_id uuid not null references auth.users on delete cascade,
  plan_id uuid not null references plans on delete cascade,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  keep_forever boolean not null default false,
  primary key (user_id, plan_id)
);
create index idx_history_plan_id on history(plan_id);


-- Create RLS policy to allow users to only access their own plan history
alter table history enable row level security;

create policy "Users can view their own plan history." on history
  for select using ((select auth.uid()) = user_id);

create policy "Users can insert their own plan history." on history
  for insert with check ((select auth.uid()) = user_id);

create policy "Users can update their own plan history." on history
  for update using ((select auth.uid()) = user_id);

create policy "Users can delete their own plan history." on history
  for delete using ((select auth.uid()) = user_id);

-- Update policies for plans to allow access if referenced in history
drop policy if exists plans_access_for_donators on public.plans;
create policy "User can query Plans in his History" ON public.plans
FOR ALL
TO authenticated
USING (
  EXISTS (
    SELECT 1
    FROM public.history
    WHERE history.plan_id = plans.plan_id
      AND history.user_id = (SELECT auth.uid())
  )
);

-- Update policies for donations to allow access if referenced in history
create policy "User can query Donations for Plans" ON public.plans
FOR ALL
TO authenticated
USING (
  EXISTS (
    SELECT 1
    FROM public.donations
    WHERE donations.plan_id = plans.plan_id
      AND donations.user_id = (SELECT auth.uid())
  )
);

-- Keep updated_at fresh
create or replace function public.history_set_updated_at()
returns trigger
language plpgsql as $$
begin
  new.updated_at := now();
  return new;
end;
$$;

create trigger history_updated_at_trg
before update on public.history
for each row execute procedure public.history_set_updated_at();


-- Alter the 'profiles' table to add a overall plan generation, monthly generation and export tracker column
alter table profiles
add column overall_generations int not null default 0,
add column monthly_generations int not null default 0,
add column exports int not null default 0;

-- Alter the 'plans' table to add a updated_at column
alter table plans
add column updated_at timestamptz not null default now(),
add column exports int not null default 0;


-- Create a pg_cron job that cleans up old plan history entries and the respective plan data
select cron.schedule(
    'clean-up-old-plans',
    '0 0 * * *', -- Every day at midnight
    $$
    begin
        -- Delete history entries older than 30 days
        delete from public.history where created_at < now() - interval '30 days' and keep_forever = false and export_count = 0;

        -- Delete plans that are not referenced anymore
        delete from public.plans p
        where not exists (select 1 from public.history h where h.plan_id = p.plan_id)
          and not exists (select 1 from public.donations d where d.plan_id = p.plan_id)
          and not exists (select 1 from public.feedback f where f.plan_id = p.plan_id)
          and not exists (select 1 from public.scraped s where s.plan_id = p.plan_id);
    end;
    $$
);

-- Create a pg_cron job that resets all monthly generations at the beginning of the month
select cron.schedule(
    'reset-monthly-generations',
    '0 0 1 * *', -- Every month on the 1st at midnight
    $$
    begin
        update public.profiles
        set monthly_generations = 0;
    end;
    $$
);

-- Create a function that increases the profile plan generation tracker upon new plan creation
create or replace function public.increment_plan_generations()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
begin
    update profiles
    set monthly_generations = monthly_generations + 1,
        overall_generations = overall_generations + 1,
        updated_at = now()
    where user_id = new.user_id;
    return new;
end;
$$;

create trigger on_plan_history_created
    after insert on public.history
    for each row execute procedure public.increment_plan_generations();
