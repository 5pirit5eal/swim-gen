begin;

select plan(14);

create temporary table test_feedback_context (
  owner_id uuid not null,
  recipient_id uuid not null,
  plan_id uuid not null,
  url_hash uuid not null
);

insert into test_feedback_context
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
  (select count(*)::integer from test_feedback_context),
  1,
  'seed contains an owner, recipient, and owned plan'
);

select owner_id, recipient_id, plan_id, url_hash
from test_feedback_context

\gset test_

select plan_id as unrelated_plan_id
from public.plans
where plan_id <> :'test_plan_id'
limit 1

\gset test_

insert into shared_plans (plan_id, user_id, url_hash)
values (:'test_plan_id', :'test_owner_id', :'test_url_hash');

select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_recipient_id', 'role', 'authenticated')::text,
  true
);
set local role authenticated;

select is(
  (select count(*)::integer from record_shared_plan(:'test_url_hash', 'link')),
  1,
  'recipient accepts the share before submitting feedback'
);

select throws_ok(
  format(
    'insert into feedback (user_id, plan_id, rating, comment) values (%L, %L, 5, %L)',
    :'test_recipient_id', :'test_plan_id', 'forged'
  ),
  '42501',
  null,
  'authenticated clients cannot insert feedback directly'
);

select throws_ok(
  format(
    'update feedback set rating = 1 where user_id = %L and plan_id = %L',
    :'test_owner_id', :'test_plan_id'
  ),
  '42501',
  null,
  'authenticated clients cannot update feedback directly'
);

select throws_ok(
  format(
    'delete from feedback where user_id = %L and plan_id = %L',
    :'test_owner_id', :'test_plan_id'
  ),
  '42501',
  null,
  'authenticated clients cannot delete feedback directly'
);

select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_recipient_id', 'role', 'authenticated')::text,
  true
);

select throws_ok(
  format(
    'insert into feedback (user_id, plan_id, rating, comment) values (%L, %L, 5, %L)',
    :'test_owner_id', :'test_plan_id', 'forged owner row'
  ),
  '42501',
  null,
  'authenticated clients cannot insert feedback for another user'
);

select throws_ok(
  format(
    'select public.submit_feedback(%L, %L, 5, false, 5, %L)',
    :'test_recipient_id', :'test_plan_id', 'forged function call'
  ),
  '42501',
  null,
  'authenticated clients cannot execute the trusted feedback operation'
);

set local role postgres;

select is(
  public.submit_feedback(:'test_owner_id', :'test_plan_id', 5, true, 7, 'owner feedback'),
  true,
  'owners can submit feedback for their plans'
);

select is(
  public.submit_feedback(:'test_recipient_id', :'test_plan_id', 4, true, 6, 'recipient feedback'),
  true,
  'accepted recipients can submit feedback for shared plans'
);

select is(
  public.submit_feedback(:'test_recipient_id', :'test_unrelated_plan_id', 4, true, 6, 'unrelated'),
  false,
  'users cannot submit feedback for an existing unrelated plan'
);

set local role authenticated;
delete from shared_history
where user_id = :'test_recipient_id'
  and plan_id = :'test_plan_id';

set local role postgres;
select is(
  public.submit_feedback(:'test_recipient_id', :'test_plan_id', 3, false, 4, 'revoked'),
  false,
  'revoked recipients cannot update their previous feedback'
);

set local role authenticated;
select is(
  (select count(*)::integer
   from feedback
   where user_id = :'test_recipient_id'
     and plan_id = :'test_plan_id'),
  1,
  'feedback remains stored after the recipient relationship is revoked'
);

select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_owner_id', 'role', 'authenticated')::text,
  true
);

select is(
  (select count(*)::integer
   from feedback
   where user_id = :'test_owner_id'
     and plan_id = :'test_plan_id'),
  1,
  'owners can read their own feedback rows'
);

select set_config(
  'request.jwt.claims',
  json_build_object('role', 'anon')::text,
  true
);
set local role anon;

select throws_ok(
  $$select count(*) from feedback$$,
  '42501',
  null,
  'anonymous users cannot read feedback'
);

select * from finish();

rollback;
