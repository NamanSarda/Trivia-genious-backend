-- Filename: <timestamp>_add_foreign_keys.sql

-- Add foreign key constraints to the quizzes table
ALTER TABLE quizzes
    ADD CONSTRAINT fk_quizzes_author
    FOREIGN KEY (author_id)
    REFERENCES auth_user(id)
    ON DELETE CASCADE;

-- Add foreign key constraints to the participant table
ALTER TABLE participant
    ADD CONSTRAINT fk_participant_quiz
    FOREIGN KEY (quiz_id)
    REFERENCES quizzes(id)
    ON DELETE CASCADE;

ALTER TABLE participant
    ADD CONSTRAINT fk_participant_user
    FOREIGN KEY (user_id)
    REFERENCES auth_user(id)
    ON DELETE CASCADE;

-- Add foreign key constraints to the question table
ALTER TABLE question
    ADD CONSTRAINT fk_question_quiz
    FOREIGN KEY (quiz_id)
    REFERENCES quizzes(id)
    ON DELETE CASCADE;
