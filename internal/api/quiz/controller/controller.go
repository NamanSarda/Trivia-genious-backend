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
	query "github.com/ayan-sh03/triviagenious-backend/internal/api/quiz/sql"
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

	var iderr error

	quiz.AuthorID, iderr = GetUserIDFromCookie(r)

	log.Print("ID from cookie", quiz.AuthorID)

	if iderr != nil {
		log.Fatal(iderr)
	}

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

	fmt.Println(quiz)

	// fmt.Printf("%+v\n", quiz.Participants)
	fmt.Printf("%+v\n", quiz.Questions)

	quizId, dberr := query.AddQuiz(&quiz)
	if dberr != nil {
		log.Fatal("Failed to add query")
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error "+dberr.Error())
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Quiz added successfully", "id": quizId})
}

func DeleteQuiz(w http.ResponseWriter, r *http.Request) {

	user_id, err := GetUserIDFromCookie(r)
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Cannot retrive ID from cookie")
		return
	}

	// var user_id int32 = 0
	var quiz models.Quiz

	vars := mux.Vars(r)
	quizIDStr, ok := vars["quizID"]

	// Check if the "quizID" parameter is present
	if !ok {
		http.Error(w, "Quiz ID not found in URL", http.StatusBadRequest)
		return
	}

	tx := config.Connect()

	err = tx.QueryRow("SELECT id,author_id FROM quizzes WHERE id = $1", quizIDStr).
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

func AddParticipants(w http.ResponseWriter, r *http.Request) {

	// get user id from cookie

	var participant models.Participant
	var err error
	participant.UserID, err = GetUserIDFromCookie(r)

	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&participant); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Incomplete Request Body")
		return
	}

	//validation
	if participant.QuizID == 0 {
		util.RespondWithError(w, http.StatusBadRequest, "No Quiz id in the Request Body ")
		return
	}

	//!

	fmt.Println(participant)

	err = query.AddParticipant(&participant)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error While Adding Participant")
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Added Participant Successfully"})
}

func UpdateParticipantScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	participantIDStr, ok := vars["participantID"]

	var score int32

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&score); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Provide Updated Score! ")
		return
	}

	if !ok {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid Participant Id")
		return
	}

	participantId, err := strconv.Atoi(participantIDStr)

	if err != nil {
		log.Fatal("Unable to Convert string to int")
	}

	err = query.UpdateScore(int32(participantId), score)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error While Executing Query !")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Updated Successfully"})

}

func DeleteParticipants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	partIDStr, quizIdstr := vars["partId"], vars["quizId"]

	partId, err := strconv.Atoi(partIDStr)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	quizId, err := strconv.Atoi(quizIdstr)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = query.DeleteParticipant(int32(quizId), int32(partId))

	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Deleted Successfully"})

}

func GetQuizById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quizIDStr, ok := vars["quizID"]

	if !ok {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid String in the Url")
		return
	}

	quiz_id, err := strconv.Atoi(quizIDStr)

	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid String in the Url")
		return
	}

	quiz, sqlErr := query.GetQuizById(int32(quiz_id))

	if sqlErr != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal server error "+sqlErr.Error())
		return

	}

	util.RespondWithJSON(w, http.StatusOK, quiz)

}

func AddQuestions(w http.ResponseWriter, r *http.Request) {
	userId, _ := GetUserIDFromCookie(r)

	vars := mux.Vars(r)
	quizIdStr, ok := vars["quizId"]
	if !ok || quizIdStr == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Missing or empty 'quizId' parameter")
		return
	}

	quizId, err := strconv.Atoi(quizIdStr)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid 'quizId' parameter: "+err.Error())
		return
	}
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Please Login Again")

		return
	}

	var questions []models.Question
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&questions)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Print the decoded questions
	for _, q := range questions {
		fmt.Printf("ID: %d, Description: %s, Answer: %s\n", q.ID, q.Description, q.Answer)
	}

	err = query.AddQuestions(userId, int32(quizId), questions)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error "+err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"message": "OK", "questions": questions})
}

func GetAllQuestionsFromQuizId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	quizIdStr := vars["quizId"]

	quizId, err := strconv.Atoi(quizIdStr)

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error"+err.Error())

	}

	questions, err := query.GetQuestionsById(int32(quizId))

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error"+err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "OK", "questions": questions})
}
