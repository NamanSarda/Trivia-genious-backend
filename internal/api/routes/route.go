package routes

import (
	aicontroller "github.com/ayan-sh03/triviagenious-backend/internal/api/ai/controller"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/middleware"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/quiz/controller"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	quizRouter := r.PathPrefix("/api/v1/quiz").Subrouter()

	quizRouter.Use(middleware.Auth)

	quizRouter.HandleFunc("", controller.CreateQuiz).Methods("POST")
	quizRouter.HandleFunc("/{quizID}", controller.DeleteQuiz).Methods("DELETE")

	aiRouter := r.PathPrefix("/api/v1/ai").Subrouter()

	aiRouter.HandleFunc("", aicontroller.GetQuestionFromAi).Methods("POST")
	// aiRouter.HandleFunc("/pdf", aicontroller.Reader).Methods("GET")

	return r
}
