begin;

select plan(28);

create temporary table test_profile_context (
  owner_id uuid not null,
  other_id uuid not null
);

insert into test_profile_context
select
  owner_user.id,
  other_user.id
from auth.users owner_user
join auth.users other_user
  on other_user.email = 'feedback-test-swimmer@example.com'
where owner_user.email = 'css-test-swimmer@example.com';

select is(
  (select count(*)::integer from test_profile_context),
  1,
  'seed contains two users for profile authorization tests'
);

select owner_id, other_id
from test_profile_context

\gset profile_

select ok(
  not has_table_privilege('anon', 'public.profiles', 'INSERT'),
  'anonymous clients cannot insert profiles'
);

select ok(
  not has_table_privilege('anon', 'public.profiles', 'UPDATE'),
  'anonymous clients cannot update profiles'
);

select ok(
  not has_table_privilege('authenticated', 'public.profiles', 'INSERT'),
  'authenticated clients cannot insert profiles'
);

select ok(
  not has_table_privilege('authenticated', 'public.profiles', 'UPDATE'),
  'authenticated clients do not have table-wide profile update privilege'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'username', 'UPDATE'),
  'authenticated clients can update username'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'experience', 'UPDATE'),
  'authenticated clients can update experience'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'preferred_language', 'UPDATE'),
  'authenticated clients can update preferred_language'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'preferred_strokes', 'UPDATE'),
  'authenticated clients can update preferred_strokes'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'categories', 'UPDATE'),
  'authenticated clients can update categories'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'css_200m_seconds', 'UPDATE'),
  'authenticated clients can update css_200m_seconds'
);

select ok(
  has_column_privilege('authenticated', 'public.profiles', 'css_400m_seconds', 'UPDATE'),
  'authenticated clients can update css_400m_seconds'
);

select ok(
  not has_column_privilege('authenticated', 'public.profiles', 'user_id', 'UPDATE'),
  'authenticated clients cannot update user_id'
);

select ok(
  not has_column_privilege('authenticated', 'public.profiles', 'updated_at', 'UPDATE'),
  'authenticated clients cannot update updated_at'
);

select ok(
  not has_column_privilege('authenticated', 'public.profiles', 'overall_generations', 'UPDATE'),
  'authenticated clients cannot update overall_generations'
);

select ok(
  not has_column_privilege('authenticated', 'public.profiles', 'monthly_generations', 'UPDATE'),
  'authenticated clients cannot update monthly_generations'
);

select ok(
  not has_column_privilege('authenticated', 'public.profiles', 'exports', 'UPDATE'),
  'authenticated clients cannot update exports'
);

select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'profile_owner_id', 'role', 'authenticated')::text,
  true
);
set local role authenticated;

select lives_ok(
  format(
    $sql$
      update public.profiles
      set username = username,
          experience = 'profile-column-test',
          preferred_language = 'de',
          preferred_strokes = array['freestyle'],
          categories = array['endurance'],
          css_200m_seconds = 180,
          css_400m_seconds = 390
      where user_id = %L
    $sql$,
    :'profile_owner_id'
  ),
  'authenticated users can update all allowlisted profile columns'
);

select is(
  (
    select count(*)::integer
    from public.profiles
    where user_id = :'profile_owner_id'
      and experience = 'profile-column-test'
      and preferred_language = 'de'
      and preferred_strokes = array['freestyle']
      and categories = array['endurance']
      and css_200m_seconds = 180
      and css_400m_seconds = 390
  ),
  1,
  'allowlisted profile values are persisted for the owner'
);

select lives_ok(
  format(
    $sql$
      do $do$
      declare
        affected_rows integer;
      begin
        update public.profiles
        set experience = 'forged'
        where user_id = %L;
        get diagnostics affected_rows = row_count;
        if affected_rows <> 0 then
          raise exception 'cross-user profile update affected %% rows', affected_rows;
        end if;
      end
      $do$;
    $sql$,
    :'profile_other_id'
  ),
  'authenticated users cannot update another profile'
);

select throws_ok(
  format(
    'insert into public.profiles (user_id, username) values (%L, %L)',
    :'profile_owner_id',
    'forged_profile'
  ),
  '42501',
  null,
  'authenticated clients cannot insert profiles directly'
);

select throws_ok(
  format(
    'update public.profiles set updated_at = now() - interval %L where user_id = %L',
    '1 year',
    :'profile_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot forge updated_at'
);

select throws_ok(
  format(
    'update public.profiles set overall_generations = 999999 where user_id = %L',
    :'profile_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot forge overall_generations'
);

select throws_ok(
  format(
    'update public.profiles set monthly_generations = 999999 where user_id = %L',
    :'profile_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot forge monthly_generations'
);

select throws_ok(
  format(
    'update public.profiles set exports = 999999 where user_id = %L',
    :'profile_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot forge exports'
);

select throws_ok(
  format(
    'update public.profiles set user_id = %L where user_id = %L',
    :'profile_other_id',
    :'profile_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot change profile ownership'
);

set local role postgres;

select lives_ok(
  format(
    'update public.profiles set exports = exports + 1 where user_id = %L',
    :'profile_owner_id'
  ),
  'trusted database operations can update server-managed profile counters'
);

select ok(
  has_table_privilege('service_role', 'public.profiles', 'UPDATE'),
  'service_role retains trusted profile update access'
);

select * from finish();

rollback;
