package utils

import (
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/AltSumpreme/Medistream.git/models"

	"github.com/gin-gonic/gin"
)

func RoleChecker(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetCurrentUser(c)
		if err != nil {
			log.Println("Error getting current user:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authenticated"})
			c.Abort()
			return
		}

		if !slices.Contains(allowedRoles, models.Role(user.Role)) {
			log.Printf("Access denied for user %s with role %s", user.ID, user.Role)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}

}

func GetCurrentUser(c *gin.Context) (*JWTClaims, error) {
	val, exists := c.Get("jwtPayload")
	if !exists {
		log.Println("JWT payload not found in context")
		return nil, errors.New("JWT payload not found in context")
	}
	user, ok := val.(*JWTClaims)
	if !ok {
		log.Println("Invalid JWT payload type in context")
		return nil, errors.New("invalid JWT payload type in context")
	}
	return user, nil
}
