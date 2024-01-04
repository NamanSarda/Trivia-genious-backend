package authcontroller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ayan-sh03/triviagenious-backend/internal/api/auth/helpers"
	"github.com/ayan-sh03/triviagenious-backend/internal/api/models"
	query "github.com/ayan-sh03/triviagenious-backend/internal/api/quiz/sql"
	"github.com/ayan-sh03/triviagenious-backend/internal/util"
)

func RegisterUserController(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())

		return
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = hashedPassword

	dberr := query.CreateUser(&user)

	if dberr != nil {
		log.Println("Error occured creating user", dberr)
		if strings.Contains(dberr.Error(), "\"users_email_key\"") {
			util.RespondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		util.RespondWithError(w, http.StatusInternalServerError, dberr.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"user": user}) // User Model
}

func LoginController(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbUser, userErr := helpers.CheckEmailExists(user.Email)
	if userErr != nil {
		util.RespondWithError(w, http.StatusInternalServerError, userErr.Error())
		return
	}

	// check password
	securityErr := util.CheckPassword(user.Password, dbUser.Password)
	if securityErr != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials : "+securityErr.Error())
		return
	}

	token, err := util.GenerateJWT(dbUser.Email, dbUser.ID) //! Changed

	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Internal server error : Error While generating Token"+err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]string{"token": "Bearer " + token}) // Only token
}
