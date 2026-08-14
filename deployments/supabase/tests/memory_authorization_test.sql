begin;

select plan(13);

create temporary table test_memory_authorization_context (
  owner_id uuid not null,
  other_id uuid not null,
  plan_id uuid not null,
  message_id uuid not null
);

insert into test_memory_authorization_context
select
  owner_user.id,
  other_user.id,
  donation.plan_id,
  gen_random_uuid()
from auth.users owner_user
join auth.users other_user
  on other_user.email = 'feedback-test-swimmer@example.com'
join donations donation
  on donation.user_id = owner_user.id
where owner_user.email = 'css-test-swimmer@example.com'
limit 1;

select is(
  (select count(*)::integer from test_memory_authorization_context),
  1,
  'seed contains an owner, unrelated user, and owner plan'
);

select owner_id, other_id, plan_id, message_id
from test_memory_authorization_context

\gset test_

set local role postgres;

insert into memory (
  id,
  plan_id,
  user_id,
  role,
  content,
  plan_snapshot
)
values (
  :'test_message_id',
  :'test_plan_id',
  :'test_owner_id',
  'ai',
  'private owner response',
  jsonb_build_object('plan_id', :'test_plan_id', 'title', 'private snapshot')
);

set local role authenticated;
select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_other_id', 'role', 'authenticated')::text,
  true
);

select is(
  (select count(*)::integer from plans where plan_id = :'test_plan_id'),
  0,
  'unrelated users cannot read another account plan'
);

select throws_ok(
  format('select count(*) from memory where plan_id = %L', :'test_plan_id'),
  '42501',
  null,
  'unrelated users cannot read another account messages or snapshots'
);

select throws_ok(
  format(
    'insert into memory (plan_id, user_id, role, content) values (%L, %L, %L, %L)',
    :'test_plan_id', :'test_other_id', 'user', 'forged message'
  ),
  '42501',
  null,
  'unrelated users cannot append messages to another account plan'
);

select throws_ok(
  format(
    'insert into memory (plan_id, user_id, role, content) values (%L, %L, %L, %L)',
    :'test_plan_id', :'test_owner_id', 'user', 'forged owner message'
  ),
  '42501',
  null,
  'clients cannot impersonate the plan owner in memory rows'
);

select throws_ok(
  format(
    'update memory set content = %L, plan_snapshot = %L where id = %L',
    'forged snapshot', '{"plan_id":"' || :'test_plan_id' || '","title":"forged"}', :'test_message_id'
  ),
  '42501',
  null,
  'unrelated users cannot alter messages or snapshots'
);

select throws_ok(
  format('delete from memory where id = %L', :'test_message_id'),
  '42501',
  null,
  'unrelated users cannot delete another account messages'
);

set local role anon;
select set_config('request.jwt.claims', json_build_object('role', 'anon')::text, true);

select throws_ok(
  format('select count(*) from memory where id = %L', :'test_message_id'),
  '42501',
  null,
  'anonymous users cannot read chat messages or snapshots'
);

select throws_ok(
  format('delete from memory where id = %L', :'test_message_id'),
  '42501',
  null,
  'anonymous users cannot delete chat messages'
);

set local role postgres;

select is(
  (select count(*)::integer
   from memory
   where id = :'test_message_id'
     and user_id = :'test_owner_id'
     and plan_id = :'test_plan_id'
     and plan_snapshot->>'plan_id' = :'test_plan_id'),
  1,
  'trusted backend can read the owner message and matching snapshot'
);

select is(
  (select content from memory where id = :'test_message_id'),
  'private owner response',
  'failed cross-account writes leave the owner message unchanged'
);

select is(
  (select count(*)::integer from memory where plan_id = :'test_plan_id'),
  1,
  'failed cross-account writes do not append messages'
);

select is(
  (select count(*)::integer from plans where plan_id = :'test_plan_id'),
  1,
  'failed cross-account writes do not remove the owner plan'
);

select * from finish();

rollback;
