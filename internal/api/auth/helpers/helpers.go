package helpers

import (
	"database/sql"
	"log"

	"github.com/ayan-sh03/triviagenious-backend/config"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/models"
)

func CheckEmailExists(emailToCheck string) (*models.User, error) {
	var user models.User
	tx := config.Connect()
	err := tx.Get(&user, "SELECT id, email, username, password FROM users WHERE email = $1", emailToCheck)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with the given email
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}

	// User found with the given email
	return &user, nil
}
