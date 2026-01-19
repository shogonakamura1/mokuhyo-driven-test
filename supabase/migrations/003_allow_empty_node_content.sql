-- Allow empty content for nodes
-- Remove the check constraint that requires content to be at least 1 character

ALTER TABLE nodes
DROP CONSTRAINT IF EXISTS nodes_content_check;

-- Add a new constraint that allows empty content but still enforces max length
ALTER TABLE nodes
ADD CONSTRAINT nodes_content_check CHECK (char_length(content) <= 200);
