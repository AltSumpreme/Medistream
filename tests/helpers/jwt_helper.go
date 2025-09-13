package helpers

import (
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func InjectJWT(claims *utils.JWTClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	}
}
