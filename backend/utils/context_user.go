package utils

import (
	"errors"
	"net/http"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/gin-gonic/gin"
)

func RoleChecker(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetCurrentUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		for _, r := range allowedRoles {
			if user.Role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}

func GetCurrentUser(c *gin.Context) (*models.User, error) {
	val, exists := c.Get("jwtPayload")
	if !exists {
		return nil, errors.New("user not found in context")
	}
	user, ok := val.(*models.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
