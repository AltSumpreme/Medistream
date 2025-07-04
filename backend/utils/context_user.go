package utils

import (
	"errors"
	"net/http"
	"slices"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/gin-gonic/gin"
)

func RoleChecker(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetCurrentUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authenticated"})
			c.Abort()
			return
		}

		if !slices.Contains(allowedRoles, user.Role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}

}

func GetCurrentUser(c *gin.Context) (*models.User, error) {
	val, exists := c.Get("jwtPayload")
	if !exists {
		return nil, errors.New("JWT payload not found in context")
	}
	user, ok := val.(*models.User)
	if !ok {
		return nil, errors.New("invalid JWT payload type in context")
	}
	return user, nil
}
