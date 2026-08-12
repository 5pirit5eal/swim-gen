begin;

select plan(26);

create temporary table test_shared_plan_context (
  owner_id uuid not null,
  recipient_id uuid not null,
  plan_id uuid not null,
  url_hash uuid not null
);

insert into test_shared_plan_context
select
  owner_user.id,
  recipient_user.id,
  donation.plan_id,
  gen_random_uuid()
from auth.users owner_user
join auth.users recipient_user
  on recipient_user.email = 'feedback-test-swimmer@example.com'
join donations donation
  on donation.user_id = owner_user.id
where owner_user.email = 'css-test-swimmer@example.com'
limit 1;

select is(
  (select count(*)::integer from test_shared_plan_context),
  1,
  'seed contains two users and an owner plan'
);

select owner_id, recipient_id, plan_id, url_hash
from test_shared_plan_context

\gset test_

insert into shared_plans (plan_id, user_id, url_hash)
values (:'test_plan_id', :'test_owner_id', :'test_url_hash');

select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_recipient_id', 'role', 'authenticated')::text,
  true
);
set local role authenticated;

select throws_ok(
  $$select count(*) from shared_plans$$,
  '42501',
  null,
  'recipients cannot enumerate shared plan rows'
);

select is(
  (select count(*)::integer from plans where plan_id = :'test_plan_id'),
  0,
  'recipients cannot read plans before accepting a share'
);

select is(
  (select count(*)::integer from get_shared_plan_by_hash(:'test_url_hash')),
  1,
  'valid share hash resolves through the scoped RPC'
);

select is(
  (select count(*)::integer from get_shared_plan_by_hash(gen_random_uuid())),
  0,
  'invalid share hash returns no rows'
);

select is(
  (select count(*)::integer from record_shared_plan(:'test_url_hash', 'link')),
  1,
  'recipient can accept a valid share'
);

select is(
  (select count(*)::integer
   from shared_history
   where user_id = :'test_recipient_id'
     and plan_id = :'test_plan_id'
     and shared_by = :'test_owner_id'),
  1,
  'share acceptance derives the recipient and owner relationship'
);

select is(
  (select count(*)::integer from plans where plan_id = :'test_plan_id'),
  1,
  'recipients can read an accepted shared plan'
);

set local role postgres;

select is(
  (select share_count from shared_plans where plan_id = :'test_plan_id'),
  1,
  'share acceptance updates the current recipient count'
);

set local role authenticated;

select is(
  (select count(*)::integer from record_shared_plan(:'test_url_hash', 'link')),
  1,
  'repeated share acceptance is idempotent'
);

set local role postgres;

select is(
  (select share_count from shared_plans where plan_id = :'test_plan_id'),
  1,
  'repeated acceptance does not inflate the count'
);

set local role authenticated;

select throws_ok(
  format(
    'insert into shared_history (user_id, plan_id, shared_by) values (%L, %L, %L)',
    :'test_recipient_id', :'test_plan_id', :'test_recipient_id'
  ),
  '42501',
  null,
  'recipients cannot forge shared history rows'
);

select throws_ok(
  format(
    'update shared_plans set share_count = 99 where plan_id = %L',
    :'test_plan_id'
  ),
  '42501',
  null,
  'recipients cannot update share metadata'
);

select throws_ok(
  format(
    'delete from plans where plan_id = %L',
    :'test_plan_id'
  ),
  '42501',
  null,
  'recipients cannot delete canonical plans'
);

select throws_ok(
  format(
    'update plans set title = %L where plan_id = %L',
    'forged title', :'test_plan_id'
  ),
  '42501',
  null,
  'recipients cannot update canonical plans'
);

select throws_ok(
  format(
    'insert into history (user_id, plan_id) values (%L, %L)',
    :'test_recipient_id', :'test_plan_id'
  ),
  '42501',
  null,
  'recipients cannot manufacture ownership history'
);

select throws_ok(
  format(
    'insert into donations (user_id, plan_id) values (%L, %L)',
    :'test_recipient_id', :'test_plan_id'
  ),
  '42501',
  null,
  'recipients cannot manufacture donation ownership'
);

select throws_ok(
  format(
    'insert into shared_plans (plan_id, user_id) values (%L, %L)',
    :'test_plan_id', :'test_recipient_id'
  ),
  '42501',
  null,
  'recipients cannot publish arbitrary plans'
);

select throws_ok(
  format(
    'update history set plan_id = %L where user_id = %L and plan_id = %L',
    :'test_plan_id', :'test_recipient_id', :'test_plan_id'
  ),
  '42501',
  null,
  'recipients cannot repoint history to arbitrary plans'
);

delete from shared_history
where user_id = :'test_recipient_id'
  and plan_id = :'test_plan_id';

set local role postgres;

select is(
  (select share_count from shared_plans where plan_id = :'test_plan_id'),
  0,
  'removing recipient history decrements the current count'
);

select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_owner_id', 'role', 'authenticated')::text,
  true
);
set local role authenticated;

select is(
  (select count(*)::integer from plans where plan_id = :'test_plan_id'),
  1,
  'owners can read their plans'
);

select is(
  (select count(*)::integer from get_shared_plan_by_hash(:'test_url_hash')),
  1,
  'owners can use the shared-plan lookup RPC'
);

select is(
  (select count(*)::integer from record_shared_plan(:'test_url_hash', 'link')),
  0,
  'owners cannot record themselves as recipients'
);

select set_config(
  'request.jwt.claims',
  json_build_object('role', 'anon')::text,
  true
);
set local role anon;

select throws_ok(
  $$select count(*) from shared_plans$$,
  '42501',
  null,
  'anonymous users cannot enumerate shared plan rows'
);

select throws_ok(
  format('select count(*) from plans where plan_id = %L', :'test_plan_id'),
  '42501',
  null,
  'anonymous users cannot query canonical plans directly'
);

select is(
  (select count(*)::integer from get_shared_plan_by_hash(:'test_url_hash')),
  1,
  'anonymous users can resolve a valid bearer hash'
);

select * from finish();

rollback;
