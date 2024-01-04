package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}

// GenerateJWT generates a JSON Web Token (JWT) for the given email and username.
//
// Parameters:
// - email: a string representing the email of the user.
// - username: a string representing the username of the user.
//
// Returns:
// - tokenString: a string representing the generated JWT.
// - err: an error object indicating any error that occurred during JWT generation.
func GenerateJWT(email string, id int32) (tokenString string, err error) {
	fmt.Printf("printing id in generate %v \n", strconv.Itoa(int(id)))

	expirationTime := time.Now().Add(50 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		ID:    strconv.Itoa(int(id)),

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}
func ValidateToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "", errors.New("couldn't parse claims")
	}

	// log.Printf("printing claims %v ", claims.ID)
	// log.Print(claims.StandardClaims)

	// fmt.Println("Email:", claims.Email)
	// fmt.Println("ID:", claims.ID)

	if claims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("token expired")
	}
	// log.Println("printing id " + claims.ID)

	return claims.ID, nil
}
