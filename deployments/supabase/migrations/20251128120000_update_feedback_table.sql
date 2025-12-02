-- Add was_swam and difficulty_rating columns to feedback table
ALTER TABLE feedback
ADD COLUMN IF NOT EXISTS was_swam boolean DEFAULT false,
ADD COLUMN IF NOT EXISTS difficulty_rating int CHECK (difficulty_rating BETWEEN 1 AND 10);

-- Keep updated_at fresh
create or replace function public.update_feedback_updated_at()
returns trigger
language plpgsql as $$
begin
    new.updated_at := now();
    return new;
end;
$$;



CREATE OR REPLACE TRIGGER update_feedback_updated_at_trg
    BEFORE UPDATE ON feedback
    FOR EACH ROW
    EXECUTE FUNCTION update_feedback_updated_at();


-- Add exported_at column to history table
ALTER TABLE history
ADD COLUMN IF NOT EXISTS exported_at timestamptz;
