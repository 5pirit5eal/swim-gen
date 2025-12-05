-- Add removed_from_history column to feedback table
-- This tracks when a user removes a plan from their history but the plan is preserved due to existing feedback
ALTER TABLE feedback
ADD COLUMN IF NOT EXISTS removed_from_history boolean DEFAULT false;

-- Add index for efficient filtering
CREATE INDEX IF NOT EXISTS idx_feedback_removed_from_history ON feedback(removed_from_history) WHERE removed_from_history = true;
