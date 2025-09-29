package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("Invalid JWT: %v", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("jwtPayload", &utils.JWTClaims{
			UserID: claims.UserID,
			Role:   claims.Role,
			Exp:    claims.Exp,
		})
		c.Next()
	}

}
