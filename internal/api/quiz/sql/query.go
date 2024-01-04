package query

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ayan-sh03/triviagenious-backend/config"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/models"
	"github.com/lib/pq"
)

func AddQuiz(quiz *models.Quiz) (int32, error) {

	var Db = config.Connect()

	tx := Db.MustBegin()
	tagsJSON, err := json.Marshal(quiz.Tags)
	if err != nil {
		// Handle the error
		fmt.Println("Error converting tags to JSON:", err)
		return 0, err
	}

	// Convert the JSON string to PostgreSQL jsonb format
	tagsString := string(tagsJSON)

	var quizID int32
	err = tx.QueryRow("INSERT INTO quizzes (author_id, time_limit, tags, difficulty) VALUES ($1, $2, $3, $4) RETURNING id", quiz.AuthorID, quiz.TimeLimit, tagsString, quiz.Difficulty).Scan(&quizID)
	if err != nil {
		// Handle the error
		log.Fatal(err)
	}

	fmt.Println("Quiz id is ", quizID)

	// for _, participant := range quiz.Participants {
	// 	_, err := tx.Exec("INSERT INTO participant (user_id, score) VALUES ($1, $2)", participant.UserID, participant.Score)
	// 	if err != nil {
	// 		// Handle the error
	// 		fmt.Println("Error inserting participant:", err)
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }

	// optionString := ""
	// for _, question := range quiz.Questions {
	// 	// Convert options for each question to JSON string
	// 	optionJson, err := json.Marshal(question.Options)
	// 	if err != nil {
	// 		// Handle the error
	// 		fmt.Println("Error converting options to JSON:", err)
	// 		return err
	// 	}

	// 	// Convert the JSON string to PostgreSQL jsonb format
	// 	optionString = string(optionJson)

	// 	// Insert each question with its unique set of options
	// 	_, err = tx.Exec("INSERT INTO question (description, score, answer, options) VALUES ($1, $2, $3, $4)", question.Description, question.Score, question.Answer, optionString)
	// 	if err != nil {
	// 		// Handle the error
	// 		fmt.Println("Error inserting question:", err)
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return quizID, nil
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

	query := `
		SELECT
		q.id AS quiz_id,
		q.author_id,
		q.time_limit,
		q.difficulty,
		ARRAY(SELECT jsonb_array_elements_text(q.tags)) AS tags,
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
	WHERE
		q.id = $1
		AND p.id IS NOT NULL  -- Exclude rows where participant_id is NULL
		AND que.id IS NOT NULL;  -- Exclude rows where question_id is NULL

    `
	row := tx.QueryRow(query, id)

	var quiz models.Quiz
	var participant models.Participant
	var question models.Question

	// Use pq.Array for scanning JSONB array into a slice of strings
	err := row.Scan(
		&quiz.ID,
		&quiz.AuthorID,
		&quiz.TimeLimit,
		&quiz.Difficulty,
		pq.Array(&quiz.Tags),
		&participant.ID,
		&participant.UserID,
		&participant.Score,
		&question.ID,
		&question.Description,
		&question.Score,
		&question.Answer,
		pq.Array(&question.Options),
	)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Add the participant and question to the quiz
	quiz.Participants = append(quiz.Participants, participant)
	quiz.Questions = append(quiz.Questions, question)

	return &quiz, nil
}

func DeleteParticipant(quizId, participantId int32) error {
	tx := config.Connect()
	_, err := tx.NamedExec("DELETE FROM participants WHERE id = :id AND quiz_id = :quiz_id", map[string]interface{}{
		"id":      participantId,
		"quiz_id": quizId,
	})
	if err != nil {

		return err
	}

	return nil
}

func CreateUser(user *models.User) error {

	tx := config.Connect()

	_, err := tx.Exec("Insert into users (username,email,password) values ($1,$2,$3)", &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Println(err)
		return err

	}
	return nil

}

func AddQuestions(userId, quizId int32, questions []models.Question) error {
	tx := config.Connect()
	fmt.Println(quizId)

	// for _, question := range questions {
	// 	// Insert each question with its unique set of options
	// 	_, err := tx.Exec(
	// 		"INSERT INTO question (quiz_id, description, score, answer, options) VALUES ($1, $2, $3, $4, $5)",
	// 		quizId, question.Description, question.Score, question.Answer, pq.Array(question.Options),
	// 	)
	// 	if err != nil {
	// 		// Handle the error
	// 		fmt.Println("Error inserting question:", err)
	// 		// tx.Rollback()
	// 		return err
	// 	}
	// }

	optionString := ""
	for _, question := range questions {
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
		_, err = tx.Exec("INSERT INTO question (quiz_id,description, score, answer, options) VALUES ($1, $2, $3, $4,$5) RETURNING id", quizId, question.Description, question.Score, question.Answer, optionString)
		if err != nil {
			// Handle the error
			fmt.Println("Error inserting question:", err)
			// tx.Rollback()
			return err
		}
	}

	return nil
}

func GetQuestionsById(quizId int32) ([]models.Question, error) {

	var questions []models.Question
	tx := config.Connect()
	err := tx.Select(
		&questions,
		"SELECT * FROM questions WHERE quiz_id = $1", quizId,
	)
	if err != nil {
		return nil, err
	}

	return questions, nil

}
