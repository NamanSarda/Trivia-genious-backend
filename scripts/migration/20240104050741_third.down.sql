ALTER TABLE participant
DROP CONSTRAINT IF EXISTS fk_participant_quiz;

ALTER TABLE question
DROP CONSTRAINT IF EXISTS fk_question_quiz;

-- Remove quiz_id columns
ALTER TABLE participant
DROP COLUMN IF EXISTS quiz_id;

ALTER TABLE question
DROP COLUMN IF EXISTS quiz_id;