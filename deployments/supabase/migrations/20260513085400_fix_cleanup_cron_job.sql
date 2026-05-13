-- Unschedule the old job
SELECT cron.unschedule('clean-up-old-plans');

-- Schedule the fixed job
SELECT cron.schedule(
    'clean-up-old-plans',
    '0 0 * * *', -- Every day at midnight
    $$
    DO $do$
    BEGIN
      -- Delete history entries older than 2 days
      DELETE FROM public.history WHERE created_at < now() - interval '2 days' AND keep_forever = false;

      -- Delete plans that are not referenced anymore
      DELETE FROM public.plans p
      WHERE NOT EXISTS (SELECT 1 FROM public.history h WHERE h.plan_id = p.plan_id)
        AND NOT EXISTS (SELECT 1 FROM public.donations d WHERE d.plan_id = p.plan_id)
        AND NOT EXISTS (SELECT 1 FROM public.feedback f WHERE f.plan_id = p.plan_id)
        AND NOT EXISTS (SELECT 1 FROM public.scraped s WHERE s.plan_id = p.plan_id);
    END
    $do$;
    $$
);
