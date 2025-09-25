package factories

import (
	"log"
	"os"
	"time"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID string, role string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SECRET_KEY")

	tok, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("failed to generate JWT: %v", err)
	}

	return tok
}
func MakeJWT(userID uuid.UUID, role models.Role) *utils.JWTClaims {
	token := GenerateJWT(userID.String(), string(role))
	claims, _ := utils.ValidateJWT(token)
	return claims
}
