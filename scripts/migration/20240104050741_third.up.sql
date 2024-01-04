-- Up migration
ALTER TABLE participant
ADD COLUMN quiz_id INT;

ALTER TABLE question
ADD COLUMN quiz_id INT;

-- Add foreign key constraints
ALTER TABLE participant
ADD CONSTRAINT fk_participant_quiz
FOREIGN KEY (quiz_id) REFERENCES quizzes(id);

ALTER TABLE question
ADD CONSTRAINT fk_question_quiz
FOREIGN KEY (quiz_id) REFERENCES quizzes(id);
