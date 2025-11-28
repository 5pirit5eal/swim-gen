-- Create table 'memory' to store chat history
create table if not exists memory (
  id uuid primary key default gen_random_uuid(),
  plan_id uuid not null references plans(plan_id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  role text not null check (role in ('user', 'ai')),
  content text not null,
  previous_message_id uuid references memory(id),
  next_message_id uuid references memory(id),
  plan_snapshot jsonb,
  created_at timestamptz not null default now()
);

-- Indexes for performance
create index idx_memory_plan_id on memory(plan_id);
create index idx_memory_user_id on memory(user_id);

-- RLS Policies
alter table memory enable row level security;

create policy "Users can view their own memory." on memory
  for select using ((select auth.uid()) = user_id);

create policy "Users can insert their own memory." on memory
  for insert with check ((select auth.uid()) = user_id);

create policy "Users can update their own memory." on memory
  for update using ((select auth.uid()) = user_id);

create policy "Users can delete their own memory." on memory
  for delete using ((select auth.uid()) = user_id);
