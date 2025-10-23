create schema if not exists private;
create table private.rate_limits (
  ip inet,
  request_at timestamp
);

-- add an index so that lookups are fast
create index rate_limits_ip_request_at_idx on private.rate_limits (ip, request_at desc);

-- function to check rate limits
create function public.check_request()
  returns void
  language plpgsql
  security definer
  as $$
declare
  req_method text := current_setting('request.method', true);
  req_ip inet := split_part(
    current_setting('request.headers', true)::json->>'x-forwarded-for',
    ',', 1)::inet;
  count_in_five_mins integer;
begin
  if req_method = 'GET' or req_method = 'HEAD' or req_method is null then
    -- rate limiting can't be done on GET and HEAD requests
    return;
  end if;

  select
    count(*) into count_in_five_mins
  from private.rate_limits
  where
    ip = req_ip and request_at between now() - interval '5 minutes' and now();

  if count_in_five_mins > 100 then
    raise sqlstate 'PGRST' using
      message = json_build_object(
        'message', 'Rate limit exceeded, try again after a while')::text,
      detail = json_build_object(
        'status',  420,
        'status_text', 'Enhance Your Calm')::text;
  end if;

  insert into private.rate_limits (ip, request_at) values (req_ip, now());
end;
  $$;

-- set the function to be called before each request
alter role authenticator
  set pgrst.db_pre_request = 'public.check_request';

notify pgrst, 'reload config';

-- pg_cron job to clean up old rate limit entries
select
  cron.schedule(
    'rate-limit-cleanup',
    '0 0 * * *', -- every day at midnight
    $$
      delete from private.rate_limits where request_at < now() - interval '1 day';
    $$
  );
