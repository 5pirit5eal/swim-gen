-- Add allow_sharing column to donations table
-- This column indicates whether the user allows their uploaded plan to be used in the RAG system

ALTER TABLE donations
ADD COLUMN IF NOT EXISTS allow_sharing BOOLEAN NOT NULL DEFAULT FALSE;

-- Add comment to document the column purpose
COMMENT ON COLUMN donations.allow_sharing IS 'Indicates whether the user allows their uploaded plan to be used in the RAG system for generating recommendations';
