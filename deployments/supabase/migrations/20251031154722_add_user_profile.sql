-- Create a table for public profiles
create table profiles (
  user_id uuid references auth.users on delete cascade not null primary key,
  updated_at timestamptz not null default now(),
  username text unique not null,
  experience text,
  preferred_language text,
  preferred_strokes text[] not null default '{}',
  categories text[] not null default '{}',
  constraint username_length check (char_length(username) >= 3)
);
-- Case-insensitive index on username
create unique index profiles_username_ci_idx on profiles (lower(username));

-- Set up Row Level Security (RLS)
-- See https://supabase.com/docs/guides/database/postgres/row-level-security for more details.
alter table profiles enable row level security;

create policy "Public profiles are viewable by everyone." on profiles
  for select using (true);

create policy "Users can insert their own profile." on profiles
  for insert with check ((select auth.uid()) = user_id);

create policy "Users can update own profile." on profiles
  for update using ((select auth.uid()) = user_id);

-- Keep updated_at fresh
create or replace function public.profiles_set_updated_at()
returns trigger
language plpgsql as $$
begin
  new.updated_at := now();
  return new;
end;
$$;

create trigger profiles_updated_at_trg
before update on public.profiles
for each row execute procedure public.profiles_set_updated_at();

-- This trigger automatically creates a profile entry when a new user signs up via Supabase Auth.
-- See https://supabase.com/docs/guides/auth/managing-user-data#using-triggers for more details.
create or replace function public.handle_new_user()
returns trigger
security definer
set search_path = public, auth
language plpgsql as $$
declare
  v_username text;
begin
  -- Try meta, fallback to generated value
  v_username := coalesce(new.raw_user_meta_data->>'username', 'user_' || substr(new.id::text, 1, 8));

  begin
    insert into public.profiles (user_id, username)
    values (new.id, v_username);
  exception
    when unique_violation then
      raise exception using
        message = 'Username already taken. Please choose another one.',
        errcode = '23505';
  end;
  return new;
end;
$$;

create trigger on_auth_user_created
    after insert on auth.users
    for each row execute procedure public.handle_new_user();
