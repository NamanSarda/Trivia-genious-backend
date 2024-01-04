package aicontroller

import (
	"fmt"
	"net/http"

	"github.com/ayan-sh03/triviagenious-backend/internal/util"
)

// 1. take file input
// 2. convert it into text
/// 3. feed it into the api

func GetQuestionFromAi(w http.ResponseWriter, r *http.Request) {

	response := util.CopyFileFromRequest(w, r)

	if response == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Error Reading File")
		return
	}

	extractedText := util.ExtractTextFromPdf(response.Name())

	fmt.Println(extractedText)
	res := util.ExecuteQuery(extractedText)

	if res == "" {
		util.RespondWithError(w, http.StatusInternalServerError, "Could not fetch from AI")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, res)

}
