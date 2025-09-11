package helpers

import (
	"time"

	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupRouterWithClaims(baseRouter *gin.Engine, userID uuid.UUID, role string) *gin.Engine {
	// Injecting Jwt claims into the context

	baseRouter.Use(func(c *gin.Context) {
		claims := &utils.JWTClaims{
			UserID: userID,
			Role:   role,
			Exp:    time.Now().Add(time.Hour).Unix(),
		}
		c.Set("jwtPayload", claims)
		c.Next()
	})
	return baseRouter
}
