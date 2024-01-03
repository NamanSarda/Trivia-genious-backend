CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    author_id INT,
    time_limit INT,
    difficulty VARCHAR(255),
    tags JSONB
);

CREATE TABLE participant (
    id SERIAL PRIMARY KEY,
    quiz_id INT,
    user_id INT,
    score INT
);

CREATE TABLE question (
    id SERIAL PRIMARY KEY,
    quiz_id INT,
    description TEXT,
    score INT,
    answer VARCHAR(255),
    options JSONB
);
