package middleware

import (
	"fmt"
	"net/http"

	"github.com/ayan-sh03/triviagenious-backend/internal/util"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: request does not contain an access token", http.StatusUnauthorized)
			return
		}

		userId, err := util.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Log the retrieved user ID before setting it in the cookie
		fmt.Printf("Retrieved User ID from token: %v\n", userId)

		// Set the user ID in a cookie
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: userId, // Convert the user ID to string
		}
		http.SetCookie(w, cookie)

		next.ServeHTTP(w, r)
	})
}
