package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/auth"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {

	{
		rg.POST("/signup", auth.SignUp)
		rg.POST("/login", auth.Login)
		// rg.POST("/verify", auth.VerifyToken)
		rg.POST("/refresh", auth.RefreshAccessToken)
		rg.POST("/logout", auth.Logout)

	}
}
