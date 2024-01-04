package routes

import (
	aicontroller "github.com/ayan-sh03/triviagenious-backend/internal/api/ai/controller"
	authcontroller "github.com/ayan-sh03/triviagenious-backend/internal/api/auth/controller"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/middleware"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/quiz/controller"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	quizRouter := r.PathPrefix("/api/v1/quiz").Subrouter()
	userRouter := r.PathPrefix("/api/v1/user").Subrouter()
	participantRouter := r.PathPrefix("/api/v1/participant").Subrouter()
	// participantRouter := r.PathPrefix("/api/v1/").Subrouter()

	quizRouter.Use(middleware.Auth)

	userRouter.HandleFunc("/register", authcontroller.RegisterUserController).Methods("POST")
	userRouter.HandleFunc("/login", authcontroller.LoginController).Methods("POST")

	quizRouter.HandleFunc("", controller.CreateQuiz).Methods("POST")
	quizRouter.HandleFunc("/{quizID}", controller.DeleteQuiz).Methods("DELETE")
	quizRouter.HandleFunc("/{quizID}", controller.GetQuizById).Methods("GET")
	quizRouter.HandleFunc("/{quizId:[0-9]+}", controller.AddQuestions).Methods("POST")

	participantRouter.HandleFunc("/", controller.AddParticipants).Methods("POST")
	participantRouter.HandleFunc("/{participantId:[0-9]+}", controller.DeleteParticipants).Methods("DELETE")
	participantRouter.HandleFunc("/score/{participantId:[0-9]+}", controller.UpdateParticipantScore).Methods("POST")

	aiRouter := r.PathPrefix("/api/v1/ai").Subrouter()

	aiRouter.HandleFunc("", aicontroller.GetQuestionFromAi).Methods("POST")

	return r
}
