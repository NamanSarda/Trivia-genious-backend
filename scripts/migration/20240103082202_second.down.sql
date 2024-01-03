-- Filename: <timestamp>_drop_foreign_keys.sql

-- Drop foreign key constraints from the quizzes table
ALTER TABLE quizzes DROP CONSTRAINT IF EXISTS fk_quizzes_author;

-- Drop foreign key constraints from the participant table
ALTER TABLE participant DROP CONSTRAINT IF EXISTS fk_participant_quiz;
ALTER TABLE participant DROP CONSTRAINT IF EXISTS fk_participant_user;

-- Drop foreign key constraints from the question table
ALTER TABLE question DROP CONSTRAINT IF EXISTS fk_question_quiz;
