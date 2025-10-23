-- Enable RLS on Feedback with user-specific policies
ALTER TABLE public.feedback ENABLE ROW LEVEL SECURITY;
CREATE POLICY feedback_select_own ON public.feedback FOR SELECT TO authenticated USING ((SELECT auth.uid()) = user_id);
CREATE POLICY feedback_insert_own ON public.feedback FOR INSERT TO authenticated WITH CHECK ((SELECT auth.uid()) = user_id);
CREATE POLICY feedback_update_own ON public.feedback FOR UPDATE TO authenticated USING ((SELECT auth.uid()) = user_id) WITH CHECK ((SELECT auth.uid()) = user_id);
CREATE POLICY feedback_delete_own ON public.feedback FOR DELETE TO authenticated USING ((SELECT auth.uid()) = user_id);

-- Enable RLS on Donations with user-specific policies
ALTER TABLE public.donations ENABLE ROW LEVEL SECURITY;
CREATE POLICY donations_select_own ON public.donations FOR SELECT TO authenticated USING ((SELECT auth.uid()) = user_id);
CREATE POLICY donations_insert_own ON public.donations FOR INSERT TO authenticated WITH CHECK ((SELECT auth.uid()) = user_id);
CREATE POLICY donations_update_own ON public.donations FOR UPDATE TO authenticated USING ((SELECT auth.uid()) = user_id) WITH CHECK ((SELECT auth.uid()) = user_id);
CREATE POLICY donations_delete_own ON public.donations FOR DELETE TO authenticated USING ((SELECT auth.uid()) = user_id);

-- Enable RLS on tables that should not be publicly readable
ALTER TABLE public.scraped ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.embeddings ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.embedders ENABLE ROW LEVEL SECURITY;

-- Enable RLS on Plans with policy to only allow write access for plan_ids contained in donations
ALTER TABLE public.plans ENABLE ROW LEVEL SECURITY;

CREATE POLICY plans_access_for_donators ON public.plans
FOR ALL
TO authenticated
USING (
  EXISTS (
    SELECT 1
    FROM public.donations
    WHERE donations.plan_id = plans.plan_id
      AND donations.user_id = auth.uid()
  )
)
WITH CHECK (
  EXISTS (
    SELECT 1
    FROM public.donations
    WHERE donations.plan_id = plans.plan_id
      AND donations.user_id = auth.uid()
  )
);



-- User-specific indexes for performance
CREATE INDEX idx_feedback_user_id ON public.feedback(user_id);
CREATE INDEX idx_donations_user_id ON public.donations(user_id);
