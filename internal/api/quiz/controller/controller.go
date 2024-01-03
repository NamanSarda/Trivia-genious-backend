package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ayan-sh03/triviagenious-backend/config"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/models"
	"github.com/ayan-sh03/triviagenious-backend/internal/util"
	"github.com/gorilla/mux"
)

func GetUserIDFromCookie(r *http.Request) (int32, error) {
	// Retrieve the cookie from the request
	// cookies := r.Cookies()

	// Print all available cookies
	// log.Println("All Cookies:", cookies)

	cookie, err := r.Cookie("user_id")
	if err != nil {
		log.Println("Error: Cookie not found -", err)
		return 0, err
	}

	// Get the value of the "user_id" cookie
	userIdStr := cookie.Value

	// Parse the user ID string into an int32
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		log.Println("Error: Unable to parse user ID -", err)
		return 0, err
	}

	// Return the user ID as int32
	return int32(userId), nil
}
func CreateQuiz(w http.ResponseWriter, r *http.Request) {

	var quiz models.Quiz

	// var iderr error

	// quiz.AuthorID, iderr = GetUserIDFromCookie(r)

	// log.Printf("ID from cookie", quiz.AuthorID)

	// if iderr != nil {
	// 	log.Fatal(iderr)
	// }

	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return
	}

	// Close the request body
	defer r.Body.Close()

	err = json.Unmarshal(jsonData, &quiz)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Printf("%+v\n", quiz.Participants)
	fmt.Printf("%+v\n", quiz.Questions)

	dberr := AddQuiz(&quiz)
	if dberr != nil {
		log.Fatal("Failed to add query")
	}
}

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

func DeleteQuiz(w http.ResponseWriter, r *http.Request) {

	// user_id, err := GetUserIDFromCookie(r)
	// if err != nil {
	// 	util.RespondWithError(w, http.StatusUnauthorized, "Cannot retrive ID from cookie")
	// 	return
	// }
	var user_id int32 = 0
	var quiz models.Quiz

	vars := mux.Vars(r)
	quizIDStr, ok := vars["quizID"]

	// Check if the "quizID" parameter is present
	if !ok {
		http.Error(w, "Quiz ID not found in URL", http.StatusBadRequest)
		return
	}

	tx := config.Connect()

	err := tx.QueryRow("SELECT id,author_id FROM quizzes WHERE id = $1", quizIDStr).
		Scan(&quiz.ID, &quiz.AuthorID)

	if err != nil {
		// Handle the error (e.g., no rows found)
		if err == sql.ErrNoRows {
			fmt.Println("No rows found.")
		} else {
			// Handle other errors
			fmt.Println("Error:", err)
		}
	} else {
		// Now you can access the authorID
		fmt.Printf("Author ID: %d\n", quiz.AuthorID)

		if user_id != quiz.AuthorID {
			util.RespondWithError(w, http.StatusUnauthorized, "Not Authorized to Delete the Quiz")
			return
		}

	}

	_, err = tx.Exec("DELETE FROM quizzes WHERE id = $1", quiz.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Could not delete Quiz ")
		return
	}

	util.RespondWithJSON(w, http.StatusNoContent, "Deleted Successfully")
}
