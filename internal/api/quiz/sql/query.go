package query

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ayan-sh03/triviagenious-backend/config"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/models"
)

func AddQuiz(quiz *models.Quiz) error {

	var Db = config.Connect()

	tx := Db.MustBegin()
	tagsJSON, err := json.Marshal(quiz.Tags)
	if err != nil {
		// Handle the error
		fmt.Println("Error converting tags to JSON:", err)
		return err
	}

	// Convert the JSON string to PostgreSQL jsonb format
	tagsString := string(tagsJSON)

	tx.MustExec("INSERT INTO quizzes (author_id, time_limit, tags, difficulty) VALUES ($1, $2, $3, $4)", quiz.AuthorID, quiz.TimeLimit, tagsString, quiz.Difficulty)

	for _, participant := range quiz.Participants {
		_, err := tx.Exec("INSERT INTO participant (user_id, score) VALUES ($1, $2)", participant.UserID, participant.Score)
		if err != nil {
			// Handle the error
			fmt.Println("Error inserting participant:", err)
			tx.Rollback()
			return err
		}
	}

	optionString := ""
	for _, question := range quiz.Questions {
		// Convert options for each question to JSON string
		optionJson, err := json.Marshal(question.Options)
		if err != nil {
			// Handle the error
			fmt.Println("Error converting options to JSON:", err)
			return err
		}

		// Convert the JSON string to PostgreSQL jsonb format
		optionString = string(optionJson)

		// Insert each question with its unique set of options
		_, err = tx.Exec("INSERT INTO question (description, score, answer, options) VALUES ($1, $2, $3, $4)", question.Description, question.Score, question.Answer, optionString)
		if err != nil {
			// Handle the error
			fmt.Println("Error inserting question:", err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func AddParticipant(part *models.Participant) error {

	tx := config.Connect()

	_, err := tx.Exec("INSERT INTO participant (quiz_id,user_id,score) values ($1,$2,$3) ", part.QuizID, part.UserID, part.Score)

	if err != nil {
		return err
	}

	return nil
}

func UpdateScore(participantId int32, score int32) error {

	//update score with the given score

	tx := config.Connect()

	query := "UPDATE participant SET score = :newScore WHERE id = :participantID"

	// Named parameters to bind values
	params := map[string]interface{}{
		"newScore":      score,
		"participantID": participantId,
	}

	_, err := tx.NamedExec(query, params)

	if err != nil {
		return err
	}

	return nil
}

func GetQuizById(id int32) (*models.Quiz, error) {

	tx := config.Connect()

	query := `SELECT
   q.id AS quiz_id,
    q.author_id,
    q.time_limit,
    q.difficulty,
    q.tags,
	  p.id AS participant_id,
    p.user_id,
    p.score AS participant_score,
    que.id AS question_id,
    que.description AS question_description,
    que.score AS question_score,
    que.answer AS correct_answer,
    que.options AS question_options
FROM
    quizzes q
LEFT JOIN LATERAL (
    SELECT id, user_id, score
    FROM participant
WHERE quiz_id = q.id
) p ON true
LEFT JOIN LATERAL (
    SELECT id, description, score, answer, options
    FROM question
    WHERE quiz_id = q.id
) que ON true
WHERE q.id = $1;
`
	row := tx.QueryRow(query, id)

	var quiz models.Quiz
	// var part models.Participant
	// var question models.Question

	err := row.Scan(&quiz.ID, &quiz.AuthorID, &quiz.Difficulty, &quiz.Tags, &quiz.TimeLimit, &quiz.Participants, &quiz.Questions)

	fmt.Println(quiz.Participants)
	fmt.Println(quiz.Questions)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &quiz, nil
}
