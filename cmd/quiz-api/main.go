package main

import (
	"log"

	"github.com/ayan-sh03/triviagenious-backend/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Print("Environment Variable loaded successfully")
	}

	app.Run()

}
