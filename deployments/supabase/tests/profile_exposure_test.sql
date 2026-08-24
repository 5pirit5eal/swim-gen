begin;

select plan(18);

create temporary table test_profile_exposure_context (
  owner_id uuid not null,
  other_id uuid not null
);

insert into test_profile_exposure_context
select
  owner_user.id,
  other_user.id
from auth.users owner_user
join auth.users other_user
  on other_user.email = 'feedback-test-swimmer@example.com'
where owner_user.email = 'css-test-swimmer@example.com';

select is(
  (select count(*)::integer from test_profile_exposure_context),
  1,
  'seed contains two users for profile exposure tests'
);

select owner_id, other_id
from test_profile_exposure_context

\gset profile_exposure_

select ok(
  not has_table_privilege('anon', 'public.profiles', 'SELECT'),
  'anonymous clients cannot read the base profiles table'
);

select ok(
  not has_table_privilege('authenticated', 'public.profiles', 'SELECT'),
  'authenticated clients do not have table-wide profile read privilege'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'user_id', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'updated_at', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'username', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'experience', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'preferred_language', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'preferred_strokes', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'categories', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'overall_generations', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'monthly_generations', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'exports', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'css_200m_seconds', 'SELECT')
    and has_column_privilege('authenticated', 'public.profiles', 'css_400m_seconds', 'SELECT'),
  'authenticated clients have explicit profile column read grants'
);

select ok(
  has_table_privilege('service_role', 'public.profiles', 'SELECT'),
  'service_role retains trusted profile read access'
);

select ok(
  has_table_privilege('service_role', 'public.profiles', 'UPDATE'),
  'service_role retains trusted profile update access'
);

select ok(
  has_table_privilege('anon', 'public.public_profiles', 'SELECT'),
  'anonymous clients can read the public profile projection'
);

select ok(
  has_table_privilege('authenticated', 'public.public_profiles', 'SELECT'),
  'authenticated clients can read the public profile projection'
);

select is(
  (select count(*)::integer
   from information_schema.columns
   where table_schema = 'public'
     and table_name = 'public_profiles'),
  1,
  'public profile projection contains one column'
);

select is(
  (select column_name
   from information_schema.columns
   where table_schema = 'public'
     and table_name = 'public_profiles'),
  'username',
  'public profile projection exposes only username'
);

select is(
  (select count(*)::integer
   from information_schema.columns
   where table_schema = 'public'
     and table_name = 'public_profiles'
     and column_name in ('user_id', 'experience', 'preferred_language', 'exports')),
  0,
  'public profile projection contains no stable identifiers or private fields'
);

set local role authenticated;
select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'profile_exposure_owner_id', 'role', 'authenticated')::text,
  true
);

select is(
  (select count(user_id)::integer
   from public.profiles
   where user_id = :'profile_exposure_owner_id'),
  1,
  'authenticated users can read their own profile'
);

select is(
  (select count(user_id)::integer
   from public.profiles
   where user_id = :'profile_exposure_other_id'),
  0,
  'authenticated users cannot read another profile'
);

set local role anon;
select set_config('request.jwt.claims', json_build_object('role', 'anon')::text, true);

select throws_ok(
  $$select username from public.profiles$$,
  '42501',
  null,
  'anonymous users cannot query the base profiles table'
);

select is(
  (select username from public.public_profiles where username = 'css_test_swimmer'),
  'css_test_swimmer',
  'anonymous users can read a username from the public projection'
);

select is(
  (select count(*)::integer
   from public.get_public_profile_username(:'profile_exposure_owner_id')),
  1,
  'anonymous users can use the narrow username lookup RPC'
);

select is(
  (select username
   from public.get_public_profile_username(:'profile_exposure_owner_id')),
  'css_test_swimmer',
  'username lookup RPC returns the expected username'
);

select is(
  (select count(*)::integer
   from public.get_public_profile_username(gen_random_uuid())),
  0,
  'username lookup RPC returns no row for an unknown user'
);

select * from finish();

rollback;
