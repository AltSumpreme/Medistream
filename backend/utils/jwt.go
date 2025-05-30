package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    "user",
		"exp":     time.Now().Add(time.Hour * 4).Unix(), // Token expiration time set to 4 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
