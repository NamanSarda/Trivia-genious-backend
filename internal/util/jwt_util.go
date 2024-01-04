package util

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UnixTime int64

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	var unixTime float64
	if err := json.Unmarshal(b, &unixTime); err != nil {
		return err
	}
	*t = UnixTime(int64(unixTime))
	return nil
}

type JWTClaim struct {
	Username string   `json:"username"`
	UserId   string   `json:"user_id"`
	IssuedAt UnixTime `json:"iat"`
	Email    string   `json:"email"`
	jwt.StandardClaims
}

func ValidateToken(signedToken string) (string, error) {
	log.Printf("Token: %v", signedToken)
	// log.Printf("Claims: %v", claims)

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "", errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("token expired")
	}

	log.Print("Done here!")

	return claims.UserId, nil
}
