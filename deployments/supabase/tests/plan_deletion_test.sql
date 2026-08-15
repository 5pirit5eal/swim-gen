begin;

select plan(12);

create temporary table test_plan_deletion_context (
  owner_id uuid not null,
  recipient_id uuid not null,
  private_plan_id uuid not null,
  shared_plan_id uuid not null,
  donated_plan_id uuid not null,
  feedback_plan_id uuid not null,
  url_hash uuid not null
);

insert into test_plan_deletion_context
select
  owner_user.id,
  recipient_user.id,
  gen_random_uuid(),
  gen_random_uuid(),
  gen_random_uuid(),
  gen_random_uuid(),
  gen_random_uuid()
from auth.users owner_user
join auth.users recipient_user
  on recipient_user.email = 'feedback-test-swimmer@example.com'
where owner_user.email = 'css-test-swimmer@example.com'
limit 1;

select is(
  (select count(*)::integer from test_plan_deletion_context),
  1,
  'seed contains an owner and recipient context'
);

select owner_id, recipient_id, private_plan_id, shared_plan_id, donated_plan_id, feedback_plan_id, url_hash
from test_plan_deletion_context

\gset test_

set local role postgres;

-- 1. Create a private plan with history and memory
insert into plans (plan_id, title, description, plan_table)
values (:'test_private_plan_id', 'Private Plan', 'Desc', '[]'::jsonb);

insert into history (user_id, plan_id)
values (:'test_owner_id', :'test_private_plan_id');

insert into memory (plan_id, user_id, role, content)
values (:'test_private_plan_id', :'test_owner_id', 'user', 'chat message');

-- 2. Create a shared plan with owner share and recipient history
insert into plans (plan_id, title, description, plan_table)
values (:'test_shared_plan_id', 'Shared Plan', 'Desc', '[]'::jsonb);

insert into history (user_id, plan_id)
values (:'test_owner_id', :'test_shared_plan_id');

insert into shared_plans (user_id, plan_id, url_hash)
values (:'test_owner_id', :'test_shared_plan_id', :'test_url_hash');

insert into shared_history (user_id, plan_id, shared_by, share_method)
values (:'test_recipient_id', :'test_shared_plan_id', :'test_owner_id', 'link');

-- 3. Create a donated plan
insert into plans (plan_id, title, description, plan_table)
values (:'test_donated_plan_id', 'Donated Plan', 'Desc', '[]'::jsonb);

insert into donations (user_id, plan_id)
values (:'test_owner_id', :'test_donated_plan_id');

-- 4. Create a feedback-bearing plan
insert into plans (plan_id, title, description, plan_table)
values (:'test_feedback_plan_id', 'Feedback Plan', 'Desc', '[]'::jsonb);

insert into history (user_id, plan_id)
values (:'test_owner_id', :'test_feedback_plan_id');

insert into feedback (user_id, plan_id, rating, comment, was_swam, difficulty_rating, removed_from_history)
values (:'test_owner_id', :'test_feedback_plan_id', 5, 'Great plan', true, 5, false);

-- Test A: Recipient cannot delete the shared plan from plans table
set local role authenticated;
select set_config(
  'request.jwt.claims',
  json_build_object('sub', :'test_recipient_id', 'role', 'authenticated')::text,
  true
);

select throws_ok(
  format('delete from plans where plan_id = %L', :'test_shared_plan_id'),
  '42501',
  null,
  'recipient cannot delete owner shared plan directly'
);

-- Test B: Recipient can delete only their own row in shared_history
delete from shared_history where plan_id = :'test_shared_plan_id';

select is(
  (select count(*)::integer from shared_history where user_id = :'test_recipient_id' and plan_id = :'test_shared_plan_id'),
  0,
  'recipient can remove plan from their own shared_history'
);

set local role postgres;

select is(
  (select count(*)::integer from shared_plans where plan_id = :'test_shared_plan_id'),
  1,
  'recipient deletion of shared_history does not affect shared_plans'
);

select is(
  (select count(*)::integer from plans where plan_id = :'test_shared_plan_id'),
  1,
  'recipient deletion of shared_history does not delete owner plan'
);

-- Test C: Owner deletion of private plan cascades to history and memory
set local role postgres;
delete from plans where plan_id = :'test_private_plan_id';

select is(
  (select count(*)::integer from history where plan_id = :'test_private_plan_id'),
  0,
  'deleting private plan cascades to history'
);

select is(
  (select count(*)::integer from memory where plan_id = :'test_private_plan_id'),
  0,
  'deleting private plan cascades to memory'
);

-- Test D: Owner deletion of shared plan cascades to shared_plans and shared_history
-- Re-insert recipient shared history to test cascade
insert into shared_history (user_id, plan_id, shared_by, share_method)
values (:'test_recipient_id', :'test_shared_plan_id', :'test_owner_id', 'link');

delete from plans where plan_id = :'test_shared_plan_id';

select is(
  (select count(*)::integer from shared_plans where plan_id = :'test_shared_plan_id'),
  0,
  'deleting shared plan cascades to shared_plans'
);

select is(
  (select count(*)::integer from shared_history where plan_id = :'test_shared_plan_id'),
  0,
  'deleting shared plan cascades to shared_history'
);

-- Test E: Owner deletion of donated plan cascades to donations
delete from plans where plan_id = :'test_donated_plan_id';

select is(
  (select count(*)::integer from donations where plan_id = :'test_donated_plan_id'),
  0,
  'deleting donated plan cascades to donations'
);

-- Test F: Feedback-bearing plan deletion logic (preserving plan and marking feedback)
delete from history where plan_id = :'test_feedback_plan_id' and user_id = :'test_owner_id';
update feedback set removed_from_history = true where plan_id = :'test_feedback_plan_id';

select is(
  (select count(*)::integer from plans where plan_id = :'test_feedback_plan_id'),
  1,
  'feedback plan entity is preserved after removal from history'
);

select is(
  (select removed_from_history from feedback where plan_id = :'test_feedback_plan_id'),
  true,
  'feedback row is marked with removed_from_history = true'
);

select * from finish();

rollback;
