-- Add question column to nodes
ALTER TABLE nodes
ADD COLUMN IF NOT EXISTS question text;

ALTER TABLE nodes
ADD CONSTRAINT nodes_question_check
CHECK (question IS NULL OR char_length(question) <= 30);
