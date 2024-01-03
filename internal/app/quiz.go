package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ayan-sh03/triviagenious-backend/config"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/routes"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB := config.Connect()

	if err := DB.Ping(); err != nil {
		log.Fatal("Could not ping db")
	}
	r := routes.SetupRoutes()
	// CORS middleware
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Wrap your router with the CORS middleware
	corsRouter := handlers.CORS(headersOk, originsOk, methodsOk)(r)
	fmt.Println("Server started on port : 8080")
	log.Fatal(http.ListenAndServe(":8080", corsRouter))
}
