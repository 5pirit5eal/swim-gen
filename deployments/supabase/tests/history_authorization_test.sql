begin;

select plan(27);

create temporary table test_history_context (
  owner_id uuid not null,
  other_id uuid not null,
  owner_plan_id uuid not null,
  other_plan_id uuid not null
);

insert into test_history_context
select
  owner_user.id,
  other_user.id,
  gen_random_uuid(),
  gen_random_uuid()
from auth.users owner_user
join auth.users other_user
  on other_user.email = 'feedback-test-swimmer@example.com'
where owner_user.email = 'css-test-swimmer@example.com';

select is(
  (select count(*)::integer from test_history_context),
  1,
  'seed contains two users for history authorization tests'
);

select owner_id, other_id, owner_plan_id, other_plan_id
from test_history_context

\gset history_

set local role postgres;

insert into public.plans (plan_id, title, description, plan_table)
values
  (:'history_owner_plan_id', 'Owner Plan', 'Owner Desc', '[]'::jsonb),
  (:'history_other_plan_id', 'Other Plan', 'Other Desc', '[]'::jsonb);

insert into public.history (user_id, plan_id, keep_forever)
values
  (:'history_owner_id', :'history_owner_plan_id', false),
  (:'history_other_id', :'history_other_plan_id', false);

-- Privilege checks for anonymous role
select ok(
  not has_table_privilege('anon', 'public.history', 'SELECT'),
  'anonymous clients cannot select history'
);

select ok(
  not has_table_privilege('anon', 'public.history', 'INSERT'),
  'anonymous clients cannot insert history'
);

select ok(
  not has_table_privilege('anon', 'public.history', 'UPDATE'),
  'anonymous clients cannot update history'
);

select ok(
  not has_table_privilege('anon', 'public.history', 'DELETE'),
  'anonymous clients cannot delete history'
);

-- Table-level privilege checks for authenticated role
select ok(
  not has_table_privilege('authenticated', 'public.history', 'INSERT'),
  'authenticated clients cannot insert history'
);

select ok(
  not has_table_privilege('authenticated', 'public.history', 'DELETE'),
  'authenticated clients cannot delete history'
);

select ok(
  not has_table_privilege('authenticated', 'public.history', 'UPDATE'),
  'authenticated clients do not have table-wide history update privilege'
);

select ok(
  has_table_privilege('authenticated', 'public.history', 'SELECT'),
  'authenticated clients can select history'
);

-- Column-level privilege checks for authenticated role
select ok(
  has_column_privilege('authenticated', 'public.history', 'keep_forever', 'UPDATE'),
  'authenticated clients can update keep_forever'
);

select ok(
  not has_column_privilege('authenticated', 'public.history', 'user_id', 'UPDATE'),
  'authenticated clients cannot update user_id'
);

select ok(
  not has_column_privilege('authenticated', 'public.history', 'plan_id', 'UPDATE'),
  'authenticated clients cannot update plan_id'
);

select ok(
  not has_column_privilege('authenticated', 'public.history', 'created_at', 'UPDATE'),
  'authenticated clients cannot update created_at'
);

select ok(
  not has_column_privilege('authenticated', 'public.history', 'updated_at', 'UPDATE'),
  'authenticated clients cannot update updated_at'
);

select ok(
  not has_column_privilege('authenticated', 'public.history', 'exported_at', 'UPDATE'),
  'authenticated clients cannot update exported_at'
);

-- RLS and Runtime checks as authenticated owner
select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'history_owner_id', 'role', 'authenticated')::text,
  true
);
set local role authenticated;

select is(
  (
    select count(*)::integer
    from public.history
    where user_id = :'history_owner_id'
      and plan_id = :'history_owner_plan_id'
  ),
  1,
  'authenticated users can read their own history'
);

select is(
  (
    select count(*)::integer
    from public.history
    where user_id = :'history_other_id'
  ),
  0,
  'authenticated users cannot read another user history'
);

select lives_ok(
  format(
    $sql$
      update public.history
      set keep_forever = true
      where user_id = %L and plan_id = %L
    $sql$,
    :'history_owner_id',
    :'history_owner_plan_id'
  ),
  'authenticated users can toggle keep_forever on their own history'
);

select is(
  (
    select keep_forever
    from public.history
    where user_id = :'history_owner_id'
      and plan_id = :'history_owner_plan_id'
  ),
  true,
  'keep_forever update persists for the owner'
);

select lives_ok(
  format(
    $sql$
      do $do$
      declare
        affected_rows integer;
      begin
        update public.history
        set keep_forever = true
        where user_id = %L;
        get diagnostics affected_rows = row_count;
        if affected_rows <> 0 then
          raise exception 'cross-user history update affected %% rows', affected_rows;
        end if;
      end
      $do$;
    $sql$,
    :'history_other_id'
  ),
  'authenticated users cannot update another user history row'
);

-- Verify that restricted operations fail with 42501 (insufficient_privilege)
select throws_ok(
  format(
    'insert into public.history (user_id, plan_id) values (%L, %L)',
    :'history_owner_id',
    :'history_owner_plan_id'
  ),
  '42501',
  null,
  'authenticated clients cannot insert history directly'
);

select throws_ok(
  format(
    'delete from public.history where user_id = %L',
    :'history_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot delete history directly'
);

select throws_ok(
  format(
    'update public.history set user_id = %L where plan_id = %L',
    :'history_other_id',
    :'history_owner_plan_id'
  ),
  '42501',
  null,
  'authenticated clients cannot update user_id on history'
);

select throws_ok(
  format(
    'update public.history set plan_id = %L where user_id = %L',
    :'history_other_plan_id',
    :'history_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot update plan_id on history'
);

select throws_ok(
  format(
    'update public.history set created_at = now() where user_id = %L',
    :'history_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot update created_at on history'
);

select throws_ok(
  format(
    'update public.history set updated_at = now() where user_id = %L',
    :'history_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot update updated_at on history'
);

select throws_ok(
  format(
    'update public.history set exported_at = now() where user_id = %L',
    :'history_owner_id'
  ),
  '42501',
  null,
  'authenticated clients cannot update exported_at on history'
);

select * from finish();

rollback;
