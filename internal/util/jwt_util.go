package util

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	jwt.StandardClaims
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
