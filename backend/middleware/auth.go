package middleware

import (
	"log"
	"net/http"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil || cookie == "" {
			log.Println("Authorization cookie missing or malformed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}
		tokenStr := cookie

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			log.Printf("Invalid JWT: %v", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("jwtPayload", &models.User{
			ID:   claims.UserID,
			Role: models.Role(claims.Role),
		})
		c.Next()
	}

}
