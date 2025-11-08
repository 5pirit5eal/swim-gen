-- Create table 'history' to manage user-created plans with reference to plans table
create table history (
  user_id uuid not null references auth.users on delete cascade,
  plan_id uuid not null references plans on delete cascade,
  created_at timestamptz not null default now(),
  export_count int not null default 0,
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


-- Alter the 'profiles' table to add a overall plan generation, monthly generation and export tracker column
alter table profiles
add column overall_generations int not null default 0,
add column monthly_generations int not null default 0,
add column exports int not null default 0;


-- Create a pg_cron job that cleans up old plan history entries and the respective plan data
select cron.schedule(
    'clean-up-old-plans',
    '0 0 * * *', -- Every day at midnight
    $$
    begin
        -- Delete history entries older than 90 days
        delete from public.history where created_at < now() - interval '90 days';

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
