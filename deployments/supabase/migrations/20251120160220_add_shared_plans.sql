-- Create table 'shared_plans' to manage user-shared plans with reference to plans table
create table shared_plans (
  user_id uuid not null references auth.users on delete cascade primary key,
  plan_id uuid not null references plans on delete cascade,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  url_hash uuid not null default gen_random_uuid(),
  share_count int not null default 0
);
create unique index idx_shared_plans_id on shared_plans(plan_id);
create unique index idx_shared_plans_url_hash on shared_plans(url_hash);

-- Create RLS policy to allow all users to view shared plans and only owners to create or delete

alter table shared_plans enable row level security;
create policy "Anyone can view shared plans." on shared_plans
  for select using (true);
create policy "Users can insert their own shared plans." on shared_plans
  for insert with check ((select auth.uid()) = user_id);
create policy "Users can delete their own shared plans." on shared_plans
  for delete using ((select auth.uid()) = user_id);
-- Allow incrementing share_count for all users
create policy "Anyone can update share_count on shared plans." on shared_plans
  for update using (true)
  with check (true);

-- Keep updated_at fresh
create or replace function public.shared_plans_set_updated_at()
returns trigger
language plpgsql as $$
begin
    new.updated_at := now();
    return new;
end;
$$;

create trigger shared_plans_updated_at_trg
before update on public.shared_plans
for each row execute procedure public.shared_plans_set_updated_at();

-- Create table 'shared_history' to track to whom plans have been shared
create table shared_history (
  user_id uuid not null references auth.users on delete cascade,
  plan_id uuid not null references plans on delete cascade,
  share_method text not null default 'link', -- e.g., 'link', 'email'
  shared_by uuid not null references shared_plans(user_id) on delete cascade,
  created_at timestamptz not null default now(),
  primary key (user_id, plan_id)
);
create index idx_shared_id on shared_history (user_id);

-- Create RLS policy to allow users to only access their own shared history
alter table shared_history enable row level security;
create policy "Users can view their own shared history." on shared_history
  for select using ((select auth.uid()) = user_id);
create policy "Users can insert their own shared history." on shared_history
  for insert with check ((select auth.uid()) = user_id);
create policy "Users can delete their own shared history." on shared_history
  for delete using ((select auth.uid()) = user_id);
-- No updates allowed on shared_history
create policy "No updates allowed on shared history." on shared_history
  for update using (false)
  with check (false);


-- Increment share_count in shared_plans when a new entry is added to shared_history
create or replace function public.increment_shared_plan_count()
returns trigger
language plpgsql as $$
begin
    update public.shared_plans
    set share_count = share_count + 1
    where plan_id = new.plan_id;
    return new;
end;
$$;

create trigger increment_shared_plan_count_trg
after insert on public.shared_history
for each row execute procedure public.increment_shared_plan_count();

-- Update policies for plans to allow access if referenced in shared_plans
create policy "User can query plans in any shared plans" ON public.plans
FOR ALL
TO authenticated
USING (
    EXISTS (
        SELECT 1
        FROM public.shared_plans
        WHERE shared_plans.plan_id = plans.plan_id
    )
);

-- Update policies for plans to allow access if referenced in shared_history
create policy "User can query plans in his shared history" ON public.plans
FOR ALL
TO authenticated
USING (
  EXISTS (
    SELECT 1
    FROM public.shared_history
    WHERE shared_history.plan_id = plans.plan_id
      AND shared_history.user_id = (SELECT auth.uid())
  )
);

